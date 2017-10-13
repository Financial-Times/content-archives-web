package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Financial-Times/ft-s3o-go/s3o"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatal("$AWS_REGION must be set")
	}

	awsBucketName := os.Getenv("AWS_BUCKET_NAME")
	if awsBucketName == "" {
		log.Fatal("$AWS_BUCKET_NAME must be set")
	}

	awsBucketPrefix := os.Getenv("AWS_BUCKET_PREFIX")
	if awsBucketPrefix == "" {
		log.Fatal("$AWS_BUCKET_PREFIX must be set")
	}

	s3Reader := NewS3Reader(awsRegion, awsBucketName, awsBucketPrefix)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.GET("/", homepageHandler(s3Reader))
	router.GET("/download/:prefix/:name", downloadHandler(s3Reader))

	handler := s3o.Handler(router)
	http.ListenAndServe(":"+port, handler)
}

func homepageHandler(s3Reader S3Reader) func(c *gin.Context) {
	return func(c *gin.Context) {
		zipFiles, err := s3Reader.RetrieveArchivesFromS3()
		if err != nil {
			c.HTML(http.StatusInternalServerError, "Unable to get archives list from S3", nil)
		}
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"zipFiles": zipFiles,
		})
	}
}

func downloadHandler(s3Reader S3Reader) func(c *gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("name")
		bytes, err := s3Reader.DownloadArchiveFromS3(name)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "Unable to download archive from S3", nil)
		}

		c.Header("Content-Disposition", "attachment; filename="+name)
		c.Data(http.StatusOK, "application/zip", bytes)
	}
}
