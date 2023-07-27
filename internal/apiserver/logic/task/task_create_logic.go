package task

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"financial_statement/internal/apiserver/code"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	dblog "financial_statement/internal/apiserver/db_log"
	jobTask "financial_statement/internal/apiserver/job/task"
	"financial_statement/internal/apiserver/svc"
	"financial_statement/internal/apiserver/types/task"
	storage "financial_statement/internal/pkg/storage"
	"financial_statement/pkg/errors"
	excel "financial_statement/pkg/excel"
	"financial_statement/pkg/log"
	md5 "financial_statement/pkg/md5"
	pdf "financial_statement/pkg/pdf"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskCreateLogic struct {
	ctx    context.Context
	ginCtx *gin.Context
	svcCtx *svc.ServiceContext
}

func NewTaskCreateLogic(ginCtx *gin.Context, serviceContext *svc.ServiceContext) TaskCreateLogic {
	return TaskCreateLogic{
		ctx:    context.Background(),
		ginCtx: ginCtx,
		svcCtx: serviceContext,
	}
}

// 下载外部文件
func getReferenceFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		log.Debugf("创建任务-->文件下载失败，http status：%d", resp.StatusCode)
		return "", fmt.Errorf("文件下载失败，http status：%d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(body), nil
}

func getFiles(files []task.CreateTaskFileItem) (base64List []string, err error) {
	for _, file := range files {
		// 判断文件是否存在
		if len(file.Base64) == 0 && len(file.Reference) == 0 {
			return nil, errors.WithCodeMsg(code.FileDoesNotExist)
		}
		base64 := file.Base64
		if len(base64) == 0 {
			base64, err = getReferenceFile(file.Reference)
			if err != nil {
				return nil, errors.WithCodeMsg(code.FileDoesNotExist)
			}
		}
		base64List = append(base64List, base64)
	}
	return base64List, nil
}

