package unauthorization

import (
	"financial_statement/internal/apiserver/logic/user/unauthorization"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

// UserLogoutHandle 用户登出
func UserLogoutHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		logic := unauthorization.NewUserLogoutLogic(c, svcCtx)
		err := logic.UserLogout()
		response.HandleResponse(c, nil, err)
	}
}
