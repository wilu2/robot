package server

import (
	"net"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Host            string
	Port            int
	Mode            string
	Db              string
	FileStorage     string
	Middlewares     []string
	Health          bool
	Language        string
	EnableProfiling bool
	EnableMetrics   bool
	IsSaas          bool
}

func NewConfig() *Config {
	return &Config{
		Host:            "127.0.0.1",
		Port:            8080,
		Health:          true,
		Db:              "mysql",
		FileStorage:     "local",
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		Language:        "en",
		EnableProfiling: true,
		EnableMetrics:   true,
		IsSaas:          false,
	}
}

// CompletedConfig GenericApiServer的已完成配置
type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

// New 从配置中设置 GenericAPIServer 的新实例
func (c CompletedConfig) New() (*GenericApiServer, error) {
	gin.SetMode(c.Mode)

	s := &GenericApiServer{
		Address:         net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		health:          c.Health,
		middlewares:     c.Middlewares,
		enableProfiling: c.EnableProfiling,
		enableMetrics:   c.EnableMetrics,
		Engine:          gin.New(),
	}

	initGenericApiServer(s)

	return s, nil
}
