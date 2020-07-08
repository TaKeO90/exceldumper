package xlsxhandler

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

//TODO: be able to read xcell file and return data structure you need to transform to a csv file

type XlsxData interface {
	Open() error
	Dump()
}

var (
	csvRow [][]string
)

type XlsxFileInfo struct {
	FileName  string
	SheetName string
	File      *excelize.File
	Chan      chan ChanResult
}

type ChanResult struct {
	CsvData [][]string
	Err     error
}

func New(filename, sheetName string, c chan ChanResult) *XlsxFileInfo {
	x := &XlsxFileInfo{}
	x.FileName = filename
	x.SheetName = sheetName
	x.File = nil
	x.Chan = c
	return x
}

func dataRefactoring(data [][]string, index int) []string {
	var r []string
	for i := range data {
		r = append(r, data[i][index])
	}
	return r
}

func checkLength(data [][]string) (bool, int) {
	var less bool
	l := 0
	for j := range data {
		if j == len(data)-1 {
			break
		}
		if len(data[j]) != len(data[j+1]) {
			l = j + 1
			less = true
		}
	}
	return less, l + 1
}

func (xl *XlsxFileInfo) Open() error {
	f, err := excelize.OpenFile(xl.FileName)
	xl.File = f
	if err != nil {
		return err
	}
	return nil
}

func (xl *XlsxFileInfo) Dump() {
	var totalRows [][]string
	cols, err := xl.File.Cols(xl.SheetName)
	c := &ChanResult{}
	if err != nil {
		c.CsvData, c.Err = csvRow, err
		xl.Chan <- *c
	}
	for cols.Next() {
		col, err := cols.Rows()
		c.CsvData, c.Err = csvRow, err
		if err != nil {
			xl.Chan <- *c
		}
		var r []string
		for i, rowcell := range col {
			if i == 0 {
				r = r[:0]
			}
			r = append(r, rowcell)
		}
		totalRows = append(totalRows, r)
	}
	if less, l := checkLength(totalRows); less {
		err := fmt.Errorf("The following column number %d has less data than others", l)
		c.CsvData, c.Err = csvRow, err
	} else {
		for j := range totalRows[l] {
			d := dataRefactoring(totalRows, j)
			csvRow = append(csvRow, d)
		}
	}
	c.CsvData, c.Err = csvRow, nil
	xl.Chan <- *c
}
