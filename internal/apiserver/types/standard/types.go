// Code generated by goctl. DO NOT EDIT.
package standard

type CreateStandardReq struct {
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
}

type CopyStandardReq struct {
	ID   uint32 `uri:"id"`
	Name string `json:"name"`
}

type CreateStandardResp struct {
	StandardID uint32 `json:"standard_id"`
}

type UpdateStandardReq struct {
	ID         uint32 `uri:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
	Status     int    `json:"status"`
}

type GetStandardReq struct {
	ID     uint32 `uri:"id"`
	Status int    `json:"status"`
}

type ListStandardReq struct {
	Page    uint32 `form:"page"`
	PerPage uint32 `form:"per_page"`
}

type ListStandardResp struct {
	Standards []StandardListInfo `json:"standards"`
	Total     int64              `json:"total"`
}

type StandardListInfo struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	ExternalID string `json:"external_id"`
	IsDefault  int    `json:"is_default"`
	Status     int    `json:"status"`
}

type GetStatementsReq struct {
	ID uint32 `uri:"id"`
}

type GetStatementsResp struct {
	Statements []StatementItems `json:"statements"`
}

type StatementItems struct {
	Type       int       `json:"type"`
	StandardID int       `json:"standard_id"`
	Status     int       `json:"status"`
	ID         uint32    `json:"id"`
	Titles     []Title   `json:"titles"`
	Formulas   []Formula `json:"formulas"`
}

type SetDefaultStandardReq struct {
	ID uint32 `uri:"id"`
}

type Title struct {
	ID         uint32 `json:"id"`
	Name       string `json:"name"`
	ExternalId string `json:"external_id"`
	Aliases    string `json:"aliases"` //别名，json字符串，示例：["name1","name2"]
	Status     int    `json:"status"`
}

type Formula struct {
	ID          uint32 `json:"id"`
	StatementId uint32 `json:"statement_id"`
	Left        string `json:"left" binding:"max=2048"`  //公式左边
	Right       string `json:"right" binding:"max=2048"` //公式右边
	Status      int    `json:"status"`
}
