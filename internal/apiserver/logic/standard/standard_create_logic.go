package standard

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard"
	"financial_statement/pkg/errors"
	"time"

	"github.com/gin-gonic/gin"
)

type StandardCreateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardCreateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardCreateLogic {
	return StandardCreateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// StandardCreate 创建准则
func (l *StandardCreateLogic) StandardCreate(req *standard.CreateStandardReq) (resp standard.CreateStandardResp, err error) {
	var (
		q              = query.Use(l.svcCtx.Db)
		tStandards     = query.Use(l.svcCtx.Db).Standard
		statementTypes = []int{
			consts.StatementTypeBalanceSheet,
			consts.StatementTypeCashFlowStatement,
			consts.StatementTypeIncomeStatement,
		}
	)
	if count, _ := tStandards.WithContext(l.ctx).Where(tStandards.Name.Eq(req.Name)).Count(); count > 0 {
		err = errors.WithCodeMsg(code.BadRequest, "已存在相同名称的数据")
		return
	}
	err = q.Transaction(func(tx *query.Query) error {
		standard := &model.Standard{
			Name:       req.Name,
			IsDefault:  consts.StandardNotDefault,
			ExternalID: req.ExternalID,
			Status:     consts.StandardStatusNormal,
		}
		if err = tx.Standard.WithContext(l.ctx).Create(standard); err != nil {
			return errors.WithCodeMsg(code.Internal)
		}

		//创建三张财务报表
		for _, t := range statementTypes {
			tx.StandardStatement.WithContext(l.ctx).Create(&model.StandardStatement{
				StandardID:    standard.ID,
				Type:          int32(t),
				Status:        consts.StatementStatusNormal,
				TitleStatus:   consts.StatementStatusNotConfigured,
				FormulaStatus: consts.StatementStatusNotConfigured,
				CreateAt:      time.Now().Unix(),
				UpdateAt:      time.Now().Unix(),
			})
		}
		resp.StandardID = standard.ID
		return nil
	})
	if err != nil {
		return
	}
	return
}