// 保存文件并，返回文件路径
func saveFiles(base64Str string, suffix string, fileDir string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	_uuid, _ := uuid.NewUUID()
	md5 := md5.MD5(_uuid.String())
	filePath := filepath.Join(fileDir, md5+"."+suffix)
	err = storage.FileStorage.Save(bytes, filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// TaskCreate 创建财报识别任务
func (l *TaskCreateLogic) TaskCreate(req *task.CreateTaskReq) (resp task.CreateTaskResp, err error) {
	var (
		q     = query.Use(l.svcCtx.Db)
		tTask = query.Use(l.svcCtx.Db).Task
		tUser = query.Use(l.svcCtx.Db).User
	)
	var u *model.User
	if l.ginCtx.Keys["user"] != nil {
		u = l.ginCtx.Keys["user"].(*model.User)
	} else {
		u, _ = tUser.WithContext(l.ctx).Where(tUser.Account.Eq("admin")).First()
	}
	if len(req.Files) == 0 {
		err = errors.WithCodeMsg(code.FileDoesNotExist)
		return
	}
	resp.FileType = req.FileType
	base64List := make([]string, 0)
	// 存储用户原始样本
	originalFiles := make([]string, 0)
	if strings.ToLower(req.FileType) == consts.TaskFileTypeExcel ||
		strings.ToLower(req.FileType) == consts.TaskFileTypePdf {
		tempList, e := getFiles(req.Files[:1])
		if e != nil {
			err = errors.WithCodeMsg(code.Internal, e.Error())
			return
		}

		fileDir := filepath.Join(strconv.FormatInt(int64(u.ID), 10), strconv.FormatInt(int64(resp.TaskID), 10), time.Now().Format("20060102"))
		suffix := "pdf"
		if strings.ToLower(req.FileType) == consts.TaskFileTypeExcel {
			suffix = "xlsx"
		}
		filePath, _ := saveFiles(tempList[0], suffix, fileDir)
		originalFiles = append(originalFiles, filePath) // 存储用户原始样本

		fileBytes, e := base64.StdEncoding.DecodeString(tempList[0])
		if e != nil {
			err = errors.WithCodeMsg(code.Internal, e.Error())
			return
		}

		if strings.ToLower(req.FileType) == consts.TaskFileTypeExcel {
			fileBytes, e = excel.ConvertToPDF(fileBytes)
			if e != nil {
				err = errors.WithCodeMsg(code.Internal, e.Error())
				return
			}
		}
		pdfPages, e := pdf.GetPages(fileBytes)
		if e != nil {
			err = errors.WithCodeMsg(code.BadRequest, "已损坏的文件，解析失败")
			return
		}
		if len(pdfPages) == 0 {
			err = errors.WithCodeMsg(code.BadRequest, "文件解析失败")
			return
		}
		for _, pageBase64 := range pdfPages {
			base64List = append(base64List, base64.StdEncoding.EncodeToString(pageBase64))
		}
	} else {
		base64List, err = getFiles(req.Files)
		if err != nil {
			log.Debugf("创建任务-->遇到错误:%s ,返回文件不存在", err.Error())
			err = errors.WithCodeMsg(code.FileDoesNotExist)
			return
		}
	}

	now := time.Now()
	err = q.Transaction(func(tx *query.Query) error {
		// 创建task
		t := &model.Task{
			FileFormat:    strings.ToLower(req.FileType),
			StandardID:    req.StandardID,
			ExternalInfo:  "",
			Files:         "[]",
			TaskName:      req.Name,
			Async:         int32(req.Async),
			AsyncStatus:   consts.TaskAsyncStatusUnsynchronized,
			Status:        consts.TaskStatusCreated,
			CreaterUserID: u.ID,
			Error:         new(string),
			CreatedAt:     now.Unix(),
			UpdatedAt:     now.Unix(),
		}
		if err := tx.Task.WithContext(l.ctx).Create(t); err != nil {
			return err
		}

		// 创建每页数据，并存储图片(图片按照uid/data/filename存储)
		for _, base64Str := range base64List {
			fileDir := filepath.Join(strconv.FormatInt(int64(u.ID), 10), strconv.FormatInt(int64(t.ID), 10), now.Format("20060102"))
			filePath, err := saveFiles(base64Str, "jpg", fileDir)
			originalFiles = append(originalFiles, filePath) // 存储用户原始样本
			if err != nil {
				return err
			}
			page := &model.Page{
				TaskID:    t.ID,
				FilePath:  filePath,
				OcrResult: new(string),
				Status:    consts.PageStatusCreated,
				CreateAt:  time.Now().Unix(),
				UpdateAt:  time.Now().Unix(),
			}
			err = tx.Page.WithContext(l.ctx).Create(page)
			if err != nil {
				return err
			}
			resp.Pages = append(resp.Pages, task.CreateTaskPagesItmeResp{
				Id:  page.ID,
				Url: filePath,
			})
		}

		resp.TaskID = t.ID
		for _, file := range req.Files {
			if len(file.FileNumber) > 0 {
				resp.FileInfo = append(resp.FileInfo, file.FileNumber)
			}
		}
		var externalInfo = consts.ExternalInfo{
			ConsumerId: req.ConsumerId,
			FileInfo:   resp.FileInfo,
			ExtraInfo:  req.ExtraInfo,
		}
		originalFilesJson, _ := json.Marshal(originalFiles)
		externalInfoJson, _ := json.Marshal(externalInfo)
		tx.Task.WithContext(l.ctx).Where(tTask.ID.Eq(t.ID)).Updates(model.Task{
			Files:        string(originalFilesJson),
			ExternalInfo: string(externalInfoJson),
		}) // 存储用户原始样本与附属信息

		return nil
	})
	if err != nil {
		err = errors.WithCodeMsg(code.Internal, err.Error())
		return
	}

	// 发布一个异步处理任务
	// l.svcCtx.JobMgr.EnqueueTask(jobTask.NewRecognizeTask(resp.TaskID))
	jobTask.NewRecognizeTask(&jobTask.RecoginzeTaskOption{
		TaskId:   resp.TaskID,
		TaskType: jobTask.TaskTypeRecognized,
	})
	dblog.DbLog(l.ctx, resp.TaskID, fmt.Sprintf("Create Task By User:%s", u.Account))
	return
}
