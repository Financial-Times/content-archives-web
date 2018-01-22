package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/client/metadata"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	region       = "REGION"
	bucket       = "BUCKET"
	bucketPrefix = "BUCKET_PREFIX"
	fileName     = "content_archive.zip"
	testURL      = "https://ft.com"
)

var (
	key                = "file1"
	lastModified       = time.Date(2017, time.October, 17, 10, 51, 19, 0, time.UTC)
	size         int64 = 23442
	content            = s3.Object{
		Key:          &key,
		LastModified: &lastModified,
		Size:         &size,
	}
)

type s3ReaderMock struct {
	mock.Mock
}

func (m *s3ReaderMock) Read(input *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {
	args := m.Called(input)
	if result, ok := args.Get(0).(*s3.ListObjectsOutput); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *s3ReaderMock) GetDownloadURL(input *s3.GetObjectInput) (*request.Request, error) {
	args := m.Called(input)
	if result, ok := args.Get(0).(*request.Request); ok {
		return result, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestRetrieveArchivesFromS3Success(t *testing.T) {
	s3ReaderObj := new(s3ReaderMock)
	s3Result := &s3.ListObjectsOutput{
		Contents: []*s3.Object{&content},
	}
	s3ReaderObj.On("Read", mock.AnythingOfType("*s3.ListObjectsInput")).Return(s3Result, nil)

	s3service := service{region, bucket, bucketPrefix, s3ReaderObj}
	zipFiles, err := s3service.RetrieveArchivesFromS3()

	assert.NoError(t, err)
	assert.Equal(t, 1, len(zipFiles))
	assert.Equal(t, key, zipFiles[0].Name)
	assert.Equal(t, "22.9K", zipFiles[0].Size)
	assert.Equal(t, "2017-10-17 10:51:19", zipFiles[0].LastModified)
	s3ReaderObj.AssertExpectations(t)
}

func TestRetrieveArchivesFromS3WithError(t *testing.T) {
	s3ReaderObj := new(s3ReaderMock)
	s3ReaderObj.On("Read", mock.AnythingOfType("*s3.ListObjectsInput")).Return(&s3.ListObjectsOutput{}, fmt.Errorf("new error"))

	s3service := service{region, bucket, bucketPrefix, s3ReaderObj}
	zipFiles, err := s3service.RetrieveArchivesFromS3()
	assert.Error(t, err)
	assert.Equal(t, 0, len(zipFiles))
	s3ReaderObj.AssertExpectations(t)
}

func TestDownloadArchiveFromS3Success(t *testing.T) {
	s3ReaderObj := new(s3ReaderMock)
	fn := func(r *request.Request) error { return nil }
	newReq, _ := http.NewRequest(http.MethodGet, testURL, nil)

	req := &request.Request{
		ClientInfo:  metadata.ClientInfo{Endpoint: testURL},
		Operation:   &request.Operation{BeforePresignFn: fn},
		HTTPRequest: newReq,
	}
	s3ReaderObj.On("GetDownloadURL", mock.AnythingOfType("*s3.GetObjectInput")).Return(req, nil)

	s3service := service{region, bucket, bucketPrefix, s3ReaderObj}
	result, err := s3service.GetDownloadURLForFile(fileName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	s3ReaderObj.AssertExpectations(t)
}

func TestDownloadArchiveFromS3WithError(t *testing.T) {
	s3ReaderObj := new(s3ReaderMock)
	s3ReaderObj.On("GetDownloadURL", mock.AnythingOfType("*s3.GetObjectInput")).Return(nil, fmt.Errorf("new error"))

	s3service := service{region, bucket, bucketPrefix, s3ReaderObj}
	result, err := s3service.GetDownloadURLForFile(fileName)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	s3ReaderObj.AssertExpectations(t)
}
