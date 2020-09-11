package handler

import (
	"strings"

	"github.com/TaKeO90/exceldumper/server/dialer"
	"github.com/gin-gonic/gin"
)

func checkPath(path string) map[string]string {
	m := make(map[string]string)
	fstE, sndE := strings.TrimSpace(strings.Split(path, "/")[1]), strings.TrimSpace(strings.Split(path, "/")[2])
	m[fstE] = sndE
	return m
}

// HandleRequest handler function.
func HandleRequest(c *gin.Context) {
	path := c.Request.URL.Path
	m := checkPath(path)
	switch m["file"] {
	case "upload":
		dialer.UploadFile(c)
	case "download":
		dialer.DownloadFile(c)
	}
}
