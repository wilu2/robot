package setting

import (
	"financial_statement/internal/apiserver/logic/setting"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

// SettingInfoHandle 查看配置
func SettingInfoHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		logic := setting.NewSettingInfoLogic(c, svcCtx)
		resp, err := logic.SettingInfo()
		response.HandleResponse(c, resp, err)
	}
}
