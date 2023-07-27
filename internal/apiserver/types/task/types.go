// Code generated by goctl. DO NOT EDIT.
package task

import "mime/multipart"

type CreateTaskReq struct {
	StandardID         uint32               `json:"standard_id" binding:"required"`                     //准则id
	ExternalStandardID string               `json:"external_standard_id"`                               //外部准则id
	Async              int                  `json:"async" default:"0" binding:"oneof=0 1"`              //是否开启同步
	Files              []CreateTaskFileItem `json:"files" binding:"required"`                           //财报文件
	FileType           string               `json:"file_type" binding:"required,oneof=image pdf excel"` //web端选择pdf时，请先调用pdf转image接口，让用户选择识别哪些页后，再调用该接口，将用户选择的image传给后端，并且类型指定为image。
	Name               string               `json:"name" binding:"required"`                            //财报任务名称
	ExtraInfo          string               `json:"extra_info"`                                         //扩展信息（任务识别结果回发时会带上）
	ConsumerId         string               `json:"consumer_id"`                                        //附属信息：客户渠道id
}

type CreateTaskFileItem struct {
	Reference  string `json:"reference"`   //文件外部引用 reference和base64二选一。都提供的话，base64优先级更高。
	Base64     string `json:"base64"`      //文件base64
	FileNumber string `json:"file_number"` //外部的file number信息
}

type CreateTaskResp struct {
	TaskID   uint32                    `json:"task_id"` //财报id
	FileType string                    `json:"file_type" binding:"oneof=image pdf excel"`
	Pages    []CreateTaskPagesItmeResp `json:"pages"`
	FileInfo []string                  `json:"file_info"` //对应请求时的外部文件的file number
}

type TaskReIdentifyReq struct {
	TaskID     uint32 `json:"task_id"`                        //财报id
	StandardID uint32 `json:"standard_id" binding:"required"` //准则id
}

type TaskReNameReq struct {
	TaskID   uint32 `uri:"id"`         //财报id
	TaskName string `json:"task_name"` //财报名称
}

type CreateTaskPagesItmeResp struct {
	Id  uint32 `json:"id"`
	Url string `json:"url"`
}

type EditTaskReq struct {
	TaskID             uint32                   `uri:"id"`
	CancellationFlag   bool                     `form:"cancellation_flag"` // 作废标志，设置为true时为作废
	CurrentIndex       int                      `form:"current_index"`     //提交时，当前编辑的财报索引，下标从0开始
	FinancialStatement []TaskFinancialStatement `json:"financial_statement"`
}

type GetTaskReq struct {
	TaskID uint32 `uri:"id"`
}

type GetTaskResp struct {
	TaskId             uint32                   `json:"task_id"`
	Status             int                      `json:"status"`
	FinancialStatement []TaskFinancialStatement `json:"financial_statement"` //财报识别结果
	Standard           Standard                 `json:"standard"`            //准则信息
	FileFormat         string                   `json:"file_format"`         //文件类型
	CreatedAt          int                      `json:"create_at"`           //创建时间
	UpdateAt           int                      `json:"update_at"`           //更新时间
	External           []string                 `json:"external"`            //文件扩展信息
	Images             []Image                  `json:"images"`              //所有图片页
	Groups             []TaskGroupsResp         `json:"groups"`              //分组信息 //如果为空则没有分组或者没有调整过分组
	ExtraInfo          string                   `json:"extra_info"`          //扩展信息，创建任务时的扩展信息
	TaskName           string                   `json:"task_name"`           //任务名
	OperType           string                   `json:"operType,omitempty"`  // 第三方业务字段（兴业），始终返回B
}

type TaskGroupsResp struct {
	GroupName string              `json:"group_name"`
	GroupId   string              `json:"group_id"`
	Files     []TaskGroupPageInfo `json:"files"`
}

type Standard struct {
	Id         uint32 `json:"id"` //准则id
	ExternalId string `json:"external_id"`
}

