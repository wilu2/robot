// Code generated by goctl. DO NOT EDIT.
package formulas

import (
	"financial_statement/internal/apiserver/handler/standard/statement/formulas"
	"financial_statement/internal/apiserver/middleware"
	"financial_statement/internal/apiserver/svc"

	"github.com/gin-gonic/gin"
)

func RegisterStandard_statement_formulasRoute(e *gin.Engine, svcCtx *svc.ServiceContext) {
	g := e.Group("/v2/standard/statement")
	g.Use(middleware.AuthorizationMiddleware)
	g.POST("/:statement_id/formula", formulas.FormulaCreateHandle(svcCtx))
	g.GET("/:statement_id/formulas", formulas.FormulaListHandle(svcCtx))
	g.DELETE("/:statement_id/formula/:formula_id", formulas.DisableFormulaHandle(svcCtx))
	g.PUT("/:statement_id/formula/:formula_id", formulas.UpdateFormulaHandle(svcCtx))

}
