package dialer

import (
	"fmt"
	"net/http"
	"path/filepath"
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

type BadRequest struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func handleOption(c *gin.Context) {
	c.Header("Allow", "POST, OPTIONS")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Content-Length")
	c.Header("Content-Type", "application/json")
	if c.Request.Method == "OPTIONS" {
		c.Status(http.StatusOK)
	}
}

// UploadFile get the uploaded file by the user.
func UploadFile(c *gin.Context) {
	handleOption(c)
	if c.Request.Method == "POST" {
		ch := make(chan service.ChanResult)
		c.Request.ParseMultipartForm(14)
		contactN := c.Request.FormValue("cntNumber")
		// get file from http request.
		file, fileHeader, err := c.Request.FormFile("file")
		errHandler(err, c)
		fExt := filepath.Ext(fileHeader.Filename)
		defer file.Close()
		errHandler(err, c)
		// convert string to int cause we need an int as number of element that we need to know
		// the number of contacts that we need to convert to a vcf file.
		cntN, err := strconv.Atoi(contactN)
		errHandler(err, c)
		if file != nil && fExt == ".xlsx" {
			go service.DumpExcelData(file, fileHeader.Filename, cntN, ch)
			results := <-ch
			errHandler(results.Err, c)
			errHandler(err, c)
			if results.Ok {
				c.JSON(http.StatusCreated, &VcfCreated{true, results.File})
			}
		} else {
			c.JSON(http.StatusBadRequest, &BadRequest{false, "file type not correct"})
		}
	}
}

//DownloadFile let the user Download the output file.
func DownloadFile(c *gin.Context) {
	handleOption(c)
	if c.Request.Method == "POST" {
		c.Request.ParseMultipartForm(0)
		filename := c.Request.FormValue("filename")
		err := service.GetFile(filename, c)
		errHandler(err, c)
	}
}
