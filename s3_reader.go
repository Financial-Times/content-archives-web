package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Reader interface for read/download files from AWS S3
type S3Reader interface {
	Read(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	Download(input *s3.GetObjectInput) (*aws.WriteAtBuffer, error)
}

type reader struct {
}

// NewS3Reader returns a new S3Reader instance
func NewS3Reader() S3Reader {
	return &reader{}
}

func (r *reader) Read(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	return svc.ListObjects(input)
}

func (r *reader) Download(input *s3.GetObjectInput) (*aws.WriteAtBuffer, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	downloader := s3manager.NewDownloader(sess)
	var a []byte
	buffer := aws.NewWriteAtBuffer(a)
	_, err = downloader.Download(buffer, input)
	return buffer, err
}
