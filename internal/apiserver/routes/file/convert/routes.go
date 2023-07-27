// Code generated by goctl. DO NOT EDIT.
package convert

import (
	"financial_statement/internal/apiserver/handler/file/convert"
	"financial_statement/internal/apiserver/middleware"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterFile_convertRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/v2/file_convert")
	g.Use(middleware.AuthorizationMiddleware)
	g.POST("/pdf2image", convert.Pdf2ImageHandle(svcCtx))

}
