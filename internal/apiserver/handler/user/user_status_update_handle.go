package user

import (
	"financial_statement/internal/apiserver/logic/user"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	userType "financial_statement/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// UserStatusUpdateHandle 修改用户状态，有效期
func UserStatusUpdateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.UpdateUserStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewUserStatusUpdateLogic(c, svcCtx)
		err := logic.UserStatusUpdate(&req)
		response.HandleResponse(c, nil, err)
	}
}
