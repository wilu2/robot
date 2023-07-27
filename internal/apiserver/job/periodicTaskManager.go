package job

import (
	filecleaner "financial_statement/internal/apiserver/job/file-cleaner"
	"financial_statement/pkg/log"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

var (
	once              sync.Once
	fileCleanerEnable bool
)

func initPeriodicTask(opt asynq.RedisConnOpt) {

	once.Do(func() {
		fileCleanerEnable = viper.GetBool("file-cleaner.enable")
	})

	scheduler := asynq.NewScheduler(opt, nil)

	if fileCleanerEnable {
		//filecleaner 定时任务 每天执行一次
		task := asynq.NewTask(filecleaner.Filecleaner, nil)
		entryID, err := scheduler.Register("0 0 1-31 * *", task) //分,时,日,月,星
		if err != nil {
			log.Errorf("initPeriodicTask with error:%s", err.Error())
		}
		log.Infof("registered an entry: %q\n", entryID)
		//filecleaner 定时任务 END
	}

	go func() {
		if err := scheduler.Run(); err != nil {
			log.Fatalf("定时任务运行出错：%s", err.Error())
		}
	}()
}
