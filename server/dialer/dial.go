package dialer

import (
	"net/http"
	"strconv"

	"github.com/TaKeO90/exceldumper/server/service"
	"github.com/gin-gonic/gin"
)

func errHandler(err error, c *gin.Context) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Something went Wrong")
	}
}

// VcfCreated structure that represent the response to the client request.
type VcfCreated struct {
	Ok   bool   `json:"ok"`
	File string `json:"file"`
}

// UploadFile get the uploaded file by the user.
func UploadFile(c *gin.Context) {
	c.Request.ParseMultipartForm(14)
	contactN := c.Request.FormValue("cntNumber")
	// get file from http request.
	file, fileHeader, err := c.Request.FormFile("file")
	errHandler(err, c)
	defer file.Close()
	errHandler(err, c)
	// convert string to int cause we need an int as number of element that we need to know
	// the number of contacts that we need to convert to a vcf file.
	cntN, err := strconv.Atoi(contactN)
	errHandler(err, c)
	fileN, ok, err := service.DumpExcelData(file, fileHeader.Filename, cntN)
	errHandler(err, c)
	if ok {
		// if success we assign to FileToDownload the output vcf file.
		service.FileToDownload = fileN
		c.JSON(http.StatusCreated, &VcfCreated{true, fileN})
	}
}

//DownloadFile let the user Download the output file.
func DownloadFile(c *gin.Context) {
	c.Request.ParseMultipartForm(0)
	filename := c.Request.FormValue("filename")
	err := service.GetFile(filename, c)
	errHandler(err, c)
}
