package dbdriver

import (
	"context"
	"database/sql"
	_ "dm"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const DriverName = "dm"

var (
	// CreateClauses create clauses
	CreateClauses = []string{"INSERT", "VALUES", "ON CONFLICT"}
	// QueryClauses query clauses
	QueryClauses = []string{}
	// UpdateClauses update clauses
	UpdateClauses = []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"}
	// DeleteClauses delete clauses
	DeleteClauses = []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"}

	defaultDatetimePrecision = 3
)

type DmDriver struct{}

// func init() {
// 	sql.Register(DriverName, &DmDriver{})
// }

// func (d DmDriver) Open(dsn string) (driver.Conn, error) {
// 	return nil, nil
// }

type Config struct {
	DSN        string
	DriverName string
	Conn       gorm.ConnPool
}

type Dialector struct {
	*Config
}

func Open(dsn string) gorm.Dialector {
	return &Dialector{Config: &Config{DSN: dsn}}
}

func New(config Config) gorm.Dialector {
	return &Dialector{Config: &config}
}

func (dialector Dialector) Name() string {
	return "dm8"
}

func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
	fmt.Printf("Initialize: %+v\n", dialector)

	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		CreateClauses: CreateClauses,
		QueryClauses:  QueryClauses,
		UpdateClauses: UpdateClauses,
		DeleteClauses: DeleteClauses,
	})

	if dialector.DriverName == "" {
		dialector.DriverName = DriverName
	}
	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		dialector.Conn, err = sql.Open("dm", dialector.DSN)
		if err != nil {
			return err
		}
		db.ConnPool = dialector.Conn
	}

	rows, err := db.ConnPool.QueryContext(context.Background(), "select * from v$version")
	if err != nil {
		return err
	}

	var version []string
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return err
		}
		version = append(version, v)
	}
	fmt.Printf("version: %+v\n", version)

	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{})

	return nil
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	fmt.Printf("Migrator: %+v\n", dialector)
	return nil
}

func (dialector Dialector) DataTypeOf(*schema.Field) string {
	fmt.Printf("DataTypeOf: %+v\n", dialector)
	return ""
}

func (dialector Dialector) DefaultValueOf(*schema.Field) clause.Expression {
	fmt.Printf("DefultValueOf: %+v\n", dialector)
	return clause.Expr{SQL: "DEFAULT"}
}

func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	fmt.Printf("BindVarTo: %+v\n", dialector)
	writer.WriteByte('?')
	return
}

func (dialector Dialector) QuoteTo(writer clause.Writer, str string) {
	fmt.Printf("QuoteTo: %+v\n", str)
	writer.WriteString(`"`)
	for _, v := range []byte(str) {
		switch v {
		case '"':
			writer.WriteString(`"`)
			writer.WriteByte(v)
		case '.':
			writer.WriteString(`"."`)
		default:
			writer.WriteByte(v)
		}
	}
	writer.WriteString(`"`)
}

func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	fmt.Printf("Explain: %+v\n", dialector)
	return logger.ExplainSQL(sql, nil, `"`, vars...)
}
