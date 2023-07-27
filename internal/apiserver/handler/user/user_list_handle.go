package user

import (
	"financial_statement/internal/apiserver/logic/user"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	userType "financial_statement/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// UserListHandle 列出用户
func UserListHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.UserListReq
		if err := c.ShouldBindQuery(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewUserListLogic(c, svcCtx)
		resp, err := logic.UserList(&req)
		response.HandleResponse(c, resp, err)
	}
}
