package dal

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	genericOptions "financial_statement/internal/pkg/options"
	"financial_statement/pkg/db"

	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Repo interface {
	GetDb() *gorm.DB
	CloseDb() error
}

type dbRepo struct {
	Db *gorm.DB
}

func (d *dbRepo) GetDb() *gorm.DB {
	return d.Db
}

func (d *dbRepo) CloseDb() error {
	db, err := d.Db.DB()
	if err != nil {
		return errors.Wrap(err, "get gorm db instance failed")
	}

	return db.Close()
}

var (
	dbFactory Repo
	once      sync.Once
)

func initMysqlDb(opts *genericOptions.MySQLOptions) (*gorm.DB, error) {
	options := &db.MysqlOptions{
		Options: db.Options{
			Host:                  opts.Host,
			Port:                  opts.Port,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.LogLevel(opts.LogLevel),
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			}),
		},
	}
	return options.New()
}

func initDm8Db(opts *genericOptions.DM8Options) (*gorm.DB, error) {
	options := &db.Dm8Options{
		Options: db.Options{
			Host:                  opts.Host,
			Port:                  opts.Port,
			Username:              opts.Username,
			Password:              opts.Password,
			Database:              opts.Database,
			MaxIdleConnections:    opts.MaxIdleConnections,
			MaxOpenConnections:    opts.MaxOpenConnections,
			MaxConnectionLifeTime: opts.MaxConnectionLifeTime,
			Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.LogLevel(opts.LogLevel),
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			}),
		},
	}
	return options.New()
}

// GetMySQLFactoryOr create mysql factory with the given config.
func GetDbFactoryOr(opts genericOptions.IDbOptions) (Repo, error) {
	if opts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store fatory")
	}

	var err error
	var dbIns *gorm.DB
	once.Do(func() {
		switch opts.(type) {
		case *genericOptions.MySQLOptions:
			dbIns, err = initMysqlDb(opts.(*genericOptions.MySQLOptions))
		case *genericOptions.DM8Options:
			dbIns, err = initDm8Db(opts.(*genericOptions.DM8Options))
		}

		// uncomment the following line if you need auto migration the given models
		// not suggested in production environment.
		// migrateDatabase(dbIns)
		dbFactory = &dbRepo{dbIns}
	})

	if err != nil {
		return dbFactory, fmt.Errorf("failed to get db,error:%s", err.Error())
	}

	return dbFactory, nil
}
