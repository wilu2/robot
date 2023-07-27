package db

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type IOptions interface {
	New() (*gorm.DB, error)
}

// Options defines optsions for mysql database.
type Options struct {
	Host                  string
	Port                  int
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
	Logger                logger.Interface
}
