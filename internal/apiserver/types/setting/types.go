// Code generated by goctl. DO NOT EDIT.
package setting

type GetSettingResp struct {
	FailedLogin  int32 `json:"failed_logins" mapstructure:",omitempty"`
	LockedTime   int32 `json:"locked_time" mapstructure:",omitempty"`
	ValidTime    int32 `json:"valid_time" mapstructure:",omitempty"`
	SessionTime  int32 `json:"session_time" mapstructure:",omitempty"`
	NotLoginTime int32 `json:"not_login_time" mapstructure:",omitempty"`
}

type UpdateSettingReq struct {
	FailedLogin  int32 `json:"failed_logins" mapstructure:",omitempty"`
	LockedTime   int32 `json:"locked_time" mapstructure:",omitempty"`
	ValidTime    int32 `json:"valid_time" mapstructure:",omitempty"`
	SessionTime  int32 `json:"session_time" mapstructure:",omitempty"`
	NotLoginTime int32 `json:"not_login_time" mapstructure:",omitempty"`
}