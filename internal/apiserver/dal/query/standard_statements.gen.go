// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"financial_statement/internal/apiserver/dal/model"
)

func newStandardStatement(db *gorm.DB) standardStatement {
	_standardStatement := standardStatement{}

	_standardStatement.standardStatementDo.UseDB(db)
	_standardStatement.standardStatementDo.UseModel(&model.StandardStatement{})

	tableName := _standardStatement.standardStatementDo.TableName()
	_standardStatement.ALL = field.NewField(tableName, "*")
	_standardStatement.ID = field.NewUint32(tableName, "id")
	_standardStatement.StandardID = field.NewUint32(tableName, "standard_id")
	_standardStatement.Type = field.NewInt32(tableName, "type")
	_standardStatement.Status = field.NewInt32(tableName, "status")
	_standardStatement.TitleStatus = field.NewInt32(tableName, "title_status")
	_standardStatement.FormulaStatus = field.NewInt32(tableName, "formula_status")
	_standardStatement.CreateAt = field.NewInt64(tableName, "create_at")
	_standardStatement.UpdateAt = field.NewInt64(tableName, "update_at")

	_standardStatement.fillFieldMap()

	return _standardStatement
}

type standardStatement struct {
	standardStatementDo standardStatementDo

	ALL           field.Field
	ID            field.Uint32
	StandardID    field.Uint32
	Type          field.Int32
	Status        field.Int32
	TitleStatus   field.Int32
	FormulaStatus field.Int32
	CreateAt      field.Int64
	UpdateAt      field.Int64

	fieldMap map[string]field.Expr
}

func (s standardStatement) Table(newTableName string) *standardStatement {
	s.standardStatementDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s standardStatement) As(alias string) *standardStatement {
	s.standardStatementDo.DO = *(s.standardStatementDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *standardStatement) updateTableName(table string) *standardStatement {
	s.ALL = field.NewField(table, "*")
	s.ID = field.NewUint32(table, "id")
	s.StandardID = field.NewUint32(table, "standard_id")
	s.Type = field.NewInt32(table, "type")
	s.Status = field.NewInt32(table, "status")
	s.TitleStatus = field.NewInt32(table, "title_status")
	s.FormulaStatus = field.NewInt32(table, "formula_status")
	s.CreateAt = field.NewInt64(table, "create_at")
	s.UpdateAt = field.NewInt64(table, "update_at")

	s.fillFieldMap()

	return s
}

func (s *standardStatement) WithContext(ctx context.Context) *standardStatementDo {
	return s.standardStatementDo.WithContext(ctx)
}

func (s standardStatement) TableName() string { return s.standardStatementDo.TableName() }

func (s standardStatement) Alias() string { return s.standardStatementDo.Alias() }

func (s *standardStatement) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *standardStatement) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 8)
	s.fieldMap["id"] = s.ID
	s.fieldMap["standard_id"] = s.StandardID
	s.fieldMap["type"] = s.Type
	s.fieldMap["status"] = s.Status
	s.fieldMap["title_status"] = s.TitleStatus
	s.fieldMap["formula_status"] = s.FormulaStatus
	s.fieldMap["create_at"] = s.CreateAt
	s.fieldMap["update_at"] = s.UpdateAt
}

func (s standardStatement) clone(db *gorm.DB) standardStatement {
	s.standardStatementDo.ReplaceDB(db)
	return s
}

type standardStatementDo struct{ gen.DO }

func (s standardStatementDo) Debug() *standardStatementDo {
	return s.withDO(s.DO.Debug())
}

func (s standardStatementDo) WithContext(ctx context.Context) *standardStatementDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s standardStatementDo) ReadDB() *standardStatementDo {
	return s.Clauses(dbresolver.Read)
}

func (s standardStatementDo) WriteDB() *standardStatementDo {
	return s.Clauses(dbresolver.Write)
}

func (s standardStatementDo) Clauses(conds ...clause.Expression) *standardStatementDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s standardStatementDo) Returning(value interface{}, columns ...string) *standardStatementDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s standardStatementDo) Not(conds ...gen.Condition) *standardStatementDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s standardStatementDo) Or(conds ...gen.Condition) *standardStatementDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s standardStatementDo) Select(conds ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s standardStatementDo) Where(conds ...gen.Condition) *standardStatementDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s standardStatementDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *standardStatementDo {
	return s.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (s standardStatementDo) Order(conds ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s standardStatementDo) Distinct(cols ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s standardStatementDo) Omit(cols ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s standardStatementDo) Join(table schema.Tabler, on ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s standardStatementDo) LeftJoin(table schema.Tabler, on ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s standardStatementDo) RightJoin(table schema.Tabler, on ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s standardStatementDo) Group(cols ...field.Expr) *standardStatementDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s standardStatementDo) Having(conds ...gen.Condition) *standardStatementDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s standardStatementDo) Limit(limit int) *standardStatementDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s standardStatementDo) Offset(offset int) *standardStatementDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s standardStatementDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *standardStatementDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s standardStatementDo) Unscoped() *standardStatementDo {
	return s.withDO(s.DO.Unscoped())
}

func (s standardStatementDo) Create(values ...*model.StandardStatement) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s standardStatementDo) CreateInBatches(values []*model.StandardStatement, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s standardStatementDo) Save(values ...*model.StandardStatement) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s standardStatementDo) First() (*model.StandardStatement, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.StandardStatement), nil
	}
}

func (s standardStatementDo) Take() (*model.StandardStatement, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.StandardStatement), nil
	}
}

func (s standardStatementDo) Last() (*model.StandardStatement, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.StandardStatement), nil
	}
}

func (s standardStatementDo) Find() ([]*model.StandardStatement, error) {
	result, err := s.DO.Find()
	return result.([]*model.StandardStatement), err
}

func (s standardStatementDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.StandardStatement, err error) {
	buf := make([]*model.StandardStatement, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s standardStatementDo) FindInBatches(result *[]*model.StandardStatement, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s standardStatementDo) Attrs(attrs ...field.AssignExpr) *standardStatementDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s standardStatementDo) Assign(attrs ...field.AssignExpr) *standardStatementDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s standardStatementDo) Joins(fields ...field.RelationField) *standardStatementDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s standardStatementDo) Preload(fields ...field.RelationField) *standardStatementDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s standardStatementDo) FirstOrInit() (*model.StandardStatement, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.StandardStatement), nil
	}
}

func (s standardStatementDo) FirstOrCreate() (*model.StandardStatement, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.StandardStatement), nil
	}
}

func (s standardStatementDo) FindByPage(offset int, limit int) (result []*model.StandardStatement, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s standardStatementDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s standardStatementDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s standardStatementDo) Delete(models ...*model.StandardStatement) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *standardStatementDo) withDO(do gen.Dao) *standardStatementDo {
	s.DO = *do.(*gen.DO)
	return s
}