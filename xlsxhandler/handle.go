package xlsxhandler

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"path/filepath"
)

//XlsxData interface has methods that we need to open and dump data from an xcell file .
type XlsxData interface {
	Open() error
	Dump()
}

// XlsxFileInfo struct with instance that we need as variables to pass to XlsxData interface methods to work.
type XlsxFileInfo struct {
	FileName  string
	SheetName string
	File      *excelize.File
	Chan      chan ChanResult
	Wg        *sync.WaitGroup
	Reader    io.Reader
	CntNumber int
}

// ChanResult result that the channel receives
type ChanResult struct {
	DumpData [][]string
	Err      error
}

func checkExcelFile(filename string) (string, bool) {
	if ext := filepath.Ext(filename); ext != ".xlsx" {
		return ext, false
	}
	return "", true
}

// New function returns pointer to XlsxFileInfo struct
func New(filename, sheetName string, c chan ChanResult, wg *sync.WaitGroup, reader io.Reader, contactNumber int) (XlsxData, error) {
	var xlxData XlsxData
	x := &XlsxFileInfo{}
	if ext, isExcel := checkExcelFile(filename); isExcel {
		x.FileName = filename
	} else {
		err := fmt.Errorf("Want <.xlsx> file got %s\n", ext)
		return nil, err
	}
	x.SheetName = sheetName
	x.File = nil
	x.Chan = c
	x.Wg = wg
	x.Reader = reader
	x.CntNumber = contactNumber
	xlxData = x
	return xlxData, nil
}

// Open open xcel file
func (xl *XlsxFileInfo) Open() (err error) {
	isWeb := os.Getenv("web")
	var f *excelize.File
	if isWeb != "true" {
		f, err = excelize.OpenFile(xl.FileName)
	} else {
		f, err = excelize.OpenReader(xl.Reader)
	}
	xl.File = f
	return
}

func checkError(err error, ch chan ChanResult, cR *ChanResult) {
	if err != nil {
		cR.Err, cR.DumpData = err, nil
		ch <- *cR
	}
}

// Dump get Data from excel file
func (xl *XlsxFileInfo) Dump() {
	c := new(ChanResult)
	defer xl.Wg.Done()
	rows, err := xl.File.Rows(xl.SheetName)
	checkError(err, xl.Chan, c)
	var counter int = 0
	fmt.Println(xl.CntNumber)
	for rows.Next() {
		counter++
		if xl.CntNumber != 0 {
			if counter == xl.CntNumber+1 {
				break
			}
		}
		row, err := rows.Columns()
		checkError(err, xl.Chan, c)
		c.DumpData = append(c.DumpData, row)
	}
	c.Err = nil
	xl.Chan <- *c
}
