package orm

import (
	"errors"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func IsUniqueConstraintFailed(tx *gorm.DB) bool {
	switch tx.Dialector.Name() {
	case "mysql":
		var mysqlErr *mysqlDriver.MySQLError
		return errors.As(tx.Error, &mysqlErr) && mysqlErr.Number == 1062
	default:
		panic("unsupported database")
	}
}
