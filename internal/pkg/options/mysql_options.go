package options

import (
	"financial_statement/internal/pkg/server"
	"financial_statement/pkg/db"
	"time"

	"github.com/spf13/pflag"
	"gorm.io/gorm"
)

type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Port                  int           `json:"port,omitempty"                     mapstructure:"port"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

// func (o *MySQLOptions)NewMySQLOptions() IDbOptions {
// 	return &MySQLOptions{
// 		Host:                  "127.0.0.1:3306",
// 		Username:              "",
// 		Password:              "",
// 		Database:              "",
// 		MaxIdleConnections:    100,
// 		MaxOpenConnections:    100,
// 		MaxConnectionLifeTime: time.Duration(10) * time.Second,
// 		LogLevel:              1, // Silent
// 	}
// }

// Validate 添加参数校验
func (o *MySQLOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Host, "mysql.host", o.Host, "MySQL 服务主机地址")
	fs.StringVar(&o.Username, "mysql.username", o.Username, "MySQL 服务用户名")
	fs.IntVar(&o.Port, "mysql.port", o.Port, "MySQL Port")
	fs.StringVar(&o.Password, "mysql.password", o.Password, "MySQL 服务密码")
	fs.StringVar(&o.Database, "mysql.database", o.Database, "MySQL 数据库名称")
	fs.IntVar(&o.MaxIdleConnections, "mysql.max-idle-connections", o.MaxOpenConnections, ""+
		"允许连接到 MySQL 的最大空闲连接数")
	fs.IntVar(&o.MaxOpenConnections, "mysql.max-open-connections", o.MaxOpenConnections, ""+
		"允许连接到 MySQL 的最大打开连接数。")
	fs.DurationVar(&o.MaxConnectionLifeTime, "mysql.max-connection-life-time", o.MaxConnectionLifeTime, ""+
		"允许连接到 MySQL 的最长连接生存时间")
	fs.IntVar(&o.LogLevel, "mysql.log-mode", o.LogLevel, "指定 gorm 日志级别")
}

func (o *MySQLOptions) ApplyTo(c *server.Config) error {
	return nil
}

// NewClient 使用给定的配置创建 mysql 实例
func (o *MySQLOptions) NewClient() (*gorm.DB, error) {
	opts := &db.MysqlOptions{
		Options: db.Options{
			Host:                  o.Host,
			Username:              o.Username,
			Port:                  o.Port,
			Password:              o.Password,
			Database:              o.Database,
			MaxIdleConnections:    o.MaxIdleConnections,
			MaxOpenConnections:    o.MaxOpenConnections,
			MaxConnectionLifeTime: o.MaxConnectionLifeTime,
			LogLevel:              o.LogLevel,
		},
	}

	return opts.New()
}
