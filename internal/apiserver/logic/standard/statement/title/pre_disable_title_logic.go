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

type PreDisableTitleLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewPreDisableTitleLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) PreDisableTitleLogic {
	return PreDisableTitleLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// PreDisableTitle 预删除
func (l *PreDisableTitleLogic) PreDisableTitle(req *title.DisableTitleReq) (resp title.PreDisableTitleResp, err error) {
	var (
		tFormulaTitleMaps = query.Use(l.svcCtx.Db).FormulaTitleMap
	)

	var formulaIds []uint32
	if err = tFormulaTitleMaps.WithContext(l.ctx).Select(tFormulaTitleMaps.FormulaID).Where(tFormulaTitleMaps.TitleID.Eq(uint32(req.TitleId))).Scan(&formulaIds); err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}
	resp.FormulaIdList = formulaIds
	return
}
