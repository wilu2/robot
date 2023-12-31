// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameStandardStatementTitle = "standard_statement_titles"

// StandardStatementTitle mapped from table <standard_statement_titles>
type StandardStatementTitle struct {
	ID          uint32  `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	StatementID uint32  `gorm:"column:statement_id;type:int(11) unsigned;not null" json:"statement_id"`
	Name        string  `gorm:"column:name;type:varchar(1024);not null" json:"name"`
	ExternalID  *string `gorm:"column:external_id;type:varchar(1024)" json:"external_id"`
	Aliases     *string `gorm:"column:aliases;type:varchar(1024)" json:"aliases"` // 别名，存储为json字符串（数据库字段类型为string，暂不用json是担心其他db时兼容性问题）
	/*
		状态：
		1：正常
		-1：停用或删除
	*/
	Status    int32 `gorm:"column:status;type:int(11);not null;default:1" json:"status"`
	OrderByID int32 `gorm:"column:order_by_id;type:int(11);not null" json:"order_by_id"` // 排序id，越小越靠前
	CreateAt  int64 `gorm:"column:create_at;type:bigint(20);not null" json:"create_at"`
	UpdateAt  int64 `gorm:"column:update_at;type:bigint(20);not null" json:"update_at"`
}

// TableName StandardStatementTitle's table name
func (*StandardStatementTitle) TableName() string {
	return TableNameStandardStatementTitle
}
