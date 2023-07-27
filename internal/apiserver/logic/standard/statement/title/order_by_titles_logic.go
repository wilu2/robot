package title

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/title"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type OrderByTitlesLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewOrderByTitlesLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) OrderByTitlesLogic {
	return OrderByTitlesLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// OrderByTitles 科目排序
func (l *OrderByTitlesLogic) OrderByTitles(req *title.OrderByTitlesReq) (err error) {
	var (
		tTitles = query.Use(l.svcCtx.Db).StandardStatementTitle
		q       = query.Use(l.svcCtx.Db)
	)
	err = q.Transaction(func(tx *query.Query) error {
		for index, id := range req.TitleIdList {
			if _, err = q.StandardStatementTitle.WithContext(l.ctx).Where(tTitles.StatementID.Eq(uint32(req.StatementId)), tTitles.ID.Eq(uint32(id))).Updates(model.StandardStatementTitle{
				OrderByID: int32(index + 1),
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
