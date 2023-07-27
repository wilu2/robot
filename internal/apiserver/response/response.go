package response

import (
	"financial_statement/internal/apiserver/code"
	"financial_statement/pkg/validation"

	"net/http"

	"github.com/gin-gonic/gin"
)

// UnifiedResponse 统一返回
type UnifiedResponse struct {
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"msg"`
}

// BadResponse 错误返回
type BadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

// HandlerParamsResponse 处理参数响应错误
func HandlerParamsResponse(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, BadResponse{
		Code:    http.StatusBadRequest,
		Message: validation.Error(err),
	})
}

func HandleResponseWithStatusOk(c *gin.Context, data any, err error) {
	if err == nil {
		c.JSON(http.StatusOK, UnifiedResponse{
			Code:    http.StatusOK,
			Data:    data,
			Message: "success",
		})
	} else {
		coder := code.ParseCoder(err)
		c.JSON(http.StatusOK, BadResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		})
	}
}

// HandleResponse 统一返回处理 {"msg":"","data":{},"code":200}
func HandleResponse(c *gin.Context, data any, err error) {
	if c.Writer.Status() != 200 { // 手动指定 http Code 如 302 重定向
		return
	}
	if err == nil {
		c.JSON(http.StatusOK, UnifiedResponse{
			Code:    http.StatusOK,
			Data:    data,
			Message: "success",
		})
	} else {
		coder := code.ParseCoder(err)
		c.JSON(coder.HTTPStatus(), BadResponse{
			Code:    coder.Code(),
			Message: coder.String(),
		})
	}
}
