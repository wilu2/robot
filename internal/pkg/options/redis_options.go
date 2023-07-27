package options

import (
	"time"

	"github.com/spf13/pflag"
)

// RedisOptions defines options for redis cluster.
type RedisOptions struct {
	Type             string        `json:"type"						mapstructure:"type"`
	MasterName       string        `json:"masterName"				mapstructure:"masterName"`
	SentinelAddrs    string        `json:"sentinelAddrs"			mapstructure:"sentinelAddrs"`
	SentinelUsername string        `json:"sentinelUsername"			mapstructure:"sentinelUsername"`
	SentinelPassword string        `json:"sentinelPassword"			mapstructure:"sentinelPassword"`
	ClusterAddrs     string        `json:"clusterAddrs"				mapstructure:"clusterAddrs"`
	Host             string        `json:"host"                     mapstructure:"host"                     description:"Redis service host address"`
	Port             int           `json:"port"                     mapstructure:"port"`
	Database         int           `json:"database"					mapstructure:"database"`
	Username         string        `json:"username"                 mapstructure:"username"`
	Password         string        `json:"password"                 mapstructure:"password"`
	PoolSize         int           `json:"poll-size"                mapstructure:"poll-size"`
	MinIdleConns     int           `json:"min-idle-conns"           mapstructure:"min-idle-conns"`
	DialTimeOut      time.Duration `json:"dial-timeout"             mapstructure:"dial-timeout"`
	ReadTimeOut      time.Duration `json:"read-timeout"             mapstructure:"read-timeout"`
	WriteTimeOut     time.Duration `json:"write-timeout"            mapstructure:"write-timeout"`
	PoolTimeOut      time.Duration `json:"pool-timeout"             mapstructure:"pool-timeout"`
	IdleCheckFreq    time.Duration `json:"idle-check-frequency"     mapstructure:"idle-check-frequency"`
	IdleTimeout      time.Duration `json:"idle-timeout"             mapstructure:"idle-timeout"`
}

// NewRedisOptions create a `zero` value instance.
func NewRedisOptions() *RedisOptions {
	return &RedisOptions{
		Type:             "default",
		MasterName:       "master",
		SentinelAddrs:    "",
		SentinelUsername: "",
		SentinelPassword: "",
		ClusterAddrs:     "",
		Host:             "127.0.0.1",
		Port:             6379,
		Username:         "",
		Password:         "",
		Database:         0,
		PoolSize:         10,
		MinIdleConns:     5,
		DialTimeOut:      time.Second * 5,
		ReadTimeOut:      time.Second * 3,
		WriteTimeOut:     time.Second * 3,
		PoolTimeOut:      time.Second * 5,
		IdleCheckFreq:    time.Second * 5,
		IdleTimeout:      time.Second * 10,
	}
}

// Validate verifies flags passed to RedisOptions.
func (o *RedisOptions) Validate() []error {
	var errs []error

	return errs
}

// AddFlags adds flags related to redis storage for a specific APIServer to the specified FlagSet.
func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.Type, "redis.type", o.Type, "Redis 部署类型")
	fs.StringVar(&o.Host, "redis.host", o.Host, "Redis 服务器地址")
	fs.IntVar(&o.Port, "redis.port", o.Port, "Redis 服务器侦听的端口")
	fs.IntVar(&o.Database, "redis.database", o.Database, "默认情况下，数据库为0，集群必须设为 0")
	fs.StringVar(&o.Username, "redis.username", o.Username, "Redis 服务器的用户名")
	fs.StringVar(&o.Password, "redis.password", o.Password, "Redis 服务器的密码")
	fs.StringVar(&o.MasterName, "redis.masterName", o.MasterName, "Redis 哨兵模式的master-name")
	fs.StringVar(&o.SentinelAddrs, "redis.sentinelAddrs", o.SentinelAddrs, "Redis 哨兵模式的地址")
	fs.StringVar(&o.SentinelUsername, "redis.sentinelUsername", o.SentinelUsername, "Redis 哨兵模式的账号")
	fs.StringVar(&o.SentinelPassword, "redis.sentinelPassword", o.SentinelPassword, "Redis 哨兵模式的密码")
	fs.StringVar(&o.ClusterAddrs, "redis.clusterAddrs", o.ClusterAddrs, "Redis 集群模式下的地址")
	fs.IntVar(&o.PoolSize, "redis.pool-size", o.PoolSize, "Redis 连接池的大小")
	// fs.IntVar(&o.MinIdleConns, "redis.min-idle-conns", o.MinIdleConns, "Redis 连接池的大小")
	// fs.IntVar(&o.DialTimeOut, "redis.dial-time-out", o.DialTimeOut, "Redis 连接池的大小")
	// fs.IntVar(&o.PoolSize, "redis.read-time-out", o.PoolSize, "Redis 连接池的大小")
	// fs.IntVar(&o.PoolSize, "redis.write-time-out", o.PoolSize, "Redis 连接池的大小")
	// fs.IntVar(&o.PoolSize, "redis.pool-time-out", o.PoolSize, "Redis 连接池的大小")
	// fs.IntVar(&o.PoolSize, "redis.idle-check-frequency", o.PoolSize, "Redis 连接池的大小")
	// fs.IntVar(&o.PoolSize, "idle-time-out", o.PoolSize, "Redis 连接池的大小")
}
