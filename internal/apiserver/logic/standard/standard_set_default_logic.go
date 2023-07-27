package standard

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type StandardSetDefaultLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardSetDefaultLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardSetDefaultLogic {
	return StandardSetDefaultLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// StandardSetDefault 设为默认准则
func (l *StandardSetDefaultLogic) StandardSetDefault(req *standard.SetDefaultStandardReq) (err error) {
	var (
		tStandards = query.Use(l.svcCtx.Db).Standard
		q          = query.Use(l.svcCtx.Db)
	)
	err = q.Transaction(func(tx *query.Query) error {
		_, err = q.Standard.WithContext(l.ctx).Where(tStandards.ID.Gte(0)).Update(tStandards.IsDefault, consts.StandardNotDefault)
		if err != nil {
			return errors.WithCodeMsg(code.Internal, err.Error())
		}
		_, err = q.Standard.WithContext(l.ctx).Where(tStandards.ID.Eq(uint32(req.ID))).Update(tStandards.IsDefault, consts.StandardIsDefault)
		if err != nil {
			return errors.WithCodeMsg(code.Internal, err.Error())
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}
