package standard

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type StandardUpdateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardUpdateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardUpdateLogic {
	return StandardUpdateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// StandardUpdate 更新准则
func (l *StandardUpdateLogic) StandardUpdate(req *standard.UpdateStandardReq) (err error) {
	var (
		tStandards = query.Use(l.svcCtx.Db).Standard
	)
	if count, _ := tStandards.WithContext(l.ctx).Where(tStandards.Name.Eq(req.Name)).Where(tStandards.ID.Neq(req.ID)).Count(); count > 0 {
		err = errors.WithCodeMsg(code.BadRequest, "已存在相同名称的数据")
		return
	}
	_, err = tStandards.WithContext(l.ctx).Where(tStandards.ID.Eq(req.ID)).Updates(model.Standard{
		ExternalID: req.ExternalID,
		Name:       req.Name,
		Status:     int32(req.Status),
	})

	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
