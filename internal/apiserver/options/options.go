package options

import (
	genericOptions "financial_statement/internal/pkg/options"
	cliflag "financial_statement/pkg/cli"

	"financial_statement/pkg/log"
)

// Options 输入命令配置参数
type Options struct {
	ServerOptions       *genericOptions.ServerOptions       `json:"server"   mapstructure:"server"`  // 基础 Server 配置
	MySQLOptions        *genericOptions.MySQLOptions        `json:"mysql"    mapstructure:"mysql"`   // MySql 配置
	DM8Options          *genericOptions.DM8Options          `json:"dameng"    mapstructure:"dameng"` // dameng 配置
	RedisOptions        *genericOptions.RedisOptions        `json:"redis"    mapstructure:"redis"`
	FeatureOptions      *genericOptions.FeatureOptions      `json:"feature"  mapstructure:"feature"`                 // 性能、指标监测功能
	Log                 *log.Options                        `json:"log"      mapstructure:"log"`                     // 日志配置
	LocalStorageOptions *genericOptions.LocalStorageOptions `json:"local-storage"      mapstructure:"local-storage"` // 文件存储配置
	NfsStorageOptions   *genericOptions.NfsStorageOptions   `json:"nfs"      mapstructure:"nfs"`                     // 文件存储配置
}

func (o Options) Flags() (fss cliflag.NamedFlagSets) {
	o.ServerOptions.AddFlags(fss.FlagSet("server"))
	o.MySQLOptions.AddFlags(fss.FlagSet("mysql"))
	o.DM8Options.AddFlags(fss.FlagSet("dameng"))
	o.RedisOptions.AddFlags(fss.FlagSet("redis"))
	o.FeatureOptions.AddFlags(fss.FlagSet("features"))
	o.Log.AddFlags(fss.FlagSet("log"))
	o.LocalStorageOptions.AddFlags(fss.FlagSet("local-storage"))
	o.NfsStorageOptions.AddFlags(fss.FlagSet("nfs"))
	return fss
}

func NewOptions() *Options {
	return &Options{
		ServerOptions:       genericOptions.NewServerRunOptions(),
		MySQLOptions:        &genericOptions.MySQLOptions{},
		DM8Options:          &genericOptions.DM8Options{},
		FeatureOptions:      genericOptions.NewFeatureOptions(),
		RedisOptions:        genericOptions.NewRedisOptions(),
		LocalStorageOptions: &genericOptions.LocalStorageOptions{},
		NfsStorageOptions:   &genericOptions.NfsStorageOptions{},
		Log:                 log.NewOptions(),
	}
}

// Validate 校验每个命令分组的输入是否正确
func (o *Options) Validate() []error {
	var errs []error

	errs = append(errs, o.ServerOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)
	errs = append(errs, o.DM8Options.Validate()...)
	errs = append(errs, o.RedisOptions.Validate()...)
	errs = append(errs, o.FeatureOptions.Validate()...)
	errs = append(errs, o.Log.Validate()...)
	errs = append(errs, o.LocalStorageOptions.Validate()...)
	errs = append(errs, o.NfsStorageOptions.Validate()...)

	return errs
}
