package csvhandler

import (
	"encoding/csv"
	"os"
	"sync"
)

//TODO: progress bar or status

type CsvFile interface {
	WriteData()
}

type CsvFileData struct {
	CsvWriter *csv.Writer
	FileData  [][]string
	Chan      chan ChanRes
	Wg        *sync.WaitGroup
}

type ChanRes struct {
	IsOk bool
	Err  error
}

func New(filename string, data [][]string, c chan ChanRes, wg *sync.WaitGroup) (*CsvFileData, error) {
	csvF := new(CsvFileData)
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	w := csv.NewWriter(f)
	csvF.CsvWriter, csvF.FileData = w, data
	csvF.Chan = c
	csvF.Wg = wg
	return csvF, nil
}

func (cF *CsvFileData) WriteData() {
	var ok bool
	chanR := new(ChanRes)
	defer cF.Wg.Done()
	for i := range cF.FileData {
		if err := cF.CsvWriter.Write(cF.FileData[i]); err == nil {
			ok = true
			chanR.Err, chanR.IsOk = nil, ok
		} else if err := cF.CsvWriter.Error(); err != nil {
			chanR.Err = err
			chanR.IsOk = false
			cF.Chan <- *chanR
		} else {
			ok = false
			chanR.Err, chanR.IsOk = err, ok
			cF.Chan <- *chanR
		}
		cF.CsvWriter.Flush()
	}
	cF.Chan <- *chanR
}
