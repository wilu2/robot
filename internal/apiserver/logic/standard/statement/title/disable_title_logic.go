package title

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard/statement/title"
	"financial_statement/pkg/errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type DisableTitleLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewDisableTitleLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) DisableTitleLogic {
	return DisableTitleLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// DisableTitle 禁用或删除科目
func (l *DisableTitleLogic) DisableTitle(req *title.DisableTitleReq) (err error) {
	var (
		tTitles = query.Use(l.svcCtx.Db).StandardStatementTitle
		// tFormulaTitleMaps = query.Use(l.svcCtx.Db).FormulaTitleMap
		tFormulas = query.Use(l.svcCtx.Db).StandardStatementFormula
		q         = query.Use(l.svcCtx.Db)
	)

	formulas, err := tFormulas.WithContext(l.ctx).Where(tFormulas.StatementID.Eq(req.StatementId)).Find()

	for _, item := range formulas {
		charList := strings.Split(item.Left, ",")
		for _, char := range charList {
			if id, _ := strconv.Atoi(char); err == nil {
				if id == int(req.TitleId) {
					err = errors.WithCodeMsg(code.BadRequest, "该科目已被用在准则公式里，请删除所有使用该科目的公式后再删除该科目。")
					return
				}
			}
		}
		charList = strings.Split(item.Right, ",")
		for _, char := range charList {
			if id, _ := strconv.Atoi(char); err == nil {
				if id == int(req.TitleId) {
					err = errors.WithCodeMsg(code.BadRequest, "该科目已被用在准则公式里，请删除所有使用该科目的公式后再删除该科目。")
					return
				}
			}
		}
	}

	err = q.Transaction(func(tx *query.Query) error {

		if _, err = q.StandardStatementTitle.WithContext(l.ctx).Where(tTitles.ID.Eq(uint32(req.TitleId)), tTitles.StatementID.Eq(uint32(req.StatementId))).Updates(model.StandardStatementTitle{
			Status: int32(req.Status),
		}); err != nil {
			return err
		}

		//禁用科目时，将所有关联的试算平衡公式也一起禁用掉
		// if req.Status == consts.TitleStatusDisabled {
		// 	var FormulaIDs []uint32
		// 	if err = q.FormulaTitleMap.WithContext(l.ctx).Select(tFormulaTitleMaps.FormulaID).Where(tFormulaTitleMaps.TitleID.Eq(uint32(req.TitleId))).Scan(&FormulaIDs); err != nil {
		// 		return err
		// 	}
		// 	if _, err = q.StandardStatementFormula.WithContext(l.ctx).Where(tFormulas.ID.In(FormulaIDs...)).Updates(model.StandardStatementFormula{
		// 		Status: consts.FormulaStatusDisabled,
		// 	}); err != nil {
		// 		return err
		// 	}
		// }

		if req.Status == consts.TitleStatusDelete {
			if _, err = q.StandardStatementTitle.WithContext(l.ctx).Where(tTitles.ID.Eq(uint32(req.TitleId)), tTitles.StatementID.Eq(uint32(req.StatementId))).Unscoped().Delete(); err != nil {
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
