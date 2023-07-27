package setting

import (
	"financial_statement/internal/apiserver/logic/setting"
	"financial_statement/internal/apiserver/response"
	"financial_statement/internal/apiserver/svc"
	settingType "financial_statement/internal/apiserver/types/setting"

	"github.com/gin-gonic/gin"
)

// SettingUpdateHandle 修改配置信息
func SettingUpdateHandle(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		var req settingType.UpdateSettingReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.HandlerParamsResponse(c, err)
			return
		}

		logic := setting.NewSettingUpdateLogic(c, svcCtx)
		err := logic.SettingUpdate(&req)
		response.HandleResponse(c, nil, err)
	}
}
