package code

import (
	"financial_statement/pkg/errors"
	"fmt"
	"net/http"
	"sync"
)

var codes = map[int]Coder{}
var codeMux = &sync.Mutex{}

func register(code int, httpStatus int, message string) {
	coder := &ErrCode{
		C:          code,
		HttpStatus: httpStatus,
		Message:    message,
	}

	MustRegister(coder)
}

// MustRegister 注册错误码
func MustRegister(coder Coder) {
	if coder.Code() == 0 {
		panic("ErrUnknown error code")
	}

	codeMux.Lock()
	defer codeMux.Unlock()

	if _, ok := codes[coder.Code()]; ok {
		panic(fmt.Sprintf("code: %d already exist", coder.Code()))
	}

	codes[coder.Code()] = coder
}

// ParseCoder 解析错误状态码
func ParseCoder(err error) Coder {
	if err == nil {
		return nil
	}

	if v, ok := err.(*errors.WithCode); ok {
		if coder, ok := codes[v.Code]; ok {
			msg := v.Error() // 如果有错误信息，则返回错误信息
			if msg != "" {
				return ErrCode{
					C:          coder.Code(),
					HttpStatus: coder.HTTPStatus(),
					Message:    msg,
				}
			}
			return coder
		}
	}

	return ErrCode{
		C:          http.StatusBadRequest,
		HttpStatus: http.StatusBadRequest,
		Message:    err.Error(),
	}
}

// Coder 返回状态码接口定义
type Coder interface {
	Code() int
	HTTPStatus() int
	String() string
}

type ErrCode struct {
	C          int
	HttpStatus int
	Message    string
}

func (coder ErrCode) Code() int {
	return coder.C
}

func (coder ErrCode) String() string {
	return coder.Message
}

func (coder ErrCode) HTTPStatus() int {
	if coder.HttpStatus == 0 {
		return http.StatusInternalServerError
	}

	return coder.HttpStatus
}
