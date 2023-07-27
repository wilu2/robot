// Code generated by goctl. DO NOT EDIT.
package unauthorization

type UserPangu struct {
	Name     string `json:"name"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type PwdLoginReq struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

type PwdLoginResp struct {
	Expiry  int64  `json:"expiry"`
	Token   string `json:"token"`
	ID      uint32 `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
	IsAdmin bool   `json:"is_admin"`
}

type SsoTicketVerificationReq struct {
	Ticket string `form:"ticket"`
}

type UpdatePwdReq struct {
	Account     string `json:"account"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
