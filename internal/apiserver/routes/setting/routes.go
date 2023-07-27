// Code generated by goctl. DO NOT EDIT.
package setting

import (
	"financial_statement/internal/apiserver/handler/setting"
	"financial_statement/internal/apiserver/middleware"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterSettingRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/v2/setting")
	g.Use(middleware.AuthorizationMiddleware)
	g.GET("/", setting.SettingInfoHandle(svcCtx))
	g.POST("/update", setting.SettingUpdateHandle(svcCtx))

}
