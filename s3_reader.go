package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Reader struct containig AWS configuration
type S3Reader struct {
	awsRegion       string
	awsBucketName   string
	awsBucketPrefix string
}

// NewS3Reader returns a new S3Reader instance
func NewS3Reader(awsRegion string, awsBucketName string, awsBucketPrefix string) S3Reader {
	return S3Reader{awsRegion, awsBucketName, awsBucketPrefix}
}

// RetrieveArchivesFromS3 returns a list of objects from S3, based on awsBucketName and awsBucketPrefix
func (s3Reader *S3Reader) RetrieveArchivesFromS3() (map[string]string, error) {
	//zipFiles := make([]string, 0)
	zipFiles := make(map[string]string)
	sess, err := session.NewSession()
	if err != nil {
		return zipFiles, err
	}

	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s3Reader.awsBucketName),
		Prefix: aws.String(s3Reader.awsBucketPrefix),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		return zipFiles, err
	}

	for _, content := range result.Contents {
		lastModified := aws.TimeValue(content.LastModified)
		zipFiles[aws.StringValue(content.Key)] = lastModified.Format("2006-01-02 15:04:05")
	}

	return zipFiles, nil
}

// DownloadArchiveFromS3 downloads a file from AWS S3 and returns an array of bytes
func (s3Reader *S3Reader) DownloadArchiveFromS3(fileName string) ([]byte, error) {
	log.Println("Starting to download archives...")

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)
	var a []byte
	buffer := aws.NewWriteAtBuffer(a)
	_, err = downloader.Download(buffer,
		&s3.GetObjectInput{
			Bucket: aws.String(s3Reader.awsBucketName),
			Key:    aws.String(s3Reader.awsBucketPrefix + "/" + fileName),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", fileName, err)
		return nil, err
	}

	return buffer.Bytes(), nil
}
