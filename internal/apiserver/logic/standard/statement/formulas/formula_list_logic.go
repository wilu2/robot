package formulas

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/formulas"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type FormulaListLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewFormulaListLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) FormulaListLogic {
	return FormulaListLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// FormulaList 获取所有公式
func (l *FormulaListLogic) FormulaList(req *formulas.FormulaListReq) (resp formulas.FormulaListResp, err error) {
	var (
		tFormula = query.Use(l.svcCtx.Db).StandardStatementFormula
	)

	query := tFormula.WithContext(l.ctx).Where(tFormula.StatementID.Eq(uint32(req.StatementId)))
	if count, e := query.Count(); e != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	} else {
		resp.Total = count
	}
	list, err := query.Offset((int(req.Page) - 1) * int(req.PerPage)).Limit(int(req.PerPage)).Find()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	for _, item := range list {
		resp.Formulas = append(resp.Formulas, formulas.Formula{
			ID:          item.ID,
			StatementId: item.StatementID,
			Left:        item.Left,
			Right:       item.Right,
			Status:      int(item.Status),
		})
	}
	return
}
