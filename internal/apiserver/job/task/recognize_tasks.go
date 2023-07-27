package jobtask

import (
	"bytes"
	"context"
	"encoding/json"
	"financial_statement/internal/apiserver/consts"
	"financial_statement/internal/apiserver/dal"
	"financial_statement/internal/apiserver/dal/model"
	"financial_statement/internal/apiserver/dal/query"
	dblog "financial_statement/internal/apiserver/db_log"
	taskresultasync "financial_statement/internal/apiserver/job/task-result-async"
	"financial_statement/internal/apiserver/types/task"
	fshelper "financial_statement/internal/pkg/fs_helper"
	storage "financial_statement/internal/pkg/storage"
	taskhelper "financial_statement/internal/pkg/task_helper"
	"financial_statement/pkg/excel"
	"financial_statement/pkg/log"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/antlabs/strsim"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"github.com/spf13/viper"
	"github.com/wxnacy/wgo/arrays"
)

// A list of task types.
const (
	RecognizeTask = "task:recognize_task"
	//识别任务
	TaskTypeRecognized = "TaskTypeRecognized"
	//重新识别的任务，重新识别的任务在回调时有特殊处理
	TaskTypeReIdentify = "TaskTypeReIdentify"
)

var (
	once              sync.Once
	tableRecognizeApi string
	db                dal.Repo
)

type RecoginzeTaskOption struct {
	TaskId   uint32
	TaskType string
}

// 传递任务id，生成一个异步处理任务
func NewRecognizeTask(taskOption *RecoginzeTaskOption) (*asynq.TaskInfo, error) {
	asynqredisIns, _ := dal.GetRedisFactoryOr(nil)
	//初始化异步任务客户端对象
	client := asynqredisIns.GetAsynqRedis()
	payload, _ := json.Marshal(taskOption)

	recognizeMaxRetry := viper.GetInt("server.recognizeMaxRetry")
	return client.Enqueue(asynq.NewTask(RecognizeTask, payload), asynq.MaxRetry(recognizeMaxRetry))
}

func ErrorHandler(ctx context.Context, payload []byte) {
	var (
		tTask = query.Use(db.GetDb()).Task
	)
	var taskOption RecoginzeTaskOption
	var id uint32
	err := json.Unmarshal(payload, &taskOption)
	if err != nil {
		if err = json.Unmarshal(payload, &id); err != nil {
			log.Errorf("HandleRecognizeTask TaskId:%s\n", err.Error())
		} else {
			taskOption.TaskId = id //兼容以前异步任务信息里只存放一个id的情况
		}
	}
	id = taskOption.TaskId
	if taskOption.TaskId > 0 {
		if _, err := tTask.WithContext(ctx).Where(tTask.ID.Eq(taskOption.TaskId)).Update(tTask.Status, consts.TaskStatusOcrFailed); err != nil {
			log.Errorf("HandleRecognizeTask TaskId:%s\n", err.Error())
		} else {
			asyncTaskResult(&taskOption)
		}
	}
	errorLog(ctx, taskOption.TaskId, "任务识别失败超出最大尝试次数，判定为识别失败")
}

// 异步任务处理
func HandleRecognizeTask(ctx context.Context, t *asynq.Task) error {
	var taskOption RecoginzeTaskOption
	var id uint32
	err := json.Unmarshal(t.Payload(), &taskOption)
	if err != nil {
		if err = json.Unmarshal(t.Payload(), &id); err != nil {
			log.Errorf("HandleRecognizeTask TaskId:%s\n", err.Error())
			return err
		} else {
			taskOption.TaskId = id //兼容以前异步任务信息里只存放一个id的情况
		}
	}
	id = taskOption.TaskId
	log.Debugf("start recognize task id : %d", taskOption.TaskId)
	if err := recognizeTask(context.Background(), taskOption.TaskId); err != nil {
		errorLog(context.Background(), taskOption.TaskId, err.Error())
		return err
	} else {
		dblog.DbLog(context.Background(), taskOption.TaskId, "task recognize finished")
		asyncTaskResult(&taskOption)
		return nil
	}
}

//同步给第三方系统
func asyncTaskResult(taskOption *RecoginzeTaskOption) {
	dblog.DbLog(context.Background(), taskOption.TaskId, "task recognize finished")
	reIdentifyWithSync := viper.GetBool("task-result-async.re-identify-with-sync")
	switch taskOption.TaskType {
	case TaskTypeReIdentify:
		if reIdentifyWithSync {
			taskresultasync.NewAsyncTask(&taskresultasync.SyncTaskInfo{
				TaskId:       taskOption.TaskId,
				TaskType:     taskresultasync.SyncTaskTypeRecognized,
				CurrentIndex: -1,
			})
		}
	case TaskTypeRecognized:
		taskresultasync.NewAsyncTask(&taskresultasync.SyncTaskInfo{
			TaskId:       taskOption.TaskId,
			TaskType:     taskresultasync.SyncTaskTypeRecognized,
			CurrentIndex: -1,
		})
	}
}

func errorLog(ctx context.Context, id uint32, err string) {
	once.Do(func() {
		db, _ = dal.GetDbFactoryOr(nil)
	})
	var (
		tTask = query.Use(db.GetDb()).Task
	)
	tTask.WithContext(ctx).Where(tTask.ID.Eq(id)).Update(tTask.Error, err)
	dblog.DbLog(ctx, id, err)
}

