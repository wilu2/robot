// Code generated by goctl. DO NOT EDIT.
package api

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

type CreateTaskPagesItmeResp struct {
	Id  uint32 `json:"id"`
	Url string `json:"url"`
}
