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
	handler := NewHandler(s3Reader)
	router := createRouting(handler)
	s3AuthHandler := s3o.Handler(router)

	http.ListenAndServe(":"+port, s3AuthHandler)
}

func createRouting(handler Handler) http.Handler {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.GET("/", handler.HomepageHandler())
	router.GET("/download/:prefix/:name", handler.DownloadHandler())

	return router
}
