// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameSetting = "settings"

// Setting mapped from table <settings>
type Setting struct {
	Key     string `gorm:"column:key;type:varchar(128);not null" json:"key"`       // 配置key
	Setting int32  `gorm:"column:setting;type:int(11);not null" json:"setting"`    // 配置值
	Remark  string `gorm:"column:remark;type:varchar(128);not null" json:"remark"` // 配置说明
}

// TableName Setting's table name
func (*Setting) TableName() string {
	return TableNameSetting
}
