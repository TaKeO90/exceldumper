package service

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/TaKeO90/exceldumper/vcfhandler"
	"github.com/TaKeO90/exceldumper/xlsxhandler"
	"github.com/gin-gonic/gin"
)

const (
	// SHEETNAME the name of the sheet.
	SHEETNAME string = "Sheet1"
	// OUTFILEEXT the vcf file extension.
	OUTFILEEXT string = ".vcf"
)

// FileToDownload the output file that the user should Download.
// NOTE: if that variable doesn't do the expected work use memcached server to store files that they need to be downloaded and think about each user need
//	a token that is just a encoded name.
var FileToDownload string

//DumpExcelData function accept a excel
func DumpExcelData(reader io.Reader, outFName string, cntNumber int) (file string, ok bool, err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	c := make(chan xlsxhandler.ChanResult)
	cv := make(chan vcfhandler.VcfChanRes)
	filename := strings.Split(outFName, ".")[0] + OUTFILEEXT
	x, err := xlsxhandler.New(outFName, SHEETNAME, c, &wg, reader, cntNumber)
	err = x.Open()
	go x.Dump()
	ext := <-c
	if ext.Err != nil {
		err = ext.Err
	}
	v, err := vcfhandler.NewVcf(ext.DumpData, &wg, cv, filename)
	go v.ExtWrite()
	vcfD := <-cv
	if vcfD.Err != nil {
		err = vcfD.Err
	}
	ok = vcfD.Ok
	file = filename
	wg.Wait()
	return
}

func checkToDownloadFile(filename string) (bool, error) {
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		return false, err
	}
	for _, f := range dir {
		if f.Name() == filename && filename == FileToDownload {
			return true, nil
		}
	}
	return false, nil
}

// GetFile get the file that the user has requested.
func GetFile(filename string, c *gin.Context) error {
	ok, err := checkToDownloadFile(filename)
	if err != nil {
		return err
	}
	if ok {
		openedFile, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer openedFile.Close()
		fileHeader := make([]byte, 512)
		openedFile.Read(fileHeader)
		fileCntType := http.DetectContentType(fileHeader)
		fileStat, err := openedFile.Stat()
		if err != nil {
			return err
		}
		fileSize := strconv.FormatInt(fileStat.Size(), 10)
		c.Header("Content-Disposition", "attachement; filename="+filename)
		c.Header("Content-Type", fileCntType)
		c.Header("Content-Length", fileSize)
		openedFile.Seek(0, 0)
		io.Copy(c.Writer, openedFile)
	}
	return nil
}
