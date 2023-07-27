// Code generated by goctl. DO NOT EDIT.
package user

type CreateUserReq struct {
	Name     string `json:"name"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
}

type UpdateUserReq struct {
	ID       uint32  `json:"id"`
	Name     *string `json:"name" mapstructure:",omitempty"`
	Password *string `json:"password" mapstructure:",omitempty"`
	Email    *string `json:"email" mapstructure:",omitempty"`
	Mobile   *string `json:"mobile" mapstructure:",omitempty"`
}

type DeleteUserReq struct {
	ID uint32 `json:"id"`
}

type ListUserResp struct {
	Users []UserInfo `json:"users"`
	Count int64      `json:"count"`
}

type GetUserReq struct {
	ID uint32 `uri:"id"`
}

type GetUserResp struct {
	UserInfo
}

type UserInfo struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Status     int    `json:"status"`
	ExpiryTime int64  `json:"expiry_time"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type UserListReq struct {
	Page        uint32 `form:"page,default=1" default:"1"`
	PerPage     uint32 `form:"per_page,default=20" default:"20"`
	OrderBy     string `form:"order_by,default=create_time" default:"create_time"`
	OrderByType string `form:"order_by_type,default=desc" default:"desc"`
}

type UpdateUserStatusReq struct {
	ID         uint32 `json:"id"`
	Status     int32  `json:"status"`
	ExpiryTime int64  `json:"expiry_time"`
}