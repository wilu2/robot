package user

import (
	"financial_statement/internal/apiserver/logic/user"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	userType "financial_statement/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// UserUpdateHandle 修改用户信息
func UserUpdateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.UpdateUserReq
		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewUserUpdateLogic(c, svcCtx)
		err := logic.UserUpdate(&req)
		response.HandleResponse(c, nil, err)
	}
}
