package main

import (
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Reader interface for read/download files from AWS S3
type S3Reader interface {
	Read(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	GetDownloadURL(input *s3.GetObjectInput) (*request.Request, error)
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

func (r *reader) GetDownloadURL(input *s3.GetObjectInput) (*request.Request, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(input)
	return req, nil
}
