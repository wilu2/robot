package formulas

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/formulas"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type UpdateFormulaLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFormulaLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateFormulaLogic {
	return UpdateFormulaLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateFormula 更新公式
func (l *UpdateFormulaLogic) UpdateFormula(req *formulas.UpdateFormulaReq) (err error) {
	var (
		tFormula         = query.Use(l.svcCtx.Db).StandardStatementFormula
		tFormulaTitleMap = query.Use(l.svcCtx.Db).FormulaTitleMap
		q                = query.Use(l.svcCtx.Db)
	)
	err = q.Transaction(func(tx *query.Query) error {
		if _, err = q.StandardStatementFormula.WithContext(l.ctx).Where(tFormula.ID.Eq(uint32(req.FormulaId)), tFormula.StatementID.Eq(uint32(req.StatementId))).Updates(model.StandardStatementFormula{
			Left:  req.Left,
			Right: req.Right,
		}); err != nil {
			return err
		}
		if _, err = q.FormulaTitleMap.WithContext(l.ctx).Where(tFormulaTitleMap.FormulaID.Eq(uint32(req.FormulaId))).Delete(); err != nil {
			return err
		}
		for _, item := range req.TitleIdList {
			if err = q.FormulaTitleMap.WithContext(l.ctx).Create(&model.FormulaTitleMap{
				FormulaID: uint32(req.FormulaId),
				TitleID:   uint32(item),
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
