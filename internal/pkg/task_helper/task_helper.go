package taskhelper

import (
	"encoding/json"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/types/task"
	"financial_statement/internal/pkg/storage"
	"financial_statement/pkg/log"
	"fmt"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/spf13/viper"
)

var (
	once               sync.Once
	TaskResultSaveType string
)

// 将识别结果存储
func SaveTaskResult(dbTask *model.Task, result []byte) (path string, err error) {
	once.Do(func() {
		TaskResultSaveType = viper.GetString("server.task-result-storage")
		if len(TaskResultSaveType) == 0 {
			TaskResultSaveType = "db"
		}
	})
	//如果存储方式为db，则直接将result返回
	if TaskResultSaveType == "db" {
		return string(result), nil
	}
	fileDir := filepath.Join(strconv.FormatInt(int64(dbTask.CreaterUserID), 10), strconv.FormatInt(int64(dbTask.ID), 10))
	filePath := filepath.Join(fileDir, fmt.Sprintf("%d.json", dbTask.ID))
	if err = storage.FileStorage.Save(result, filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

// 读取识别结果，读取时为了兼容旧数据，可以先尝试json反序列化数据，失败的话再去文件服务器拿数据
func GetTaskResult(dbTask *model.Task) (result []task.TaskFinancialStatement, err error) {
	financialStatement := make([]task.TaskFinancialStatement, 0)
	if dbTask.StandardResult != nil && len(*dbTask.StandardResult) > 0 {
		if err := json.Unmarshal([]byte(*dbTask.StandardResult), &financialStatement); err != nil {
			log.Infof("从数据库中解析任务%d的财报识别结果失败：%s ，尝试从文件服务器读取;", dbTask.ID, err.Error())
			var resultBytes []byte
			if resultBytes, err = storage.FileStorage.Get(*dbTask.StandardResult); err != nil {
				log.Errorf("从数据库以及文件服务器读取任务%d财报识别结果失败%s", dbTask.ID, err.Error())
			} else {
				if err := json.Unmarshal(resultBytes, &financialStatement); err != nil {
					log.Errorf("从数据库以及文件服务器读取任务%d财报识别结果失败%s", dbTask.ID, err.Error())
				}
			}
		}
	}
	return financialStatement, nil
}
