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

type StandardInfoLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardInfoLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardInfoLogic {
	return StandardInfoLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// StandardInfo 获取准则信息
func (l *StandardInfoLogic) StandardInfo(req *standard.GetStandardReq) (resp standard.UpdateStandardReq, err error) {
	var (
		tStandards = query.Use(l.svcCtx.Db).Standard
	)
	standard, err := tStandards.WithContext(l.ctx).Where(tStandards.ID.Eq(uint32(req.ID))).First()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	copier.Copy(&resp, &standard)
	return
}
