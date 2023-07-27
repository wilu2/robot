package filecleaner

import (
	"context"
	"encoding/json"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	dblog "financial_statement/internal/apiserver/db_log"
	"financial_statement/internal/pkg/options"
	"financial_statement/internal/pkg/storage"
	"financial_statement/pkg/log"
	"fmt"
	"sync"
	"time"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
)

const (
	Filecleaner = "task:task-cleaner"
)

var (
	once                 sync.Once
	db                   dal.Repo
	backupdb             dal.Repo
	backupTask           bool
	backupTaskDir        string
	backupTaskDbUsername string
	backupTaskDbPassword string
	backupTaskDbHost     string
	backupTaskDbPort     int
	backupTaskDbDatabase string
)

func HandleCleanerTask(ctx context.Context, t *asynq.Task) error {
	log.Debugf("----->file-cleaner start job. time:%s", time.Now().String())
	defer func() { log.Debugf("----->file-cleaner end job. time:%s", time.Now().String()) }()
	return cleanerTask(context.Background())
}

func cleanerTask(ctx context.Context) error {
	once.Do(func() {
		db, _ = dal.GetDbFactoryOr(nil)
		backupTask = viper.GetBool("file-cleaner.backup-task")
		backupTaskDir = viper.GetString("file-cleaner.backup-task-dir")
		backupTaskDbUsername = viper.GetString("file-cleaner.backup-task-db-username")
		backupTaskDbPassword = viper.GetString("file-cleaner.backup-task-db-password")
		backupTaskDbHost = viper.GetString("file-cleaner.backup-task-db-host")
		backupTaskDbPort = viper.GetInt("file-cleaner.backup-task-db-port")
		backupTaskDbDatabase = viper.GetString("file-cleaner.backup-task-db-database")
		if backupTask {
			opt := options.MySQLOptions{
				Host:     backupTaskDbHost,
				Port:     backupTaskDbPort,
				Username: backupTaskDbUsername,
				Password: backupTaskDbPassword,
				Database: backupTaskDbDatabase,
			}
			backupDb, err := dal.GetDbFactoryOr(&opt)
			if err != nil {
				log.Infof("Back up mysql config:%v", opt)
				log.Panicf("Back up mysql cann't Connection:%s", err.Error())
			}
			backupdb = backupDb
		}

	})
	var (
		tTask = query.Use(db.GetDb()).Task
		tasks []struct {
			ID    uint32
			Files string
		}
		ids []uint32
	)

	//获取6个月之前的任务
	err := tTask.WithContext(ctx).Select(tTask.ID, tTask.Files).Where(tTask.Status.Neq(consts.TaskStatusDeleted), tTask.CreatedAt.Lt(time.Now().AddDate(0, -6, 0).Unix())).
		Scan(&tasks)
	if err != nil {
		log.Errorf("执行清理任务出错：%s", err.Error())
		return err
	}

	for _, task := range tasks {
		if len(task.Files) > 0 {
			var fileList []string
			json.Unmarshal([]byte(task.Files), &fileList)
			for _, file := range fileList {
				if backupTask && len(backupTaskDir) > 0 {
					//历史文件迁移功能
					if err := storage.FileStorage.Copy(file, backupTaskDir); err != nil {
						dblog.DbLog(ctx, task.ID, fmt.Sprintf("迁移文件时出错：%s", err.Error()))
						continue
					}
				}
				if backupTask && backupdb != nil {
					//历史文件数据备份功能
					tempTask, _ := tTask.WithContext(ctx).Where(tTask.ID.Eq(task.ID)).First()
					backupTtask := query.Use(backupdb.GetDb()).Task
					if err := backupTtask.WithContext(ctx).Create(&model.Task{
						TaskName:       tempTask.TaskName,
						FileFormat:     tempTask.FileFormat,
						StandardID:     tempTask.StandardID,
						ExternalInfo:   tempTask.ExternalInfo,
						Async:          tempTask.Async,
						AsyncStatus:    tempTask.AsyncStatus,
						Status:         tempTask.Status,
						CreaterUserID:  tempTask.CreaterUserID,
						StandardResult: tempTask.StandardResult,
						Files:          tempTask.Files,
						Error:          tempTask.Error,
						CreatedAt:      tempTask.CreatedAt,
						UpdatedAt:      tempTask.UpdatedAt,
					}); err != nil {
						dblog.DbLog(ctx, task.ID, fmt.Sprintf("迁移文件到备份db时出错：%s", err.Error()))
						continue
					}
				}
				if err := storage.FileStorage.Delete(file); err != nil {
					log.Errorf("删除文件失败：%s", err.Error())
				}
			}
		}
		ids = append(ids, task.ID)
	}

	if _, err := tTask.WithContext(ctx).Where(tTask.ID.In(ids...)).Update(tTask.Status, consts.TaskStatusDeleted); err != nil {
		log.Errorf("更新task为已删除失败：%s", err.Error())
		return err
	}
	return nil
}
