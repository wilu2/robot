package formulas

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/formulas"
	"financial_statement/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
)

type DisableFormulaLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDisableFormulaLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DisableFormulaLogic {
	return DisableFormulaLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DisableFormula 删除或禁用公式
func (l *DisableFormulaLogic) DisableFormula(req *formulas.DisableFormulaReq) (err error) {
	var (
		tFormula = query.Use(l.svcCtx.Db).StandardStatementFormula
	)
	_, err = tFormula.WithContext(l.ctx).Where(tFormula.ID.Eq(uint32(req.FormulaId)), tFormula.StatementID.Eq(uint32(req.StatementId))).
		Updates(model.StandardStatementFormula{
			Status:   int32(req.Status),
			UpdateAt: time.Now().Unix()})
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	if req.Status == consts.TitleStatusDelete {
		_, err = tFormula.WithContext(l.ctx).Where(tFormula.ID.Eq(uint32(req.FormulaId)), tFormula.StatementID.Eq(uint32(req.StatementId))).Unscoped().Delete()
		if err != nil {
			err = errors.WithCodeMsg(code.Internal)
			return
		}
	}

	return
}
