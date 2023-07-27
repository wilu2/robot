package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Options defines optsions for mysql database.
type MysqlOptions struct {
	Options
}

// New create a new gorm db instance with the given options.
func (opts *MysqlOptions) New() (*gorm.DB, error) {
	dsn := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=%t&loc=%s&timeout=10s`,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		true,
		"Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: opts.Logger,
	})
	if err != nil {
		return nil, err
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
