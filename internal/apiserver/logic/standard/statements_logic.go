package standard

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type StatementsLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStatementsLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StatementsLogic {
	return StatementsLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// Statements 获取准则报表信息
func (l *StatementsLogic) Statements(req *standard.GetStatementsReq) (resp standard.GetStatementsResp, err error) {
	var (
		// tStandards  = query.Use(l.svcCtx.Db).Standard
		tStatements = query.Use(l.svcCtx.Db).StandardStatement
		tTitles     = query.Use(l.svcCtx.Db).StandardStatementTitle
		tFormulas   = query.Use(l.svcCtx.Db).StandardStatementFormula
	)

	statements, err := tStatements.WithContext(l.ctx).Where(tStatements.StandardID.Eq(uint32(req.ID))).Find()

	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	var _error error
	for _, item := range statements {

		statement := standard.StatementItems{
			Type:       int(item.Type),
			StandardID: int(item.StandardID),
			Status:     int(item.Status),
			ID:         item.ID,
			Titles:     []standard.Title{},
			Formulas:   []standard.Formula{},
		}

		titles, err := tTitles.WithContext(l.ctx).Where(tTitles.StatementID.Eq(item.ID)).Find()
		if err != nil {
			_error = errors.WithCodeMsg(code.Internal)
			break
		}
		for _, titleItem := range titles {
			statement.Titles = append(statement.Titles, standard.Title{
				ID:         titleItem.ID,
				Name:       titleItem.Name,
				ExternalId: *titleItem.ExternalID,
				Aliases:    *titleItem.Aliases,
				Status:     int(titleItem.Status),
			})
		}

		formulas, err := tFormulas.WithContext(l.ctx).Where(tFormulas.StatementID.Eq(item.ID)).Find()
		if err != nil {
			_error = errors.WithCodeMsg(code.Internal)
			break
		}
		for _, formulaItem := range formulas {
			statement.Formulas = append(statement.Formulas, standard.Formula{
				ID:          formulaItem.ID,
				StatementId: formulaItem.StatementID,
				Left:        formulaItem.Left,
				Right:       formulaItem.Right,
				Status:      int(formulaItem.Status),
			})
		}
		resp.Statements = append(resp.Statements, statement)
	}
	if _error != nil {
		err = _error
		return
	}

	return
}
