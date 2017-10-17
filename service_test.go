package main

import (
	"testing"

	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	region       = "REGION"
	bucket       = "BUCKET"
	bucketPrefix = "BUCKET_PREFIX"
	fileName     = "content_archive.zip"
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

func (m *s3ReaderMock) Download(input *s3.GetObjectInput) (*aws.WriteAtBuffer, error) {
	args := m.Called(input)
	if result, ok := args.Get(0).(*aws.WriteAtBuffer); ok {
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
	a := []byte{'a', 'b', 'c'}
	fileBuffer := aws.NewWriteAtBuffer(a)
	s3ReaderObj.On("Download", mock.AnythingOfType("*s3.GetObjectInput")).Return(fileBuffer, nil)

	s3service := service{region, bucket, bucketPrefix, s3ReaderObj}
	result, err := s3service.DownloadArchiveFromS3(fileName)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 3, len(result))
	s3ReaderObj.AssertExpectations(t)
}

func TestDownloadArchiveFromS3WithError(t *testing.T) {
	s3ReaderObj := new(s3ReaderMock)
	s3ReaderObj.On("Download", mock.AnythingOfType("*s3.GetObjectInput")).Return(nil, fmt.Errorf("new error"))

	s3service := service{region, bucket, bucketPrefix, s3ReaderObj}
	result, err := s3service.DownloadArchiveFromS3(fileName)

	assert.Error(t, err)
	assert.Nil(t, result)
	s3ReaderObj.AssertExpectations(t)
}
