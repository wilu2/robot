package config

import (
	"financial_statement/internal/apiserver/options"

	"github.com/spf13/viper"
)

type Config struct {
	*options.Options
	Nfs *NfsConfig `json:"nfs" mapstructure:"nfs"`
}

// CreateConfigFromOptions 创建一个配置实例
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	conf := &Config{
		Options: opts,
		Nfs:     nil,
	}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// NfsConfig 存储文件配置信息
type NfsConfig struct {
	Download string `json:"download,omitempty" mapstructure:"download"`
	Upload   string `json:"upload,omitempty" mapstructure:"upload"`
	Bucket   string `json:"bucket,omitempty" mapstructure:"bucket"`
}