// 识别财报任务
func recognizeTask(ctx context.Context, id uint32) error {
	once.Do(func() {
		tableRecognizeApi = viper.GetString("ocr.recognize_table_api")
		db, _ = dal.GetDbFactoryOr(nil)
	})

	var (
		tTask              = query.Use(db.GetDb()).Task
		tPage              = query.Use(db.GetDb()).Page
		tTitle             = query.Use(db.GetDb()).StandardStatementTitle
		tStandardStatement = query.Use(db.GetDb()).StandardStatement
		tStandard          = query.Use(db.GetDb()).Standard
	)

	dbTask, err := tTask.WithContext(ctx).Select(tTask.ALL).Where(tTask.ID.Eq(id)).First()
	if err != nil {
		log.Errorf("RecognizeTask Find Task(id:%d) With Err: %s", id, err.Error())
		return err
	}
	if dbTask.Status != consts.TaskStatusCreated {
		return nil
	}

	// 获取该财报任务所用的财报准则里的所有科目信息，为后续科目匹配做准备
	standardStatementTitles := new([]StandardStatementTitle)
	err = tTitle.WithContext(ctx).
		Select(tTitle.ID, tTitle.Name, tTitle.Aliases, tStandardStatement.Type).
		LeftJoin(tStandardStatement, tStandardStatement.ID.EqCol(tTitle.StatementID)).
		LeftJoin(tStandard, tStandard.ID.EqCol(tStandardStatement.StandardID)).
		Where(tTitle.Status.Eq(consts.TitleStatusNormal), tStandard.ID.Eq(dbTask.StandardID)).Scan(&standardStatementTitles)
	if err != nil {
		log.Errorf("Get task titles With Err: %s", err.Error())
		return err
	}
	var externalInfo consts.ExternalInfo
	json.Unmarshal([]byte(dbTask.ExternalInfo), &externalInfo)

	// 开始识别每一页图片
	financialStatement := make([]task.TaskFinancialStatement, 0)
	pages, err := tPage.WithContext(ctx).Where(tPage.TaskID.Eq(dbTask.ID), tPage.Status.Eq(consts.PageStatusCreated)).Find()
	if err != nil {
		log.Errorf("RecognizeTask Find Task Pages With Err: %s", err.Error())
		return err
	}

	if externalInfo.Groups != nil && len(dbTask.ExternalInfo) > 0 {
		// 调整过分组
		financialStatement, err = recognizeTaskPagesWithGroup(pages, standardStatementTitles, &externalInfo.Groups)
	} else {
		// 常规识别
		financialStatement, err = recognizeTaskPages(pages, standardStatementTitles)
	}

	// sortPeriodTitle(&financialStatement)
	fixHeaderAndCells(&financialStatement)
	if err != nil {
		return err
	}
	//记录识别耗时到externalInfo
	externalInfo.RecognizeDuration = time.Now().Unix() - dbTask.CreatedAt
	extJson, err := json.Marshal(externalInfo)
	if err != nil {
		log.Errorf("json.Marshal(externalInfo): %s", err.Error())
		return err
	}
	//序列化识别结果以保存到DB
	frJson, err := json.Marshal(&financialStatement)
	if err != nil {
		log.Errorf("json.Marshal(financialStatement): %s", err.Error())
		return err
	}
	var filePath string
	if filePath, err = taskhelper.SaveTaskResult(dbTask, frJson); err != nil {
		return err
	}
	_, err = tTask.WithContext(ctx).Where(tTask.ID.Eq(id)).Updates(model.Task{
		StandardResult: &filePath,
		Status:         consts.TaskStatusOcrSuccess,
		ExternalInfo:   string(extJson),
	})
	if err != nil {
		log.Errorf("Updates 财报结果更新到DB出错: %s", err.Error())
		return err
	}
	return nil
}
func sortPeriodTitle(financialStatement *[]task.TaskFinancialStatement) {
	for _, fs := range *financialStatement {
		sortHeaders(&fs.BalanceSheet.Headers)
		sortHeaders(&fs.CashFlowStatement.Headers)
		sortHeaders(&fs.IncomeStatement.Headers)
	}
}

func sortHeaders(headers *[]task.Header) {
	headerDates := make([]string, len(*headers))
	for i, header := range *headers {
		headerDates[i] = header.Date
	}
	headerDates = fshelper.GetPeriodHeaderOrderIndex(headerDates)
	for i, item := range *headers {
		for j, val := range headerDates {
			if item.Date == val {
				(*headers)[i].Order = j
				break
			}
		}
	}
}

// Contains 数组是否包含某元素
func contains(slice []int, s int) int {
	for index, value := range slice {
		if value == s {
			return index
		}
	}
	return -1
}

func recognizeTaskPagesWithGroup(pages []*model.Page, titles *[]StandardStatementTitle, groups *[]consts.ExternalInfoGroup) (financialStatements []task.TaskFinancialStatement, err error) {
	var periods []Period
	for _, group := range *groups {
		financialStatement := task.TaskFinancialStatement{
			BalanceSheet:      task.Statement{},
			IncomeStatement:   task.Statement{},
			CashFlowStatement: task.Statement{},
		}
		var statementType int
		for _, file := range group.Files {
			statementType = file.Type
			var page *model.Page
			var tableOcrResult RecognizeTableResult
			for _, p := range pages {
				if p.ID == file.FileId {
					page = p
					if err = RecognizeTaskWithRetry(p, &tableOcrResult); err != nil {
						return financialStatements, err
					}
				}
			}
			for _, table := range tableOcrResult.Result.Tables {
				if table.Type == tableTypePlain {
					continue
				}
				//解析财报数据
				statement, _ := getStatement(statementType, table, titles, &periods)
				// periods = append(periods, tempPeriods...)

				// 根据行列，对科目的重新排序
				sort.Stable(TitleSlice(statement.Titles))
				sort.SliceStable(statement.Titles, func(i, j int) bool {
					return statement.Titles[i].ColIndex < statement.Titles[j].ColIndex
				})

				switch statementType {
				case balanceSheet:
					//判断是新一份财报还是财报数据跨页了
					var headers = make([]string, 0)
					for _, item := range financialStatement.BalanceSheet.Headers {
						headers = append(headers, item.Date)
					}
					for _, item := range statement.Headers {
						if arrays.ContainsString(headers, item.Date) == -1 {
							financialStatement.BalanceSheet.Headers = append(financialStatement.BalanceSheet.Headers, item)
						}
					}
					if statement.Count >= financialStatement.BalanceSheet.Count {
						financialStatement.BalanceSheet.Count = statement.Count
					}
					image := task.Image{
						ImageSrc:    page.FilePath,
						RotateAngle: tableOcrResult.Result.Angle,
						FileId:      page.ID,
					}
					// if arrays.Contains(financialStatement.BalanceSheet.Images, image) == -1 {
					// 	financialStatement.BalanceSheet.Images = append(financialStatement.BalanceSheet.Images, image)
					// }
					// 调整分组支持存在两个相同图片路径，因为有复制的功能
					financialStatement.BalanceSheet.Images = append(financialStatement.BalanceSheet.Images, image)
					for index, _ := range statement.Titles {
						statement.Titles[index].PageIndex = len(financialStatement.BalanceSheet.Images) - 1
					}
					financialStatement.BalanceSheet.Titles = append(financialStatement.BalanceSheet.Titles, statement.Titles...)
				case cashFlow:
					var headers = make([]string, 0)
					for _, item := range financialStatement.CashFlowStatement.Headers {
						headers = append(headers, item.Date)
					}
					for _, item := range statement.Headers {
						if arrays.ContainsString(headers, item.Date) == -1 {
							financialStatement.CashFlowStatement.Headers = append(financialStatement.CashFlowStatement.Headers, item)
						}
					}
					if statement.Count >= financialStatement.CashFlowStatement.Count {
						financialStatement.CashFlowStatement.Count = statement.Count
					}
					image := task.Image{
						ImageSrc:    page.FilePath,
						RotateAngle: tableOcrResult.Result.Angle,
						FileId:      page.ID,
					}
					// if arrays.Contains(financialStatement.CashFlowStatement.Images, image) == -1 {
					// 	financialStatement.CashFlowStatement.Images = append(financialStatement.CashFlowStatement.Images, image)
					// }
					// 调整分组支持存在两个相同图片路径，因为有复制的功能
					financialStatement.CashFlowStatement.Images = append(financialStatement.CashFlowStatement.Images, image)
					for index, _ := range statement.Titles {
						statement.Titles[index].PageIndex = len(financialStatement.CashFlowStatement.Images) - 1
					}
					financialStatement.CashFlowStatement.Titles = append(financialStatement.CashFlowStatement.Titles, statement.Titles...)

				case income:
					var headers = make([]string, 0)
					for _, item := range financialStatement.IncomeStatement.Headers {
						headers = append(headers, item.Date)
					}
					for _, item := range statement.Headers {
						if arrays.ContainsString(headers, item.Date) == -1 {
							financialStatement.IncomeStatement.Headers = append(financialStatement.IncomeStatement.Headers, item)
						}
					}
					if statement.Count >= financialStatement.IncomeStatement.Count {
						financialStatement.IncomeStatement.Count = statement.Count
					}
					image := task.Image{
						ImageSrc:    page.FilePath,
						RotateAngle: tableOcrResult.Result.Angle,
						FileId:      page.ID,
					}
					// if arrays.Contains(financialStatement.IncomeStatement.Images, image) == -1 {
					// 	financialStatement.IncomeStatement.Images = append(financialStatement.IncomeStatement.Images, image)
					// }
					// 调整分组支持存在两个相同图片路径，因为有复制的功能
					financialStatement.IncomeStatement.Images = append(financialStatement.IncomeStatement.Images, image)
					for index, _ := range statement.Titles {
						statement.Titles[index].PageIndex = len(financialStatement.IncomeStatement.Images) - 1
					}
					financialStatement.IncomeStatement.Titles = append(financialStatement.IncomeStatement.Titles, statement.Titles...)
				}
			}
		}
		financialStatements = append(financialStatements, financialStatement)
	}

	return financialStatements, err
}

