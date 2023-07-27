package standard

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/standard"
	"financial_statement/pkg/errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StandardCopyWithFormulasLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewStandardCopyWithFormulasLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) StandardCopyWithFormulasLogic {
	return StandardCopyWithFormulasLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

func formatFormula(formula string, oldTitles map[string]string, newTitles map[string]uint32) string {
	charList := strings.Split(formula, ",")
	for i := 0; i < len(charList); i++ {
		value, ok := oldTitles[charList[i]]
		if ok {
			charList[i] = strconv.Itoa(int(newTitles[value]))
		}

	}
	return strings.Join(charList, ",")
}

// StandardCopyWithFormulas 复制准则(连同公式)
func (l *StandardCopyWithFormulasLogic) StandardCopyWithFormulas(req *standard.CopyStandardReq) (err error) {
	var (
		tStandards              = query.Use(l.svcCtx.Db).Standard
		tStandardStatement      = query.Use(l.svcCtx.Db).StandardStatement
		tStandardStatementTitle = query.Use(l.svcCtx.Db).StandardStatementTitle
		tFormula                = query.Use(l.svcCtx.Db).StandardStatementFormula
		standardStatements      []*model.StandardStatement
		standard                *model.Standard
		q                       = query.Use(l.svcCtx.Db)
	)

	if count, _ := tStandards.WithContext(l.ctx).Where(tStandards.Name.Eq(req.Name)).Count(); count > 0 {
		err = errors.WithCodeMsg(code.BadRequest, "已存在相同名称的数据")
		return
	}

	// 准则对象
	if standard, err = tStandards.WithContext(l.ctx).Where(tStandards.ID.Eq(req.ID)).First(); err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	// 三大财务准则表对象
	if standardStatements, err = tStandardStatement.WithContext(l.ctx).Where(tStandardStatement.StandardID.Eq(req.ID)).Find(); err != nil {
		err = errors.WithCodeMsg(code.Internal)
		return
	}

	if err = q.Transaction(func(tx *query.Query) error {
		newStandard := model.Standard{
			Name:       req.Name,
			ExternalID: standard.ExternalID,
			IsDefault:  consts.StandardNotDefault,
			Status:     consts.StandardStatusNormal,
		}
		tStandards.WithContext(l.ctx).Create(&newStandard)

		//创建三张财务报表
		for _, t := range standardStatements {
			newStandardStatement := model.StandardStatement{
				StandardID:    newStandard.ID,
				Type:          t.Type,
				Status:        t.Status,
				TitleStatus:   t.TitleStatus,
				FormulaStatus: t.FormulaStatus,
				CreateAt:      time.Now().Unix(),
				UpdateAt:      time.Now().Unix(),
			}
			if err := tx.StandardStatement.WithContext(l.ctx).Create(&newStandardStatement); err != nil {
				return err
			}

			oldTitleMaps := make(map[string]string)
			newTitleMaps := make(map[string]uint32)
			var titles []*model.StandardStatementTitle
			if titles, err = tx.StandardStatementTitle.WithContext(l.ctx).Where(tStandardStatementTitle.StatementID.Eq(t.ID)).Find(); err != nil {
				return err
			}
			for _, title := range titles {
				oldTitleMaps[strconv.FormatUint(uint64(title.ID), 10)] = title.Name
				t := &model.StandardStatementTitle{
					StatementID: newStandardStatement.ID,
					Name:        title.Name,
					ExternalID:  title.ExternalID,
					Aliases:     title.Aliases,
					Status:      title.Status,
					OrderByID:   title.OrderByID,
					CreateAt:    time.Now().Unix(),
					UpdateAt:    time.Now().Unix(),
				}
				if err := tx.StandardStatementTitle.WithContext(l.ctx).Create(t); err != nil {
					return err
				}
				newTitleMaps[title.Name] = t.ID
			}

			var formulas []*model.StandardStatementFormula
			if formulas, err = tx.StandardStatementFormula.WithContext(l.ctx).Where(tFormula.StatementID.Eq(t.ID)).Find(); err != nil {
				return err
			}
			for _, formula := range formulas {
				newFormula := &model.StandardStatementFormula{
					StatementID: newStandardStatement.ID,
					Left:        formatFormula(formula.Left, oldTitleMaps, newTitleMaps),
					Right:       formatFormula(formula.Right, oldTitleMaps, newTitleMaps),
					Status:      formula.Status,
					CreateAt:    time.Now().Unix(),
					UpdateAt:    time.Now().Unix(),
				}
				if err := tx.StandardStatementFormula.WithContext(l.ctx).Create(newFormula); err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}

	return
}
