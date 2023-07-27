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

type StandardDeleteLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardDeleteLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardDeleteLogic {
	return StandardDeleteLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// StandardDelete 删除准则
func (l *StandardDeleteLogic) StandardDelete(req *standard.GetStandardReq) (err error) {
	var (
		tStandards = query.Use(l.svcCtx.Db).Standard
	)

	_, err = tStandards.WithContext(l.ctx).Where(tStandards.ID.Eq(req.ID)).Update(tStandards.Status, req.Status)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