func saveOcrResult(tableOcrResult *RecognizeTableResult, id uint32) error {
	var (
		tPage     = query.Use(db.GetDb()).Page
		ocrResult consts.PageOcrResult
	)

	copier.Copy(&ocrResult, &tableOcrResult)
	if jsonb, err := json.Marshal(ocrResult); err == nil {
		if _, err := tPage.WithContext(context.Background()).Where(tPage.ID.Eq(id)).Update(tPage.OcrResult, string(jsonb)); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

// 识别财报的每一页数据
func recognizeTaskPages(pages []*model.Page, titles *[]StandardStatementTitle) (financialStatements []task.TaskFinancialStatement, err error) {

	var statementType int
	var statementTypecache []int
	var periods []Period
	financialStatement := task.TaskFinancialStatement{
		BalanceSheet:      task.Statement{},
		IncomeStatement:   task.Statement{},
		CashFlowStatement: task.Statement{},
	}
	for pageIndex, page := range pages {
		var tableOcrResult RecognizeTableResult
		runtime.GC()
		log.Debugf("page id:%d", page.ID)
		err = RecognizeTaskWithRetry(page, &tableOcrResult)
		if err != nil {
			break
		}
		saveOcrResult(&tableOcrResult, page.ID)
		for tId, table := range tableOcrResult.Result.Tables {
			log.Debugf("当前处理page：%d,第%d/%d个table", page.ID, tId+1, len(tableOcrResult.Result.Tables))
			if table.Type == tableTypePlain {
				continue
			}
			statementType = getTableType(table)
			if statementType == 0 {
				continue
			}
			// 数据跨页逻辑 START
			if index := contains(statementTypecache, statementType); index >= 0 {
				// 类型缓存中有该类型
				if statementTypecache[len(statementTypecache)-1] == statementType {
					//同一份财报，并且数据跨页了
					log.Debugf("数据跨页了！")
				} else if len(statementTypecache) >= 3 {
					//非跨页，是一份新的财报
					log.Debugf("非跨页，是一份新的财报")
					periods = make([]Period, 0)
					financialStatements = append(financialStatements, financialStatement)
					statementTypecache = make([]int, 0) //清理临时表
					//重新初始化一份财报对象
					financialStatement = task.TaskFinancialStatement{
						BalanceSheet:      task.Statement{},
						IncomeStatement:   task.Statement{},
						CashFlowStatement: task.Statement{},
					}
				}

			} else {
				// 属于同一份财报
				log.Debugf("属于同一份财报")
				periods = make([]Period, 0)
				statementTypecache = append(statementTypecache, statementType)
			}
			// 数据跨页逻辑 END

			//解析财报数据
			statement, _ := getStatement(statementType, table, titles, &periods)
			// periods = append(periods, tempPeriods...)

			log.Debugf("科目排序Start")
			// 根据行列，对科目的重新排序
			sort.Stable(TitleSlice(statement.Titles))
			sort.SliceStable(statement.Titles, func(i, j int) bool {
				return statement.Titles[i].ColIndex < statement.Titles[j].ColIndex
			})
			log.Debugf("科目排序END")

			switch statementType {
			case balanceSheet:
				//判断是新一份财报还是财报数据跨页了
				var headers = make([]string, 0)
				for _, item := range financialStatement.BalanceSheet.Headers {
					headers = append(headers, item.Date)
				}
				for _, item := range statement.Headers {
					if arrays.ContainsString(headers, item.Date) == -1 {
						financialStatement.BalanceSheet.Headers = append(financialStatement.BalanceSheet.Headers, item)
					}
				}
				if statement.Count >= financialStatement.BalanceSheet.Count {
					financialStatement.BalanceSheet.Count = statement.Count
				}
				financialStatement.BalanceSheet.Images = append(financialStatement.BalanceSheet.Images, task.Image{
					ImageSrc:    page.FilePath,
					RotateAngle: tableOcrResult.Result.Angle,
					FileId:      page.ID,
				})
				for index, _ := range statement.Titles {
					statement.Titles[index].PageIndex = len(financialStatement.BalanceSheet.Images) - 1
				}
				financialStatement.BalanceSheet.Titles = append(financialStatement.BalanceSheet.Titles, statement.Titles...)

			case cashFlow:
				var headers = make([]string, 0)
				for _, item := range financialStatement.CashFlowStatement.Headers {
					headers = append(headers, item.Date)
				}
				for _, item := range statement.Headers {
					if arrays.ContainsString(headers, item.Date) == -1 {
						financialStatement.CashFlowStatement.Headers = append(financialStatement.CashFlowStatement.Headers, item)
					}
				}
				if statement.Count >= financialStatement.CashFlowStatement.Count {
					financialStatement.CashFlowStatement.Count = statement.Count
				}
				financialStatement.CashFlowStatement.Images = append(financialStatement.CashFlowStatement.Images, task.Image{
					ImageSrc:    page.FilePath,
					RotateAngle: tableOcrResult.Result.Angle,
					FileId:      page.ID,
				})
				for index, _ := range statement.Titles {
					statement.Titles[index].PageIndex = len(financialStatement.CashFlowStatement.Images) - 1
				}
				financialStatement.CashFlowStatement.Titles = append(financialStatement.CashFlowStatement.Titles, statement.Titles...)

			case income:
				var headers = make([]string, 0)
				for _, item := range financialStatement.IncomeStatement.Headers {
					headers = append(headers, item.Date)
				}
				for _, item := range statement.Headers {
					if arrays.ContainsString(headers, item.Date) == -1 {
						financialStatement.IncomeStatement.Headers = append(financialStatement.IncomeStatement.Headers, item)
					}
				}
				if statement.Count >= financialStatement.IncomeStatement.Count {
					financialStatement.IncomeStatement.Count = statement.Count
				}
				financialStatement.IncomeStatement.Images = append(financialStatement.IncomeStatement.Images, task.Image{
					ImageSrc:    page.FilePath,
					RotateAngle: tableOcrResult.Result.Angle,
					FileId:      page.ID,
				})
				for index, _ := range statement.Titles {
					statement.Titles[index].PageIndex = len(financialStatement.IncomeStatement.Images) - 1
				}
				financialStatement.IncomeStatement.Titles = append(financialStatement.IncomeStatement.Titles, statement.Titles...)

			}

			// if newStatement {
			// 	financialStatements = append(financialStatements, financialStatement)
			// 	//新的财报分组
			// 	financialStatement = task.TaskFinancialStatement{
			// 		BalanceSheet:      task.Statement{},
			// 		IncomeStatement:   task.Statement{},
			// 		CashFlowStatement: task.Statement{},
			// 	}
			// }

		}
		log.Debugf("page id:%d END1")
		if pageIndex == len(pages)-1 &&
			(financialStatement.BalanceSheet.Count > 0 ||
				financialStatement.CashFlowStatement.Count > 0 ||
				financialStatement.IncomeStatement.Count > 0) {
			financialStatements = append(financialStatements, financialStatement)
		}
		log.Debugf("page id:%d END2")
	}
	return financialStatements, err
}

// 补全缺失的Header and cells
func fixHeaderAndCells(fs *[]task.TaskFinancialStatement) {
	for i := range *fs {
		item := &(*fs)[i]

		headerCount := 0
		for _, title := range item.BalanceSheet.Titles {
			if headerCount < len(title.Values) {
				headerCount = len(title.Values)
			}
		}
		if len(item.BalanceSheet.Headers) < headerCount {
			count := headerCount - len(item.BalanceSheet.Headers)
			for j := 0; j < count; j++ {
				item.BalanceSheet.Headers = append(item.BalanceSheet.Headers, task.Header{
					Organization: "",
					Date:         "",
					Order:        len(item.BalanceSheet.Headers) + j + 1,
				})
			}
		}

		headerCount = len(item.BalanceSheet.Headers)
		for j := range item.BalanceSheet.Titles {
			title := &(item.BalanceSheet.Titles)[j]
			if len(title.Values) < headerCount {
				count := headerCount - len(title.Values)
				for x := 0; x < count; x++ {
					title.Values = append(title.Values, task.StatementTitleValue{
						Ocr:        "",
						Supervised: "",
						Position:   []int{},
					})
				}
			}
		}

		headerCount = 0

		for _, title := range item.CashFlowStatement.Titles {
			if headerCount < len(title.Values) {
				headerCount = len(title.Values)
			}
		}
		if len(item.CashFlowStatement.Headers) < headerCount {
			count := headerCount - len(item.CashFlowStatement.Headers)
			for j := 0; j < count; j++ {
				item.CashFlowStatement.Headers = append(item.CashFlowStatement.Headers, task.Header{
					Organization: "",
					Date:         "",
					Order:        len(item.CashFlowStatement.Headers) + j + 1,
				})
			}
		}
		headerCount = len(item.CashFlowStatement.Headers)
		for j := range item.CashFlowStatement.Titles {
			title := &(item.CashFlowStatement.Titles)[j]
			if len(title.Values) < headerCount {
				count := headerCount - len(title.Values)
				for x := 0; x < count; x++ {
					title.Values = append(title.Values, task.StatementTitleValue{
						Ocr:        "",
						Supervised: "",
						Position:   []int{},
					})
				}
			}
		}

		headerCount = 0

		for _, title := range item.IncomeStatement.Titles {
			if headerCount < len(title.Values) {
				headerCount = len(title.Values)
			}
		}
		if len(item.IncomeStatement.Headers) < headerCount {
			count := headerCount - len(item.IncomeStatement.Headers)
			for j := 0; j < count; j++ {
				item.IncomeStatement.Headers = append(item.IncomeStatement.Headers, task.Header{
					Organization: "",
					Date:         "",
					Order:        len(item.IncomeStatement.Headers) + j + 1,
				})
			}
		}
		headerCount = len(item.IncomeStatement.Headers)
		for j := range item.IncomeStatement.Titles {
			title := &(item.IncomeStatement.Titles)[j]
			if len(title.Values) < headerCount {
				count := headerCount - len(title.Values)
				for x := 0; x < count; x++ {
					title.Values = append(title.Values, task.StatementTitleValue{
						Ocr:        "",
						Supervised: "",
						Position:   []int{},
					})
				}
			}
		}

		headerCount = 0
	}
}

//将excel结果转换为OCR结果，然后走后续流程
func getExcelResult(page *model.Page, ocrResult *RecognizeTableResult) error {
	var excelTable excel.Table
	if err := json.Unmarshal([]byte(*page.OcrResult), &excelTable); err != nil {
		return err
	}
	table := RecognizeTableResultTable{
		Type:         "table_with_line",
		Position:     []int{},
		TableRows:    0,
		TableCols:    0,
		HeightOfRows: []int{},
		WidthOfCols:  []int{},
		TableCells:   []RecognizeTableResultTableCell{},
	}
	for _, cell := range excelTable.TableCells {
		table.TableCells = append(table.TableCells, RecognizeTableResultTableCell{
			StartRow: cell.StartRow,
			StartCol: cell.StartCol,
			EndRow:   cell.EndRow,
			EndCol:   cell.EndCol,
			Text:     cell.Text,
			Position: []int{},
			Semantic: "",
			Order:    0,
			Lines: []struct {
				Text                string      "json:\"text\""
				Score               float64     "json:\"score\""
				Type                string      "json:\"type\""
				Position            []int       "json:\"position\""
				Angle               int         "json:\"angle\""
				Direction           int         "json:\"direction\""
				Handwritten         int         "json:\"handwritten\""
				CharScores          []float64   "json:\"char_scores\""
				CharCenters         [][]int     "json:\"char_centers\""
				CharPositions       [][]int     "json:\"char_positions\""
				CharCandidates      [][]string  "json:\"char_candidates\""
				CharCandidatesScore [][]float64 "json:\"char_candidates_score\""
				Semantic            string      "json:\"semantic\""
			}{},
		})
	}
	ocrResult.Result.Tables = append(ocrResult.Result.Tables, table)
	return nil
}

func RecognizeTaskWithRetry(page *model.Page, tableOcrResult *RecognizeTableResult) (err error) {
	maxRetries := 3
	retryInterval := 1 * time.Second

	for i := 0; i < maxRetries; i++ {
		err = getOcrResult(page, tableOcrResult)
		if err == nil {
			return nil
		}
		log.Warnf("Recognize Task failed. Retrying in %s...\n", retryInterval)
		time.Sleep(retryInterval)
	}

	log.Warnf("Task failed after %d retries，err:%s", maxRetries, err.Error())
	return err
}

func getOcrResult(page *model.Page, tableOcrResult *RecognizeTableResult) error {
	fileByte, err := storage.FileStorage.Get(page.FilePath)
	if err != nil {
		return fmt.Errorf("获取任务：%d 的页文件（id：%d）失败：%s", page.TaskID, page.ID, err.Error())
	}
	body := bytes.NewReader(fileByte)
	httpReq, err := http.NewRequest("POST", tableRecognizeApi, body)
	if err != nil {
		return fmt.Errorf("创建Request请求失败，任务：%d 的页文件（id：%d）失败：%s", page.TaskID, page.ID, err.Error())
	}
	httpResp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("OCR识别任务：%d 的页文件（id：%d）失败：%s", page.TaskID, page.ID, err.Error())
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("OCR识别任务：%d 的页文件（id：%d） http status code: %d", page.TaskID, page.ID, httpResp.StatusCode)
	}
	resBody, _ := ioutil.ReadAll(httpResp.Body)
	_ = json.Unmarshal(resBody, &tableOcrResult)
	if tableOcrResult.Code != http.StatusOK {
		return fmt.Errorf("解析OCR数据的code码不为200，status：%d msg:%s", tableOcrResult.Code, tableOcrResult.Message)
	}
	return nil
}

const (
	//资产负债表
	balanceSheetRegStr   = "流动负债.{0,3}|合并资产负债表|科目|.{0,2}流动资产.{0,2}|货币资金|.?结算备付金|.?拆出资金|交易性金融资产|以公允价值计量且其变动计入当期损益的金融资产|衍生金融资产|应[收付]票据|应[收付]账款|应[收付]款项融资|预[收付]款项|.?应[收付]保费|.?应[收付]分保账款|.?应收分保合同准备金|其他应[收付]款|(其中)?.?应[收付]利息|应[收付]股利|.?买入返售金融资产|存货|合同资产|持有待售资产|一年内到期的非流动资产|其他流动资产|流动资产合计|非流动资产.{0, 2}|.?发放贷款和垫款|债券投资|可供出售金融资产|其他债券投资|持有至到期投资|长期应收款|长期股权投资|其他权益工具投资|其他非流动金融资产|投资性房地产|固定资产|在建工程|生产性生物资产|油气资产|使用权资产|无形资产|开发支出|商誉|长期待摊费用|递延所得税资产|其他非流动资产|非流动资产合计|资产总计|资本公积|.{0,6}库存股|其他综合收益|专用储备|盈余公积|未分配利润|所有者权益（或股东权益）合计|实收资本|所有者权益合计|负债和所有者权益统计|所有者权益|应收票据|应收帐款|预付款项|应收利息|应收股利|其他应收款|工程物资|固定资产清理|短期借款|交易性金融负债|应付票据|应付账款|预收款项|合同负债|应付职工薪酬|应交税费|应付利息|应付股利|其他应付款|一年内到期的非流动负债.{0,6}|其他流动负债|流动负债合计|长期借款|应付债券|长期应付款|专项应付款|预计负债|递延所得税负债|其他非流动负债|非流动负债合计|负债合计|实收资本（或股本）|负债和所有者权益（或股东权益）总计"
	balanceSheetTcRegStr = "流動負債.{0,3}|合併資產負債表|科目|.{0,2}流動資產.{0,2}|貨幣資金|.?結算備付金|.?拆出資金|交易性金融資產|以公允價值計量且其變動計入當期損益的金融資產|衍生金融資產|應[收付]票據|應[收付]帳款|應[收付]款項融資|預[收付]款項|.?應[收付]保費|.?應[收付]分保帳款|.?應收分保合同準備金|其他應[收付]款|(其中)?.?應[收付]利息|應[收付]股利|.?買入返售金融資產|存貨|合同資產|持有待售資產|一年內到期的非流動資產|其他流動資產|流動資產合計|非流動資產.{0, 2}|.?發放貸款和垫款|債券投資|可供出售金融資產|其他債券投資|持有至到期投資|長期應收款|長期股權投資|其他權益工具投資|其他非流動金融資產|投資性房地產|固定資產|在建工程|生產性生物資產|油氣資產|使用權資產|無形資產|開發支出|商譽|長期待攤費用|递延所得税資產|其他非流動資產|非流動資產合計|資產總計|資本公積|.{0,6}庫存股|其他綜合收益|專用儲備|盈余公積|未分配利潤|所有者權益（或股東權益）合計|實收資本|所有者權益合計|負債和所有者權益統計|所有者權益|應收票據|應收帳款|預付款項|應收利息|應收股利|其他應收款|工程物資|固定資產清理|短期借款|交易性金融負債|應付票據|應付帳款|預收款項|合同負債|應付職工薪酬|應交稅費|應付利息|應付股利|其他應付款|一年內到期的非流動負債.{0,6}|其他流動負債|流動負債合計|長期借款|應付債券|長期應付款|專項應付款|預計負債|递延所得税負債|其他非流動負債|非流動負債合計|負債合計|實收資本（或股本）|負債和所有者權益（或股東權益）總計"
	//现金流量表
	cashFlowStatementRegStr   = ".*现金流量表.*|科目|销售商品.?提供劳务收到的现金|收到的税费返还|收到其他与经营活动有关的现金|.*经营活动现金流入小计.*|购买商品.?接受劳务支付的现金|支付给职工以及为职工支付的现金|支付的各项税费|支付其他与经营活动有关的现金|经营活动现金流出小计|经营活动产生的现金流量净额|取得投资收益收到的现金|处置固定资产.?无形资产和其他长期资产收回的现金净额|收到其他与投资活动有关的现金|投资活动现金流入小计|购建固定资产.?无形资产和其他长期资产支付的现金|支付其他与投资活动有关的现金|投资活动现金流出小计|投资活动产生的现金流量净额|取得借款收到的现金|收到其他与筹资活动有关的现金|筹资活动现金流入小计|偿还债务支付的现金|分配股利.?利润或偿付利息支付的现金|筹资活动现金流出小计|筹资活动产生的现金流量净额|.*现金及现金等价物净增加额.*|加.*期初现金及现金等价物的余额|期末现金及现金等价物余额|.*购买商品、接受劳务支付的现金.*|.?资活动产生的现金流量.?|收回投资收到的现金|处置固定资产、无形资产和其他长期资产收回的现金净额|处置子公司及其他营业单位收到的现金净额|购建固定资产、无形资产和其他长期资产支付的现金|投资支付的现金|吸收投资收到的现金|分配股利、利润或偿付利息支付的现金|支付其他与筹资活动有关的现金|.*汇率变动对现金的影响.*|.*现金及现金等价物余额.*|销售商品、提供劳务收到的现金|客户存款和同业存放款项净增加额|.*向中央银行借款净增加额.*|.*向其他金融机构拆入资金净增加额.*|.*收到原保险合同保费取得的现金.*|.*收到再保险业务现金净额.*|.*保户储金及投资款净增加额.*|.*处置以公允价值计量且其变动计入.*|.*收取利息、手续费及佣金的现金.*|.*拆入资金净增加额.*|.*回购业务资金净增加额.*|.*收到的税费返还.*|.*收到其他与经营活动有关的现金.*|.*客户贷款及垫款净增加额.*|.*存放中央银行和同业款项净增加额.*|.*支付原保险合同赔付款项的现金.*|短期借款、长期借款借方数字|短期借款、长期借款贷方数字|备注资料|吸收权益性资金"
	cashFlowStatementTcRegStr = ".*現金流量表.*|科目|銷售商品.?提供勞務收到的現金|收到的稅費返還|收到其他與經營活動有關的現金|.*經營活動現金流入小計.*|購買商品.?接受勞務支付的現金|支付給職工以及為職工支付的現金|支付的各項稅費|支付其他與經營活動有關的現金|經營活動現金流出小計|經營活動產生的現金流量淨額|取得投資收益收到的現金|處置固定資產.?無形資產和其他長期資產收回的現金淨額|收到其他與投資活動有關的現金|投資活動現金流入小計|購建固定資產.?無形資產和其他長期資產支付的現金|支付其他與投資活動有關的現金|投資活動現金流出小計|投資活動產生的現金流量淨額|取得借款收到的現金|收到其他與籌資活動有關的現金|籌資活動現金流入小計|償還債務支付的現金|分配股利.?利潤或償付利息支付的現金|籌資活動現金流出小計|籌資活動產生的現金流量淨額|.*現金及現金等價物淨增加額.*|加.*期初現金及現金等價物的餘額|期末現金及現金等價物餘額|.*購買商品、接受勞務支付的現金.*|.?資活動產生的現金流量.?|收回投資收到的現金|處置固定資產、無形資產和其他長期資產收回的現金淨額|處置子公司及其他營業單位收到的現金淨額|購建固定資產、無形資產和其他長期資產支付的現金|投資支付的現金|吸收投資收到的現金|分配股利、利潤或償付利息支付的現金|支付其他與籌資活動有關的現金|.*匯率變動對現金的影響.*|.*現金及現金等價物餘額.*|銷售商品、提供勞務收到的現金|客戶存款和同業存放款項淨增加額|.*向中央銀行借款淨增加額.*|.*向其他金融機構拆入資金淨增加額.*|.*收到原保險合同保費取得的現金.*|.*收到再保險業務現金淨額.*|.*保戶儲金及投資款淨增加額.*|.*處置以公允價值計量且其變動計入.*|.*收取利息、手續費及佣金的現金.*|.*拆入資金淨增加額.*|.*回購業務資金淨增加額.*|.*收到的稅費返還.*|.*收到其他與經營活動有關的現金.*|.*客戶貸款及擔款淨增加額.*|.*存放中央銀行和同業款項淨增加額.*|.*支付原保險合同賠付款項的現金.*|短期借款、長期借款借方數字|短期借款、長期借款貸方數字|備註資料|吸收權益性資金"
	//利润表
	incomeStatementRegStr   = "合并利润表.?损益表.?收益表|已赚保费|科目|营业总收入|减.?营业成本|.*营业收入|营业总成本|.*营业成本|税金及附加|销售费用|管理费用|研发费用|财务费用|其中.?利息费用|利息收入|加.?其他收益|投资收益.?损失以.?.?.?号填列0|信用减值损失|资产减值损失|资产处置收益|.{0,6}营业利润|.{0,6}营业外收入|.{0,6}营业外支出|利润总额|.{0, 6}所得税费用|净利润|持续经营净利润|归属于母公司股东的净利润|综合收益总额|归属于母公司所有者的综合收益总额|.*收益总额|.*每股收益|.{0,6}基本每股收益|.*价值变动|.{0,6}稀释每股收益|取得子公司及其他营业单位支付的现金净额|营业税金及附加|。*公允价值变动收益|投资收益|对联营企业和合营企业的投资收益|.*其他综合收益|.{0,6}非流动资产处置损失|.{0,6}利润总额|.{0,6}所得税费用|.{0,6}净利润|.{0,6}每股收益.{0,6}"
	incomeStatementTcRegStr = "合併利潤表.?損益表.?收益表|已賺保費|科目|營業總收入|減.?營業成本|.*營業收入|營業總成本|.*營業成本|稅金及附加|銷售費用|管理費用|研發費用|財務費用|其中.?利息費用|利息收入|加.?其他收益|投資收益.?損失以.?.?.?號填列0|信用減值損失|資產減值損失|資產處置收益|.{0,6}營業利潤|.{0,6}營業外收入|.{0,6}營業外支出|利潤總額|.{0, 6}所得稅費用|淨利潤|持續經營淨利潤|歸屬於母公司股東的淨利潤|綜合收益總額|歸屬於母公司所有者的綜合收益總額|.*收益總額|.*每股收益|.{0,6}基本每股收益|.價值變動|.{0,6}稀釋每股收益|取得子公司及其他營業單位支付的現金淨額|營業稅金及附加|。公允價值變動收益|投資收益|對聯營企業和合營企業的投資收益|.*其他綜合收益|.{0,6}非流動資產處置損失|.{0,6}利潤總額|.{0,6}所得稅費用|.{0,6}淨利潤|.{0,6}每股收益.{0,6}"
)

const (
	//资产负债表
	balanceSheet = 1
	//现金流量表
	cashFlow = 2
	//利润表
	income = 3
)

const (
	// 文本块
	tableTypePlain = "plain"
	// 无表格线表格
	tableTypeTableWithoutLine = "table_without_line"
	// 带表格线表格
	tableTypeTableWithLine = "table_with_line"
)

//获取该页的财报类别（1资产负债表，2现金流量表，3利润表）
func getTableType(ocr RecognizeTableResultTable) int {
	var statementType = 0
	var balanceSheetMatchScore, cashFlowStatementMatchScore, incomeStatementMatchScore = 0, 0, 0
	/*
		1、判断表格类型
		2、根据单元格里的line属性判断该单元格属性是否为科目，如果为科目则去科目列表里匹配是哪个科目
		3、整理所有科目单元格后面的列数据
	*/
	//遍历非表格区域的行数据
	// 添加对繁体科目的支持
	balanceSheetRegStr := fmt.Sprintf("%s|%s", balanceSheetRegStr, balanceSheetTcRegStr)
	cashFlowStatementRegStr := fmt.Sprintf("%s|%s", cashFlowStatementRegStr, cashFlowStatementTcRegStr)
	incomeStatementRegStr := fmt.Sprintf("%s|%s", incomeStatementRegStr, incomeStatementTcRegStr)

	balanceSheetReg := regexp.MustCompile(balanceSheetRegStr)
	cashFlowStatementReg := regexp.MustCompile(cashFlowStatementRegStr)
	incomeStatementReg := regexp.MustCompile(incomeStatementRegStr)

	for _, line := range ocr.Lines {
		// 将所有的行ocr数据进行正则匹配，得到财务三大报表的可能值
		if balanceSheetReg.MatchString(line.Text) {
			balanceSheetMatchScore += 1
		}
		if cashFlowStatementReg.MatchString(line.Text) {
			cashFlowStatementMatchScore += 1
		}
		if incomeStatementReg.MatchString(line.Text) {
			incomeStatementMatchScore += 1
		}
	}
	for _, cell := range ocr.TableCells {
		for _, line := range cell.Lines {
			// 将所有的行ocr数据进行正则匹配，得到财务三大报表的可能值
			if balanceSheetReg.MatchString(line.Text) {
				balanceSheetMatchScore += 1
			}
			if cashFlowStatementReg.MatchString(line.Text) {
				cashFlowStatementMatchScore += 1
			}
			if incomeStatementReg.MatchString(line.Text) {
				incomeStatementMatchScore += 1
			}
		}
	}
	limit := 3 // 至少有几个科目列满足时，可以得到表类型
	//根据财务三大报表的可能值得到表格类型，确定是其中的哪一种财务报表
	if balanceSheetMatchScore >= limit && balanceSheetMatchScore >= cashFlowStatementMatchScore && balanceSheetMatchScore >= incomeStatementMatchScore {
		statementType = balanceSheet
	} else if cashFlowStatementMatchScore >= limit && cashFlowStatementMatchScore >= balanceSheetMatchScore && cashFlowStatementMatchScore >= incomeStatementMatchScore {
		statementType = cashFlow
	} else if incomeStatementMatchScore >= limit && incomeStatementMatchScore >= balanceSheetMatchScore && incomeStatementMatchScore >= cashFlowStatementMatchScore {
		statementType = income
	}

	return statementType
}

func getStatement(tableType int, ocr RecognizeTableResultTable, titles *[]StandardStatementTitle, periods *[]Period) (s task.Statement, p []Period) {
	statementTitles := make([]StandardStatementTitle, 0)
	for _, title := range *titles {
		if tableType == int(title.Type) {
			statementTitles = append(statementTitles, title)
		}
	}

	//用于临时存储组织数据，便于后续的日期数据找到对应的组织
	institutions := make([]Institution, 0)
	//用于临时存储时间维度的数据，便于后续金额数据找到各自的时间维度

	statement := task.Statement{
		Count:   0,
		Images:  []task.Image{},
		Headers: []task.Header{},
		Titles:  []task.StatementTitle{},
	}

	title := task.StatementTitle{
		Id:         0,
		ExternalId: "",
		Key: task.StatementTitleKey{
			Ocr:      "",
			Inferred: 0,
			Position: []int{},
		},
		Values: []task.StatementTitleValue{},
	}

	rowIndex := 0
	for _, cell := range ocr.TableCells {
		/*
			遍历每个表格里的行，确认该单元格的属性（科目、金额、时间维度、组织）
			科目：
				1、该行若存在科目，则加入到Titles集合中
				2、继续遍历后续的该行的单元格，如果遇到金额，则加入到上述的科目 values 集合中
				科目匹配规则：
					1、所有科目文字相似度计算取最高
					2、记录科目匹配置信度最高的科目id，并附带科目匹配置信度
			金额：
				如果遇到金额，则加入到上述的科目 values 集合中
			时间维度：
				1、若该行存在时间维度数据，加入到headers集合中，加入的时候需要根据列信息寻找匹配的集团信息
			组织：
				1、遇到组织数据后，加入到组织集合，记录它的列范围
		*/
		if cell.Text == "向中央银行借款" {
			fmt.Print("12")
		}

		if rowIndex != cell.StartRow {
			//新的一个title
			if len(title.Key.Position) > 0 {
				statement.Titles = append(statement.Titles, title)
			}
			title = task.StatementTitle{
				Id:         0,
				ExternalId: "",
				PageIndex:  0,
				Key: task.StatementTitleKey{
					Ocr:      "",
					Inferred: 0,
					Position: []int{},
				},
				Values: []task.StatementTitleValue{},
			}
			rowIndex = cell.StartRow
		}
		switch cell.Semantic {
		case RecognizeTableResultSemanticTypePeriod:
			organization := ""
			// 新增时间维度时，查找属于哪个组织的
			for _, institution := range institutions {
				if cell.StartCol >= institution.StartCol && cell.EndCol <= institution.EndCol {
					organization = institution.Text
					break
				}
			}

			newPeriod := true
			for _, period := range *periods {
				if period.Text == cell.Text && period.Type == 1 && period.StartCol == cell.StartCol && period.EndCol == cell.EndCol {
					newPeriod = false
					break
				}
			}
			if len(*periods) == 0 || newPeriod {
				*periods = append(*periods, Period{
					Text:     cell.Text,
					Type:     1,
					StartCol: cell.StartCol,
					EndCol:   cell.EndCol,
				})
			}

			existed := false
			for _, hander := range statement.Headers {
				if hander.Date == cell.Text && hander.Organization == organization {
					existed = true
				}
			}
			if !existed {
				statement.Headers = append(statement.Headers, task.Header{
					Organization: organization,
					Date:         cell.Text,
					Order:        cell.Order,
				})
			}

		case RecognizeTableResultSemanticTypeInstitution:
			// 组织机构
			institutions = append(institutions, Institution{
				Text:     cell.Text,
				Type:     1,
				StartCol: cell.StartCol,
				EndCol:   cell.EndCol,
			})
		case RecognizeTableResultSemanticTypeItem:
			bastMatchTitle := matchTitle(cell.Text, statementTitles)
			if len(title.Key.Position) > 0 {
				// 暂时先通过判断title的key坐标信息是否有值，来判断该行是否已经找到过title，如果找到过，则新起一个title
				if len(title.Key.Ocr) != 0 {
					statement.Titles = append(statement.Titles, title)
				}
				title = task.StatementTitle{
					Id:         0,
					ExternalId: "",
					PageIndex:  0,
					Key: task.StatementTitleKey{
						Ocr:      "",
						Inferred: 0,
						Position: []int{},
					},
					Values: []task.StatementTitleValue{},
				}
			}
			title.Id = bastMatchTitle.ID
			title.TitleName = bastMatchTitle.Name
			title.Similarity = bastMatchTitle.Score
			// title.PageIndex = pageIndex
			title.RowIndex = cell.StartRow
			title.ColIndex = cell.EndCol
			title.Key = task.StatementTitleKey{
				Ocr:      cell.Text,
				Inferred: bastMatchTitle.ID,
				Position: cell.Position,
			}
			title.Values = []task.StatementTitleValue{}
		case RecognizeTableResultSemanticTypeAmount:
			title.Values = append(title.Values, task.StatementTitleValue{
				Ocr:        cell.Text,
				Supervised: cell.Text,
				Position:   cell.Position,
			})
		default:
			if cell.StartRow >= rowIndex && cell.EndRow <= rowIndex {
				for _, item := range *periods {
					if cell.StartCol >= item.StartCol && cell.EndCol <= item.EndCol {
						title.Values = append(title.Values, task.StatementTitleValue{
							Ocr:        cell.Text,
							Supervised: cell.Text,
							Position:   cell.Position,
						})
					}
				}
			}
		}
	}

	if len(title.Key.Position) > 0 && len(title.Key.Ocr) > 0 {
		statement.Titles = append(statement.Titles, title)
	}

	if len(statement.Titles) > 0 {
		titleColCounts := make(map[int]int, 0)
		for _, t := range statement.Titles {
			if len(t.Values) > 0 {
				titleColCounts[len(t.Values)] += 1
			}
		}
		keys := make([]int, 0)
		for key, _ := range titleColCounts {
			keys = append(keys, key)
		}
		sort.SliceStable(keys, func(i, j int) bool {
			return titleColCounts[keys[i]] > titleColCounts[keys[j]]
		})
		var maxCount int
		if len(keys) > 0 {
			maxCount = keys[0]
		}

		for i := 0; i < len(statement.Titles); i++ {
			// t := statement.Titles[i]
			if len(statement.Titles[i].Values) > maxCount {
				statement.Titles[i].Values = statement.Titles[i].Values[:maxCount]
			}
			if statement.Count <= len(statement.Titles[i].Values) {
				statement.Count = len(statement.Titles[i].Values)
			}
		}
	}

	return statement, nil
}

var charList = []string{
	"△", "#", "*", "※", "：", ":", "（", "）", "(", ")", "☆", "★", "▲", "▼", "√", "×",
}

func matchTitle(ocr string, titles []StandardStatementTitle) StandardStatementTitle {
	tempTitle := StandardStatementTitle{
		ID:      0,
		Name:    "",
		Aliases: "",
		Type:    0,
		Score:   0,
	}
	for _, char := range charList {
		ocr = strings.ReplaceAll(ocr, char, "")
	}
	for _, title := range titles {
		aliases := []string{}
		if len(title.Aliases) > 0 {
			json.Unmarshal([]byte(title.Aliases), &aliases)
		}
		aliases = append(aliases, title.Name)
		for i := range aliases {
			for _, char := range charList {
				aliases[i] = strings.ReplaceAll(aliases[i], char, "")
			}
		}
		tempScore := strsim.FindBestMatch(ocr, aliases)
		if tempTitle.Score < tempScore.Match.Score && tempScore.Match.Score > 0 && tempScore.Match.Score > 0.5 {
			tempTitle = title
			tempTitle.Score = tempScore.Match.Score
		}
	}
	return tempTitle
}
