package standard

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type StandardListLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardListLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardListLogic {
	return StandardListLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// StandardList 列出准则
func (l *StandardListLogic) StandardList(req *standard.ListStandardReq) (resp standard.ListStandardResp, err error) {
	var (
		tStandards = query.Use(l.svcCtx.Db).Standard
	)

	standardQuery := tStandards.WithContext(l.ctx).Order(tStandards.ID.Desc())
	count, err := standardQuery.Count()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	standards, err := standardQuery.Offset((int(req.Page) - 1) * int(req.PerPage)).Limit(int(req.PerPage)).Find()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	copier.Copy(&resp.Standards, &standards)
	resp.Total = count
	return
}
