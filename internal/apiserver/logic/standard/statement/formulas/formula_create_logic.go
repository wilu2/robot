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

type FormulaCreateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewFormulaCreateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) FormulaCreateLogic {
	return FormulaCreateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// FormulaCreate 新增公式
func (l *FormulaCreateLogic) FormulaCreate(req *formulas.CreateFormulaReq) (resp formulas.Id, err error) {
	var (
		tformula = query.Use(l.svcCtx.Db).StandardStatementFormula
	)

	formula := &model.StandardStatementFormula{
		StatementID: uint32(req.StatementId),
		Left:        req.Left,
		Right:       req.Right,
		Status:      consts.FormulaStatusNormal,
		CreateAt:    time.Now().Unix(),
		UpdateAt:    time.Now().Unix(),
	}
	err = tformula.WithContext(l.ctx).Create(formula)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
