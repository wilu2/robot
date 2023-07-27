package excel

import (
	"bytes"
	"financial_statement/pkg/log"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

type TableCell struct {
	StartRow int    `json:"start_row"`
	StartCol int    `json:"start_col"`
	EndRow   int    `json:"end_row"`
	EndCol   int    `json:"end_col"`
	Text     string `json:"text"`
	// Semantic string `json:"semantic"`
}

type Table struct {
	TableCells []TableCell
}

type mergeCellStruct struct {
	startRowIndex int
	endRowIndex   int
	startColIndex int
	endColIndex   int
}

func GetPages(excel []byte) ([]Table, error) {
	fileRead := bytes.NewReader(excel)
	file, err := excelize.OpenReader(fileRead)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err != nil {
		return nil, err
	}
	tables := make([]Table, 0)
	sheets := file.GetSheetList()

	for _, sheet := range sheets {
		//获取合并的单元格
		mergeCellList := make([]mergeCellStruct, 0)
		mergeCells, err := file.GetMergeCells(sheet)
		for _, mergeCell := range mergeCells {
			value := mergeCell.GetCellValue()
			startAxis := mergeCell.GetStartAxis()
			endAxis := mergeCell.GetEndAxis()
			startColIndex, startRowIndex, _ := excelize.CellNameToCoordinates(startAxis)
			endColIndex, endRowIndex, _ := excelize.CellNameToCoordinates(endAxis)
			mergeCellList = append(mergeCellList, mergeCellStruct{
				startRowIndex: startRowIndex - 1,
				endRowIndex:   endRowIndex - 1,
				startColIndex: startColIndex - 1,
				endColIndex:   endColIndex - 1,
			})
			fmt.Printf("MergeCell value:%s startAxis:%s endAxis:%s startRowIndex:%d endRowIndex:%d startColIndex:%d endColIndex:%d \n", value, startAxis, endAxis, startRowIndex-1, endRowIndex-1, startColIndex-1, endColIndex-1)
		}

		table := Table{}
		rows, err := file.GetRows(sheet)
		if err != nil {
			log.Errorf("Get Rows with error:%s", err.Error())
			continue
		}

		for rowIndex, row := range rows {
			for colIndex, colCell := range row {
				for _, mergeCell := range mergeCellList {
					if rowIndex >= mergeCell.startRowIndex && rowIndex <= mergeCell.startRowIndex &&
						colIndex >= mergeCell.startColIndex && colIndex <= mergeCell.endColIndex {
						//fmt.Printf("MergeCell Text:%s startRowIndex:%d endRowIndex:%d startColIndex:%d endColIndex:%d \n", colCell, mergeCell.startRowIndex, mergeCell.endRowIndex, mergeCell.startColIndex, mergeCell.endColIndex)
						table.TableCells = append(table.TableCells, TableCell{
							StartRow: mergeCell.startRowIndex,
							StartCol: mergeCell.startColIndex,
							EndRow:   mergeCell.endRowIndex,
							EndCol:   mergeCell.endColIndex,
							Text:     strings.TrimSpace(colCell),
						})
						continue
					}
				}
				table.TableCells = append(table.TableCells, TableCell{
					StartRow: rowIndex,
					StartCol: colIndex,
					EndRow:   rowIndex,
					EndCol:   colIndex,
					Text:     strings.TrimSpace(colCell),
				})
			}
		}
		tables = append(tables, table)
	}
	return tables, nil
}

//  ConvertToPDF
//  @Description: 转换文件为pdf
//  @param filePath 需要转换的文件
//  @param outPath 转换后的PDF文件存放目录
//  @return string
func ConvertToPDF(excel []byte) ([]byte, error) {
	fileName := uuid.New().String()
	excelFilePath := filepath.Join("/usr/temp", fileName+".xls")
	outFilePath := filepath.Join("/usr/temp", fileName+".pdf")
	defer func() {
		if err := os.Remove(excelFilePath); err != nil {
			log.Warnf("Excel临时文件%s删除失败:%s", excelFilePath, err.Error())
		}
		if err := os.Remove(outFilePath); err != nil {
			log.Warnf("Pdf临时文件%s删除失败:%s", outFilePath, err.Error())
		}
	}()
	if err := ioutil.WriteFile(excelFilePath, excel, 0666); err != nil {
		return nil, err
	}
	// 1、拼接执行转换的命令
	commandName := "libreoffice"
	params := []string{"--invisible", "--headless", "--convert-to", "pdf:calc_pdf_Export:{\"SinglePageSheets\":{\"type\":\"boolean\",\"value\":\"true\"}}", excelFilePath, "--outdir", "/usr/temp/"}
	// 开始执行转换
	if _, err := interactiveToexec(commandName, params); err == nil {
		return ioutil.ReadFile(outFilePath)
	} else {
		return nil, err
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//
//  interactiveToexec
//  @Description: 执行指定命令
//  @param commandName 命令名称
//  @param params 命令参数
//  @return string 执行结果返回信息
//  @return bool 是否执行成功
//
func interactiveToexec(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	if err != nil {
		log.Fatalf("Error: <", err, "> when exec command read out buffer")
		return "", err
	} else {
		return string(buf), err
	}
}
