package jobtask

import "financial_statement/internal/apiserver/types/task"

const (
	RecognizeTableResultTableTypePlain            = "plain"
	RecognizeTableResultTableTypeTableWithLine    = "table_with_line"
	RecognizeTableResultTableTypeTableWithoutLine = "table_without_line"
)

const (
	RecognizeTableResultSemanticTypePeriod      = "period"      //时间维度，时间跨度
	RecognizeTableResultSemanticTypeInstitution = "institution" //组织（例如：本集团）
	RecognizeTableResultSemanticTypeAmount      = "amount"      //金额
	RecognizeTableResultSemanticTypeItem        = "item"        //科目
	RecognizeTableResultSemanticTypeOther       = "other"       //其他
)

type RecognizeTableResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Angle  int                         `json:"angle"`
		Width  int                         `json:"width"`
		Height int                         `json:"height"`
		Tables []RecognizeTableResultTable `json:"tables"`
		Lines  []struct {
			Text                string      `json:"text"`
			Score               float64     `json:"score"`
			Type                string      `json:"type"`
			Position            []int       `json:"position"`
			Angle               int         `json:"angle"`
			Direction           int         `json:"direction"`
			Handwritten         int         `json:"handwritten"`
			CharScores          []float64   `json:"char_scores"`
			CharCenters         [][]int     `json:"char_centers"`
			CharPositions       [][]int     `json:"char_positions"`
			CharCandidates      [][]string  `json:"char_candidates"`
			CharCandidatesScore [][]float64 `json:"char_candidates_score"`
			Semantic            string      `json:"semantic"`
		} `json:"lines"`
		InferAngle int `json:"infer_angle"`
	} `json:"result"`
}
type RecognizeTableResultTableCell struct {
	StartRow int    `json:"start_row"`
	StartCol int    `json:"start_col"`
	EndRow   int    `json:"end_row"`
	EndCol   int    `json:"end_col"`
	Text     string `json:"text"`
	Position []int  `json:"position"`
	Semantic string `json:"semantic"`
	Order    int    `json:"order"`
	Lines    []struct {
		Text                string      `json:"text"`
		Score               float64     `json:"score"`
		Type                string      `json:"type"`
		Position            []int       `json:"position"`
		Angle               int         `json:"angle"`
		Direction           int         `json:"direction"`
		Handwritten         int         `json:"handwritten"`
		CharScores          []float64   `json:"char_scores"`
		CharCenters         [][]int     `json:"char_centers"`
		CharPositions       [][]int     `json:"char_positions"`
		CharCandidates      [][]string  `json:"char_candidates"`
		CharCandidatesScore [][]float64 `json:"char_candidates_score"`
		Semantic            string      `json:"semantic"`
	} `json:"lines"`
}

type RecognizeTableResultTable struct {
	Type         string                          `json:"type"`
	Position     []int                           `json:"position"`
	TableRows    int                             `json:"table_rows"`
	TableCols    int                             `json:"table_cols"`
	HeightOfRows []int                           `json:"height_of_rows"`
	WidthOfCols  []int                           `json:"width_of_cols"`
	TableCells   []RecognizeTableResultTableCell `json:"table_cells"`
	Lines        []struct {
		Text                string      `json:"text"`
		Score               float64     `json:"score"`
		Type                string      `json:"type"`
		Position            []int       `json:"position"`
		Angle               int         `json:"angle"`
		Direction           int         `json:"direction"`
		Handwritten         int         `json:"handwritten"`
		CharScores          []float64   `json:"char_scores"`
		CharCenters         [][]int     `json:"char_centers"`
		CharPositions       [][]int     `json:"char_positions"`
		CharCandidates      [][]string  `json:"char_candidates"`
		CharCandidatesScore [][]float64 `json:"char_candidates_score"`
		Semantic            string      `json:"semantic"`
	} `json:"lines"`
}

//财报识别时，组织 数据的结构体定义，用于临时存储组织信息
type Institution struct {
	Text     string
	Type     int // 1为表格 2为文本（无表格信息）
	StartCol int
	EndCol   int
}

type Period struct {
	Text     string
	Type     int // 1为表格 2为文本（无表格信息）
	StartCol int
	EndCol   int
}

type StandardStatementTitle struct {
	ID      uint32
	Name    string
	Aliases string
	Type    int32
	Score   float64
}

// title 的排序方法

type TitleSlice []task.StatementTitle

func (t TitleSlice) Len() int {
	return len(t)
}

func (t TitleSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t TitleSlice) Less(i, j int) bool {
	return t[i].RowIndex < t[j].RowIndex
}
