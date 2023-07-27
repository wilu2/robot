package job

import (
	"context"
	"financial_statement/internal/apiserver/dal"
	filecleaner "financial_statement/internal/apiserver/job/file-cleaner"
	jobTask "financial_statement/internal/apiserver/job/task"
	taskresultasync "financial_statement/internal/apiserver/job/task-result-async"
	"financial_statement/pkg/log"

	"github.com/hibiken/asynq"
)

type JobMgr struct {
	AsynqClient *asynq.Client
	AsynqServer *asynq.Server
}

//初始化任务管理器
func InitJobMgr() *JobMgr {
	j := new(JobMgr)

	asynqredisIns, _ := dal.GetRedisFactoryOr(nil)
	asynqRedis := asynqredisIns.GetAsynqRedis()
	asynqRedisConnOpt := asynqredisIns.GetAsynqRedisConnOpt()
	//Web站点时服务端也是客户端
	//初始化异步任务服务端对象
	j.AsynqServer = asynq.NewServer(*asynqRedisConnOpt, asynq.Config{
		//Task处理失败时的处理
		ErrorHandler: asynq.ErrorHandlerFunc(handleError),
	})
	//初始化异步任务客户端对象
	j.AsynqClient = asynqRedis

	//初始化历史任务清理定时任务
	initPeriodicTask(*asynqRedisConnOpt)

	//初始化异步任务处理Handle
	mux := asynq.NewServeMux()
	//注册异步任务处理函数
	mux.HandleFunc(jobTask.RecognizeTask, jobTask.HandleRecognizeTask)
	//注册文件清理处理函数
	mux.HandleFunc(filecleaner.Filecleaner, filecleaner.HandleCleanerTask)
	mux.HandleFunc(taskresultasync.TaskAsync, taskresultasync.HandleAsyncTask)

	go func() {
		if err := j.AsynqServer.Run(mux); err != nil {
			log.Errorf("AsynqServer Run with error:%s", err.Error())
			panic(err)
		}
	}()
	return j
}

//添加一个异步处理任务到任务队列
func (j *JobMgr) EnqueueTask(t *asynq.Task) error {
	_, err := j.AsynqClient.Enqueue(t)
	if err != nil {
		log.Errorf("asynq enqueue task with error:%s", err.Error())
		return err
	}
	return nil
}

func handleError(ctx context.Context, task *asynq.Task, err error) {
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry {
		log.Errorf("retry exhausted for task %s: %w", task.Type(), err)
		switch task.Type() {
		case jobTask.RecognizeTask:
			jobTask.ErrorHandler(ctx, task.Payload())
		}
	}
}
