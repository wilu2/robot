package user

import (
	"financial_statement/internal/apiserver/logic/user"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	userType "financial_statement/internal/apiserver/types/user"

	"github.com/gin-gonic/gin"
)

// UserCreateHandle 新建用户
func UserCreateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req userType.CreateUserReq
		if err := c.ShouldBind(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := user.NewUserCreateLogic(c, svcCtx)
		err := logic.UserCreate(&req)
		response.HandleResponseWithStatusOk(c, nil, err)
	}
}
