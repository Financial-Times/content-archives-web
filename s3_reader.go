package main

import (
	"log"

	"code.cloudfoundry.org/bytefmt"
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

// ZipFile object that holds details about an archive
type ZipFile struct {
	Name         string
	Size         string
	LastModified string
}

// NewS3Reader returns a new S3Reader instance
func NewS3Reader(awsRegion string, awsBucketName string, awsBucketPrefix string) S3Reader {
	return S3Reader{awsRegion, awsBucketName, awsBucketPrefix}
}

// RetrieveArchivesFromS3 returns a list of objects from S3, based on awsBucketName and awsBucketPrefix
func (s3Reader *S3Reader) RetrieveArchivesFromS3() ([]ZipFile, error) {
	zipFiles := make([]ZipFile, 0)
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
		lastModifiedStr := lastModified.Format("2006-01-02 15:04:05")
		size := aws.Int64Value(content.Size)
		sizeMB := bytefmt.ByteSize(uint64(size))

		file := ZipFile{aws.StringValue(content.Key), sizeMB, lastModifiedStr}
		zipFiles = append(zipFiles, file)
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
