package options

import (
	"financial_statement/internal/pkg/server"

	"github.com/spf13/pflag"
)

type IStorageOptions interface {
	// idata int
	Validate() []error
	AddFlags(fs *pflag.FlagSet)
	ApplyTo(c *server.Config) error
}
