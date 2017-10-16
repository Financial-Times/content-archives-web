package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler struct containig handlers configuration
type Handler struct {
	s3Reader S3Reader
}

// NewHandler returns a new Handler instance
func NewHandler(s3Reader S3Reader) Handler {
	return Handler{s3Reader}
}

// HomepageHandler serves the homepage content
func (h *Handler) HomepageHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		zipFiles, err := h.s3Reader.RetrieveArchivesFromS3()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "Unable to get archives list from S3", nil)
		}
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"zipFiles": zipFiles,
		})
	}
}

// DownloadHandler starts the download of the specified file in request
func (h *Handler) DownloadHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" || len(strings.TrimSpace(name)) == 0 {
			c.HTML(http.StatusBadRequest, "Please specify the name of the file", nil)
		}

		bytes, err := h.s3Reader.DownloadArchiveFromS3(name)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "Unable to download archive from S3", nil)
		}

		c.Header("Content-Disposition", "attachment; filename="+name)
		c.Data(http.StatusOK, "application/zip", bytes)
	}
}
