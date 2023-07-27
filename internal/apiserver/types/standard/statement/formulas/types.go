// Code generated by goctl. DO NOT EDIT.
package formulas

type Id struct {
	ID uint32 `json:"id"`
}

type CreateFormulaReq struct {
	StatementId uint32   `uri:"statement_id"`
	Left        string   `json:"left" binding:"max=2048"`  //公式左边
	Right       string   `json:"right" binding:"max=2048"` //公式右边
	TitleIdList []uint32 `json:"title_id_list"`            //该公式所包含的科目id集合
}

type UpdateFormulaReq struct {
	StatementId uint32   `uri:"statement_id"`
	FormulaId   uint32   `uri:"formula_id"`
	Left        string   `json:"left" binding:"max=2048"`  //公式左边
	Right       string   `json:"right" binding:"max=2048"` //公式右边
	TitleIdList []uint32 `json:"title_id_list"`            //该公式所包含的科目id集合
}

type FormulaListReq struct {
	StatementId uint32 `uri:"statement_id"`
	Page        uint32 `uri:"page"`
	PerPage     uint32 `uri:"per_page"`
}

type FormulaListResp struct {
	Formulas []Formula `json:"formulas"`
	Total    int64     `json:"total"`
}

type Formula struct {
	ID          uint32 `json:"id"`
	StatementId uint32 `json:"statement_id"`
	Left        string `json:"left" binding:"max=2048"`  //公式左边
	Right       string `json:"right" binding:"max=2048"` //公式右边
	Status      int    `json:"status"`
}

type DisableFormulaReq struct {
	StatementId uint32 `uri:"statement_id"`
	FormulaId   uint32 `uri:"formula_id"`
	Status      int    `json:"status"`
}