package title

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/title"
	"financial_statement/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TitleListLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTitleListLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TitleListLogic {
	return TitleListLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TitleList 获取所有科目
func (l *TitleListLogic) TitleList(req *title.TitleListReq) (resp title.TitleListResp, err error) {
	var (
		tTitles = query.Use(l.svcCtx.Db).StandardStatementTitle
	)
	query := tTitles.WithContext(l.ctx).Where(tTitles.StatementID.Eq(uint32(req.StatementId))).Order(tTitles.OrderByID, tTitles.ID.Desc())
	if count, e := query.Count(); e != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	} else {
		resp.Total = count
	}
	titles, err := query.Offset((int(req.Page) - 1) * int(req.PerPage)).Limit(int(req.PerPage)).Find()
	for _, item := range titles {
		resp.Titles = append(resp.Titles, title.Title{
			ID:         item.ID,
			Name:       item.Name,
			ExternalId: *item.ExternalID,
			Aliases:    *item.Aliases,
			Status:     int(item.Status),
		})
	}
	return
}
