package db

import (
	dbdriver "financial_statement/pkg/dbDriver"
	"financial_statement/pkg/log"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Options defines optsions for mysql database.
type Dm8Options struct {
	Options
}

// New create a new gorm db instance with the given options.
func (opts *Dm8Options) New() (*gorm.DB, error) {
	dsn := fmt.Sprintf(`dm://%s:%s@%s:%d?connectTimeout=10000&compatibleMode=mysql&charset=1`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port)
	log.Debugf("使用dm8数据库，连接字符串：%s", dsn)

	db, err := gorm.Open(dbdriver.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   opts.Database + ".",
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}