type TaskFinancialStatement struct {
	BalanceSheet      Statement `json:"balance_sheet"`       //资产负债表
	IncomeStatement   Statement `json:"income_statement"`    //利润表
	CashFlowStatement Statement `json:"cash_flow_statement"` //现金流量表
}

type Statement struct {
	Count   int              `json:"count"`   //总列数
	Images  []Image          `json:"images"`  //图片信息(供前端展示用)
	Headers []Header         `json:"headers"` //列头
	Titles  []StatementTitle `json:"titles"`  //科目
}

type StatementTitle struct {
	Id         uint32                `json:"id"`          //匹配到的标准科目ID
	TitleName  string                `json:"title_name"`  //匹配到的标准科目名称
	Similarity float64               `json:"similarity"`  //匹配置信度 0到1之间，大于等于0.85匹配正常，小于0.85大于0匹配度低，0不匹配
	ExternalId string                `json:"external_id"` //外部科目ID
	PageIndex  int                   `json:"page_index"`  //科目所在文档的index
	RowIndex   int                   `json:"-"`
	ColIndex   int                   `json:"-"`
	Key        StatementTitleKey     `json:"key"` //科目信息
	Values     []StatementTitleValue `json:"values"`
	ModifyFlag bool                  `json:"modify_flag"` //修改标识
}

type StatementTitleKey struct {
	Ocr      string `json:"ocr"`      //ocr结果
	Inferred uint32 `json:"inferred"` //推断出来的科目id
	Position []int  `json:"position"` //坐标
}

type StatementTitleValue struct {
	Ocr        string `json:"ocr"`        //ocr识别结果
	Supervised string `json:"supervised"` //人工编辑后的值
	Position   []int  `json:"position"`   //坐标
}

type Header struct {
	Organization string `json:"organization"` //所处组织
	Date         string `json:"date"`         //期限
	Order        int    `json:"order"`        // 系统推理出来的时间维度排序，越大越靠后
}

type Image struct {
	ImageSrc    string `json:"img_src"`
	RotateAngle int    `json:"rotate_angle"`
	FileId      uint32 `json:"file_id"`
}

type ListTaskReq struct {
	Page        uint32 `form:"page"`
	PerPage     uint32 `form:"per_page"`
	KeyWords    string `form:"keywords"`
	OrderBy     string `form:"order_by"`
	OrderByType string `form:"order_by_type"`
}

type TaskGroupReq struct {
	Group  []TaskGroupInfo `json:"group"`
	TaskID uint32          `uri:"id"`
}

type TaskGroupInfo struct {
	GroupName string              `json:"group_name"`
	GroupId   string              `json:"group_id"`
	Files     []TaskGroupPageInfo `json:"files"`
}

type TaskGroupPageInfo struct {
	FileId      uint32 `json:"file_id"`
	ImageSrc    string `json:"img_src"`
	RotateAngle int    `json:"rotate_angle"`
	Type        int    `json:"type"` //1资产负债表  2现金流量表  3利润表
}

type TaskGroupPageReq struct {
	TaskID uint32                `uri:"id"`
	File   *multipart.FileHeader `form:"file"` // 文件流
}

type TaskGroupPageResp struct {
	FileId   uint32 `json:"file_id"`
	ImageSrc string `json:"img_src"`
}

type ListTaskResp struct {
	Tasks []TaskListInfo `json:"tasks"`
	Total int64          `json:"total"`
}

type TaskListInfo struct {
	ID             uint32  `json:"id"`
	FileFormat     string  `json:"file_format"`
	TaskName       string  `json:"task_name"`
	StatementTypes [][]int `json:"statement_types"`
	StandardID     uint32  `json:"standard_id"`
	Async          int     `json:"async"`
	Status         int     `json:"status"`
	Error          string  `json:"error"`
	CreatedAt      int     `json:"created_at"`
	UpdatedAt      int     `json:"updated_at"`
}