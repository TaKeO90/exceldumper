package service

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

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

type ChanResult struct {
	File string
	Ok   bool
	Err  error
}

//(file string, ok bool, err error)
//DumpExcelData function accept a excel
func DumpExcelData(reader io.Reader, outFName string, cntNumber int, c chan ChanResult) {
	chanResult := new(ChanResult)
	filename := strings.Split(outFName, ".")[0] + OUTFILEEXT
	x, err := xlsxhandler.New(outFName, SHEETNAME, reader, cntNumber)
	err = x.Open()
	data, err := x.Dump()
	if err != nil {
		chanResult.Err = err
		c <- *chanResult
	}
	v, err := vcfhandler.NewVcf(data, filename)
	ok := v.ExtWrite()
	file := filename
	chanResult.File, chanResult.Ok = file, ok
	c <- *chanResult
}

func checkToDownloadFile(filename string) (bool, error) {
	dir, err := ioutil.ReadDir(".")
	if err != nil {
		return false, err
	}
	for _, f := range dir {
		if f.Name() == filename {
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
