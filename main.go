package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

	s3Service := NewS3Service(awsRegion, awsBucketName, awsBucketPrefix)
	healthCheck := HealthCheck{awsBucketName, awsBucketPrefix}
	appHandler := NewHandler(s3Service, healthCheck)

	r := mux.NewRouter()
	// load static files
	staticH := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(staticH)

	// use middlewares to restrict access to FT members only
	r.Handle("/", appHandler.AuthHandler(appHandler.HomepageHandler))
	r.Handle("/download/{prefix}/{name}", appHandler.AuthHandler(appHandler.DownloadHandler))

	// health should be accessible for anyone
	r.HandleFunc("/__health", appHandler.HealthCheckHandler)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
