package options

import (
	"financial_statement/internal/pkg/server"

	"github.com/spf13/pflag"
)

type NfsStorageOptions struct {
	Download string `json:"download" mapstructure:"download"`
	Upload   string `json:"upload" mapstructure:"upload"`
	Bucket   string `json:"bucket" mapstructure:"bucket"`
}

// ApplyTo 从命令行读取到 Config 中
func (s *NfsStorageOptions) ApplyTo(c *server.Config) error {
	return nil
}

// Validate 添加参数校验
func (s *NfsStorageOptions) Validate() []error {
	var errors []error
	return errors
}

// AddFlags 将特定APIServer的标志添加到指定的标志集中
func (s *NfsStorageOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Download, "nfs.download", s.Download, "网络文件存储下载路径")
	fs.StringVar(&s.Bucket, "nfs.bucket", s.Bucket, "网络文件存储bucket")
	fs.StringVar(&s.Upload, "nfs.upload", s.Upload, "网络文件上传路径")
}
