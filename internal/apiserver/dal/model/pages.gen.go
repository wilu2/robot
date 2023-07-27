// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNamePage = "pages"

// Page mapped from table <pages>
type Page struct {
	ID        uint32  `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	TaskID    uint32  `gorm:"column:task_id;type:int(11) unsigned;not null" json:"task_id"`
	FilePath  string  `gorm:"column:file_path;type:varchar(1024);not null" json:"file_path"` // 文件存储路径
	OcrResult *string `gorm:"column:ocr_result;type:longtext" json:"ocr_result"`             // ocr结果
	Status    int32   `gorm:"column:status;type:int(11);not null" json:"status"`             // 状态
	CreateAt  int64   `gorm:"column:create_at;type:bigint(20);not null" json:"create_at"`
	UpdateAt  int64   `gorm:"column:update_at;type:bigint(20);not null" json:"update_at"`
}

// TableName Page's table name
func (*Page) TableName() string {
	return TableNamePage
}
