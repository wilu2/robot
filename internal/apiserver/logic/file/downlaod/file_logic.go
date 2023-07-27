package downlaod

import (
	"context"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/file/downlaod"
	"financial_statement/internal/pkg/storage"
	"financial_statement/pkg/errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type FileLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewFileLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) FileLogic {
	return FileLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// File 文件下载
func (l *FileLogic) File(req *downlaod.FileReq) (resp downlaod.FileResp, err error) {
	if l.svcCtx.Config.ServerOptions.IsSaas {
		var (
			u = l.ginCtx.Keys["user"].(*model.User)
		)
		if u.ID != uint32(req.Uid) {
			err = errors.WithCodeMsg(code.Unauthorized)
			return
		}
	}

	filePath := fmt.Sprintf("/%d/%d/%s/%s", req.Uid, req.TaskId, req.Date, req.FileName)
	fileByte, err := storage.FileStorage.Get(filePath)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	resp.File = fileByte
	return
}
