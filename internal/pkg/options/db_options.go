package options

import (
	"financial_statement/internal/pkg/server"

	"github.com/spf13/pflag"
	"gorm.io/gorm"
)

type IDbOptions interface {
	// idata int
	Validate() []error
	AddFlags(fs *pflag.FlagSet)
	ApplyTo(c *server.Config) error
	NewClient() (*gorm.DB, error)
}
