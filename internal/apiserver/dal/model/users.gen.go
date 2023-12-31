// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameUser = "users"

// User mapped from table <users>
type User struct {
	ID                uint32 `gorm:"column:id;type:int(11) unsigned;primaryKey;autoIncrement:true" json:"id"`
	Name              string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Account           string `gorm:"column:account;type:varchar(255);not null" json:"account"`
	Password          string `gorm:"column:password;type:varchar(64);not null" json:"password"`
	Salt              string `gorm:"column:salt;type:varchar(32);not null" json:"salt"`
	Email             string `gorm:"column:email;type:varchar(128);not null" json:"email"`
	Mobile            string `gorm:"column:mobile;type:varchar(32);not null" json:"mobile"`
	AuthMethod        string `gorm:"column:auth_method;type:varchar(128);not null" json:"auth_method"`
	ExpiryTime        int64  `gorm:"column:expiry_time;type:bigint(20);not null" json:"expiry_time"`                   // 过期时间
	FailedLogins      int32  `gorm:"column:failed_logins;type:int(11);not null" json:"failed_logins"`                  // 连续登录失败次数
	LastLoginFailTime int64  `gorm:"column:last_login_fail_time;type:bigint(20);not null" json:"last_login_fail_time"` // 上次登录失败时间
	LastLoginTime     int64  `gorm:"column:last_login_time;type:bigint(20);not null" json:"last_login_time"`           // 上次登录时间
	Status            int32  `gorm:"column:status;type:tinyint(4);not null;default:1" json:"status"`                   // 用户状态 0未使用，1正常，2锁定，3禁用
	CreatedAt         int64  `gorm:"column:created_at;type:bigint(20);not null" json:"created_at"`
	UpdatedAt         int64  `gorm:"column:updated_at;type:bigint(20);not null" json:"updated_at"`
	DeletedAt         *int64 `gorm:"column:deleted_at;type:bigint(20)" json:"deleted_at"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
