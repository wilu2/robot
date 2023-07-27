package title

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/title"
	"financial_statement/pkg/errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateTitleLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateTitleLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) UpdateTitleLogic {
	return UpdateTitleLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// UpdateTitle 更新科目
func (l *UpdateTitleLogic) UpdateTitle(req *title.UpdateTitleReq) (err error) {
	var (
		tTitles = query.Use(l.svcCtx.Db).StandardStatementTitle
	)
	count, err := tTitles.WithContext(l.ctx).Where(tTitles.ID.Neq(uint32(req.TitleId)), tTitles.Name.Eq(req.Name), tTitles.StatementID.Eq(uint32(req.StatementId))).Count()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	if count > 0 {
		err = errors.WithCodeMsg(code.BadRequest, fmt.Sprintf("已存在名为 %s 的数据", req.Name))
		return
	}
	_, err = tTitles.WithContext(l.ctx).Where(tTitles.ID.Eq(uint32(req.TitleId)), tTitles.StatementID.Eq(uint32(req.StatementId))).Updates(
		model.StandardStatementTitle{
			Name:       req.Name,
			ExternalID: &req.ExternalId,
			Aliases:    &req.Aliases,
			UpdateAt:   time.Now().Unix(),
		})
	if err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	return
}
