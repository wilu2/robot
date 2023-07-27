package title

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/title"
	"financial_statement/pkg/errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type TitleCreateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTitleCreateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TitleCreateLogic {
	return TitleCreateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// TitleCreate 新增科目
func (l *TitleCreateLogic) TitleCreate(req *title.CreateTitleReq) (err error) {
	var (
		q      = query.Use(l.svcCtx.Db)
		tTitle = query.Use(l.svcCtx.Db).StandardStatementTitle
	)

	if err = q.Transaction(func(tx *query.Query) error {
		for _, title := range req.Titles {
			count, err := tx.StandardStatementTitle.WithContext(l.ctx).Where(tTitle.StatementID.Eq(req.StatementId),
				tTitle.Name.Eq(title.Name)).Count()
			if err != nil {
				return err
			}
			if count > 0 {
				return errors.WithCodeMsg(code.BadRequest, fmt.Sprintf("已存在名为 %s 的数据", title.Name))
			}
			if _, err := tx.StandardStatementTitle.WithContext(l.ctx).Attrs(
				tTitle.StatementID.Value(req.StatementId),
				tTitle.Name.Value(title.Name),
				tTitle.ExternalID.Value(title.ExternalId),
				tTitle.Aliases.Value(title.Aliases),
				tTitle.Status.Value(consts.TitleStatusNormal),
				tTitle.OrderByID.Value(9999),
				tTitle.CreateAt.Value(time.Now().Unix()),
				tTitle.UpdateAt.Value(time.Now().Unix()),
			).Where(tTitle.StatementID.Eq(req.StatementId),
				tTitle.Name.Eq(title.Name)).FirstOrCreate(); err != nil {
				return err
			}

		}
		return nil
	}); err != nil {
		err = errors.WithCodeMsg(code.BadRequest, err.Error())
		return
	}
	return
}
