package task

import (
	"context"
	"encoding/base64"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	"financial_statement/pkg/errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TaskGroupPageLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskGroupPageLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskGroupPageLogic {
	return TaskGroupPageLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

type File struct {
	FileName string
	File     []byte
	FilePath string
}

func getFile(gin *gin.Context, key string) (*File, error) {
	multipartForm, err := gin.MultipartForm()
	if err != nil {
		return nil, err
	}
	for _, fileHeader := range multipartForm.File[key] {
		file, err := fileHeader.Open()
		defer file.Close()
		if err != nil {
			return nil, err
		}
		// 处理文件名
		filename := filepath.Base(fileHeader.Filename)
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return &File{
			FileName: filename,
			File:     fileBytes,
		}, nil
	}
	return nil, fmt.Errorf("未找到任何文件")
}

// TaskGroupPage 调整分组时的上传文件接口
func (l *TaskGroupPageLogic) TaskGroupPage(req *task.TaskGroupPageReq) (resp task.TaskGroupPageResp, err error) {
	var (
		u      = l.ginCtx.Keys["user"].(*model.User)
		tPage  = query.Use(l.svcCtx.Db).Page
		tTask  = query.Use(l.svcCtx.Db).Task
		now    = time.Now()
		file   *File
		dbTask *model.Task
	)
	query := tTask.WithContext(l.ctx).Where(tTask.ID.Eq(req.TaskID))
	if l.svcCtx.Config.ServerOptions.IsSaas {
		query = query.Where(tTask.CreaterUserID.Eq(u.ID))
	}
	if dbTask, err = query.First(); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	if file, err = getFile(l.ginCtx, "file"); err != nil {
		err = errors.WithCodeMsg(code.FileDoesNotExist)
		return
	}

	fileDir := filepath.Join(strconv.FormatInt(int64(u.ID), 10), strconv.FormatInt(int64(dbTask.ID), 10), now.Format("20060102"))
	filePath, err := saveFiles(base64.StdEncoding.EncodeToString(file.File), "jpg", fileDir)
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	page := &model.Page{
		TaskID:    dbTask.ID,
		FilePath:  filePath,
		OcrResult: new(string),
		Status:    consts.PageStatusCreated,
		CreateAt:  time.Now().Unix(),
		UpdateAt:  time.Now().Unix(),
	}
	if err = tPage.WithContext(l.ctx).Create(page); err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}
	resp.FileId = page.ID
	resp.ImageSrc = filePath
	return
}
