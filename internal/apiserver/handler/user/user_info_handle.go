package user

import (
	"financial_statement/internal/apiserver/logic/user"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	userType "financial_statement/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// UserInfoHandle 查看用户
func UserInfoHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.GetUserReq
		if err := c.ShouldBindUri(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewUserInfoLogic(c, svcCtx)
		resp, err := logic.UserInfo(&req)
		response.HandleResponse(c, resp, err)
	}
}
