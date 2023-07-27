package consts

import "time"

const TokenExpiry = time.Hour * 24

// 任务类型枚举
const (
	TaskFileTypeImage = "image"
	TaskFileTypePdf   = "pdf"
	TaskFileTypeExcel = "excel"
)

// 任务状态
const (
	//已删除
	TaskStatusDeleted = -1
	//创建任务中
	TaskStatusCreateTasking = 1
	//任务创建完成
	TaskStatusCreated = 10
	//识别中
	TaskStatusOcring = 20
	//识别失败
	TaskStatusOcrFailed = 90
	//识别成功
	TaskStatusOcrSuccess = 100
	//已回传
	TaskStatusCallbacked = 110
	//作废
	TaskCancellationFlag = 120
)

const (
	TaskAsyncStatusUnsynchronized = 1
	TaskAsyncStatusSynchronized   = 10
	TaskAsyncStatusFailed         = 20
)

// Pages表status字段枚举
const (
	PageStatusCreated    = 10
	PageStatusOcring     = 20
	PageStatusOcrSuccess = 100
)

// 准则表
const (
	// status字段枚举

	//禁用
	StandardStatusDisabled = -1
	//正常
	StandardStatusNormal = 1

	// 准则表is_default字段枚举

	// 非默认
	StandardNotDefault = -1
	// 默认
	StandardIsDefault = 1
)

// 财务报表
const (
	//type枚举值

	// 资产负债表
	StatementTypeBalanceSheet = 1
	//现金流动表
	StatementTypeCashFlowStatement = 2
	//利润表
	StatementTypeIncomeStatement = 3

	// status 枚举

	//正常
	StatementStatusNormal = 1
	//停用
	StatementStatusDisabled = -1

	// 已配置
	StatementStatusConfigured = 1
	//未配置
	StatementStatusNotConfigured = -1
)

const (
	TitleStatusNormal   = 1
	TitleStatusDisabled = -1
	TitleStatusDelete   = -2
)

const (
	FormulaStatusNormal   = 1
	FormulaStatusDisabled = -1
)

// Task 表中的ExternalInfo字段结构体
type ExternalInfo struct {
	ExtraInfo         string              `json:"extra_info"`         //扩展信息
	ConsumerId        string              `json:"consumer_id"`        //客户id，非业务必须字段
	FileInfo          []string            `json:"file_info"`          //文件信息，非业务必须字段
	RecognizeDuration int64               `json:"recognize_duration"` //该任务的识别耗时
	Groups            []ExternalInfoGroup `json:"groups"`
}

type ExternalInfoGroup struct {
	GroupName string `json:"group_name"`
	GroupId   string `json:"group_id"`
	Files     []struct {
		ImageSrc    string `json:"img_src"`
		RotateAngle int    `json:"rotate_angle"`
		FileId      uint32 `json:"file_id"`
		Type        int    `json:"type"` //1资产负债表  2现金流量表  3利润表
	} `json:"files"`
}

// 存储页面的基本ocr识别元素，便于后续前端使用
type PageOcrResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Angle      int `json:"angle"`
		Width      int `json:"width"`
		Height     int `json:"height"`
		InferAngle int `json:"infer_angle"`
	} `json:"result"`
}
