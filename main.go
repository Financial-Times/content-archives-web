package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatal("$AWS_REGION must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")
	router.SetFuncMap(template.FuncMap{
		"downloadArchiveFromS3": downloadArchiveFromS3,
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"zipFiles": retrieveArchivesFromS3(awsRegion),
		})
	})

	router.Run(":" + port)
}

func retrieveArchivesFromS3(awsRegion string) []string {
	sess, err := session.NewSession()

	// Create S3 service client
	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Prefix: aws.String("yearly-archives"),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		log.Printf("Unable to list buckets, %v", err)
	}
	log.Println(result)

	zipFiles := make([]string, 0)
	for _, content := range result.Contents {
		zipFiles = append(zipFiles, aws.StringValue(content.Key))
	}

	return zipFiles
}

func downloadArchiveFromS3(fileName string) {
	sess, err := session.NewSession()
	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(fileName)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", err)
	}
	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("BUCKET_NAME")),
			Key:    aws.String(fileName),
		})

	if err != nil {
		exitErrorf("Unable to download item %q, %v", fileName, err)
	}

	//_, err = io.Copy(fileName, resp.Body)

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
}
