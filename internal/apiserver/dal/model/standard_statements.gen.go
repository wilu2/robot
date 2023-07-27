// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameStandardStatement = "standard_statements"

// StandardStatement mapped from table <standard_statements>
type StandardStatement struct {
	ID         uint32 `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	StandardID uint32 `gorm:"column:standard_id;type:int(11) unsigned;not null" json:"standard_id"`
	/*
		财务准则报表类型：
		1：资产负债表
		2：利润表
		3：现金流量表
	*/
	Type int32 `gorm:"column:type;type:int(11);not null" json:"type"`
	/*
		状态：
		1：启用
		-1：停用
	*/
	Status int32 `gorm:"column:status;type:int(11);not null;default:1" json:"status"`
	/*
		科目配置状态：
		1：一配置
		-1：未配置
	*/
	TitleStatus int32 `gorm:"column:title_status;type:int(11);not null;default:1" json:"title_status"`
	/*
		试算公式状态：
		1：已配置
		-1：未配置
	*/
	FormulaStatus int32 `gorm:"column:formula_status;type:int(11);not null;default:1" json:"formula_status"`
	CreateAt      int64 `gorm:"column:create_at;type:bigint(20);not null" json:"create_at"`
	UpdateAt      int64 `gorm:"column:update_at;type:bigint(20);not null" json:"update_at"`
}

// TableName StandardStatement's table name
func (*StandardStatement) TableName() string {
	return TableNameStandardStatement
}