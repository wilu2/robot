package unauthorization

import (
	"financial_statement/internal/apiserver/logic/user/unauthorization"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	unauthorizationType "financial_statement/internal/apiserver/types/user/unauthorization"

	"github.com/gin-gonic/gin"
)

// UserUpdatePwdHandle 不使用token修改用户密码，需要旧密码
func UserUpdatePwdHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req unauthorizationType.UpdatePwdReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := unauthorization.NewUserUpdatePwdLogic(c, svcCtx)
		err := logic.UserUpdatePwd(&req)
		response.HandleResponse(c, nil, err)
	}
}
