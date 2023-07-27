// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameStandard = "standards"

// Standard mapped from table <standards>
type Standard struct {
	ID         uint32 `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	Name       string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	ExternalID string `gorm:"column:external_id;type:varchar(255);not null" json:"external_id"`
	IsDefault  int32  `gorm:"column:is_default;type:int(11);not null" json:"is_default"` // 1为默认 -1为非默认
	/*
		状态：
		1：正常
		-1：停用
	*/
	Status int32 `gorm:"column:status;type:int(11);not null;default:1" json:"status"`
}

// TableName Standard's table name
func (*Standard) TableName() string {
	return TableNameStandard
}
