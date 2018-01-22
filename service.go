package main

import (
	"log"
	"time"

	"code.cloudfoundry.org/bytefmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Service interface responsible for read/download operations using S3Reader interface
type S3Service interface {
	RetrieveArchivesFromS3() ([]ZipFile, error)
	GetDownloadURLForFile(fileName string) (string, error)
}

type service struct {
	awsRegion       string
	awsBucketName   string
	awsBucketPrefix string
	s3Reader        S3Reader
}

// ZipFile object that holds details about an archive
type ZipFile struct {
	Name         string
	Size         string
	LastModified string
}

// NewS3Service returns a new S3Service instance
func NewS3Service(awsRegion string, awsBucketName string, awsBucketPrefix string) S3Service {
	s3Reader := NewS3Reader()
	return &service{awsRegion, awsBucketName, awsBucketPrefix, s3Reader}
}

// RetrieveArchivesFromS3 returns a list of objects from S3, based on awsBucketName and awsBucketPrefix
func (s service) RetrieveArchivesFromS3() ([]ZipFile, error) {
	zipFiles := make([]ZipFile, 0)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(s.awsBucketName),
		Prefix: aws.String(s.awsBucketPrefix),
	}

	result, err := s.s3Reader.Read(input)
	if err != nil {
		return zipFiles, err
	}

	for _, content := range result.Contents {
		lastModified := aws.TimeValue(content.LastModified)
		lastModifiedStr := lastModified.Format("2006-01-02 15:04:05")
		size := aws.Int64Value(content.Size)
		sizeMB := bytefmt.ByteSize(uint64(size))

		file := ZipFile{aws.StringValue(content.Key), sizeMB, lastModifiedStr}
		zipFiles = append(zipFiles, file)
	}

	return zipFiles, nil
}

// GetDownloadURLForFile returns an download URL for an AWS S3 file
func (s service) GetDownloadURLForFile(fileName string) (string, error) {
	log.Printf("Getting download URL for file: %s", fileName)
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.awsBucketName),
		Key:    aws.String(s.awsBucketPrefix + "/" + fileName),
	}

	req, err := s.s3Reader.GetDownloadURL(input)
	if err != nil {
		log.Printf("Unable to download item %q, %v", fileName, err)
		return "", err
	}

	url, err := req.Presign(5 * time.Minute)
	if err != nil {
		log.Printf("Unable to get download URL for file: %q, %v", fileName, err)
		return "", err
	}

	return url, nil
}
