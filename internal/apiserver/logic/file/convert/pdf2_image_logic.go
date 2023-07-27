package convert

import (
	"context"
	"encoding/base64"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/file/convert"
	"financial_statement/pkg/errors"
	pdf "financial_statement/pkg/pdf"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type Pdf2ImageLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewPdf2ImageLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) Pdf2ImageLogic {
	return Pdf2ImageLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// Pdf2Image pdf转图片
func (l *Pdf2ImageLogic) Pdf2Image(req *convert.FileConvertReq) (resp convert.FileConvertResp, err error) {
	// todo: add your logic here and delete this line
	fileIo, err := req.File.Open()
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	pdfBytes, err := ioutil.ReadAll(fileIo)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	pages, err := pdf.GetPages(pdfBytes)
	if err != nil {
		err = errors.WithCodeMsg(code.BadRequest, "已损坏的文件，解析失败")
		return
	}
	if len(pages) == 0 {
		err = errors.WithCodeMsg(code.BadRequest, "文件解析失败")
		return
	}
	for _, page := range pages {
		resp.Files = append(resp.Files, base64.StdEncoding.EncodeToString(page))
	}
	return
}
