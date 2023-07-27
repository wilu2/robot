package taskresultasync

import (
	"bytes"
	"context"
	"encoding/json"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	dblog "financial_statement/internal/apiserver/db_log"
	"financial_statement/internal/apiserver/types/task"
	taskhelper "financial_statement/internal/pkg/task_helper"
	"financial_statement/pkg/log"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

const (
	TaskAsync              = "task:task-async"
	SyncTaskTypeRecognized = "1" //任务识别完成时的回调
	SyncTaskTypeSaved      = "2" //任务在编辑界面保存时的回调
	// SyncTaskTypeReIdentify = "3" //任务重新识别后的回调
)

var (
	once     sync.Once
	asyncUrl string
	db       dal.Repo
)

type SyncTaskInfo struct {
	TaskId       uint32
	TaskType     string
	CurrentIndex int
}

// 传递任务id，生成一个异步处理任务
func NewAsyncTask(taskPayload *SyncTaskInfo) (*asynq.TaskInfo, error) {
	saveWithSync := viper.GetBool("task-result-async.save-with-sync")
	if !saveWithSync && taskPayload.TaskType == SyncTaskTypeSaved {
		//如果是保存时的结果回调，但保存时回调开关没有开，则不创建异步任务
		log.Debugf("保存时的数据回调，但是回调开关没有开，所以不进行回调;config.saveWithSync:%v taskType:%v", saveWithSync, taskPayload)
		return nil, nil
	}
	asynqredisIns, _ := dal.GetRedisFactoryOr(nil)
	//初始化异步任务客户端对象
	client := asynqredisIns.GetAsynqRedis()
	payload, _ := json.Marshal(taskPayload)
	info, err := client.Enqueue(asynq.NewTask(TaskAsync, payload))
	return info, err
}

func HandleAsyncTask(ctx context.Context, t *asynq.Task) error {
	log.Debugf("----->Task async start job. time:%s", time.Now().String())
	defer func() { log.Debugf("----->Task async end job. time:%s", time.Now().String()) }()
	once.Do(func() {
		asyncUrl = viper.GetString("task-result-async.url")
		db, _ = dal.GetDbFactoryOr(nil)
	})
	if len(asyncUrl) == 0 {
		return nil
	}
	var sycnTaskInfo SyncTaskInfo
	err := json.Unmarshal(t.Payload(), &sycnTaskInfo)
	if err != nil {
		log.Errorf("HandleAsyncTask with error%s TaskId:%d\n", err.Error(), sycnTaskInfo.TaskId)
		return err
	}
	log.Debugf("start async task id : %d", sycnTaskInfo.TaskId)
	var (
		tTask = query.Use(db.GetDb()).Task
	)
	count, err := tTask.WithContext(context.Background()).Where(tTask.ID.Eq(sycnTaskInfo.TaskId), tTask.Async.Eq(1)).Count()
	if count == 0 && err == nil {
		return nil
	}
	dbTask, err := tTask.WithContext(context.Background()).Where(tTask.ID.Eq(sycnTaskInfo.TaskId), tTask.Async.Eq(1)).First()
	if err != nil {
		dblog.DbLog(context.Background(), dbTask.ID, fmt.Sprintf("Task Async 同步失败：%s", err.Error()))
		log.Errorf("获取Task Async失败：%s", err.Error())
		return err
	}
	if err := taskAsync(context.Background(), dbTask, &sycnTaskInfo); err != nil {
		dblog.DbLog(context.Background(), dbTask.ID, fmt.Sprintf("Task Async 同步失败：%s", err.Error()))
		tTask.WithContext(context.Background()).Where(tTask.ID.Eq(sycnTaskInfo.TaskId)).Update(tTask.AsyncStatus, consts.TaskAsyncStatusUnsynchronized)
		return err
	} else {
		tTask.WithContext(context.Background()).Where(tTask.ID.Eq(sycnTaskInfo.TaskId)).Update(tTask.AsyncStatus, consts.TaskAsyncStatusSynchronized)
		return nil
	}
}

func taskAsync(ctx context.Context, dbTask *model.Task, sycnTaskInfo *SyncTaskInfo) (err error) {

	data := task.GetTaskResp{}
	data.CreatedAt = int(dbTask.CreatedAt)
	// resp.External = task.ExternalInfo

	var externalInfo consts.ExternalInfo
	if err = json.Unmarshal([]byte(dbTask.ExternalInfo), &externalInfo); err != nil {
		log.Errorf("同步任务%d结果请求失败：%s;", data.TaskId, err.Error())
		return err
	}
	data.ExtraInfo = externalInfo.ExtraInfo
	for _, file := range externalInfo.FileInfo {
		if len(file) > 0 {
			data.External = append(data.External, file)
		}
	}

	data.FileFormat = dbTask.FileFormat
	data.Standard = task.Standard{
		Id:         dbTask.StandardID,
		ExternalId: "",
	}
	data.Status = int(dbTask.Status)
	data.TaskId = dbTask.ID
	data.UpdateAt = int(dbTask.UpdatedAt)
	data.OperType = "B"
	if data.FinancialStatement, err = taskhelper.GetTaskResult(dbTask); err != nil {
		return err
	}
	postDataStr, _ := json.Marshal(data)
	body := bytes.NewReader([]byte(postDataStr))

	httpReq, err := http.NewRequest("POST", asyncUrl, body)
	httpReq.Header.Add("Content-type", "application/json")

	q := httpReq.URL.Query()
	q.Add("sync_type", sycnTaskInfo.TaskType)
	q.Add("current_index", strconv.Itoa(sycnTaskInfo.CurrentIndex))
	httpReq.URL.RawQuery = q.Encode()

	log.Infof("同步任务%d到%s", data.TaskId, httpReq.URL)
	if err != nil {
		log.Errorf("同步任务%d结果请求失败：%s;", data.TaskId, err.Error())
		return err
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Errorf("同步任务%d结果请求失败：%s", data.TaskId, err.Error())
		return err
	}
	if httpResp.StatusCode != http.StatusOK {
		log.Errorf("同步任务%d结果请求失败，Http Status：%d", data.TaskId, httpResp.StatusCode)
		return err
	}
	defer httpResp.Body.Close()
	return nil
}
