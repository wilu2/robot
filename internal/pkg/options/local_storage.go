package options

import (
	"financial_statement/internal/pkg/server"

	"github.com/spf13/pflag"
)

type LocalStorageOptions struct {
	SaveDir string `json:"save-dir" mapstructure:"save-dir"`
}

// ApplyTo 从命令行读取到 Config 中
func (s *LocalStorageOptions) ApplyTo(c *server.Config) error {
	return nil
}

// Validate 添加参数校验
func (s *LocalStorageOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags 将特定APIServer的标志添加到指定的标志集中
func (s *LocalStorageOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.SaveDir, "local-storage.save-dir", s.SaveDir, "文件本地存储路径")
}
