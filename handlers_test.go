package main

import (
	"github.com/stretchr/testify/mock"
)

import "testing"

var s3Reader S3Reader = NewS3Reader("REGION", "BUCKET", "BUCKET_PREFIX")

type s3ReaderMock struct {
	mock.Mock
}

func newS3ReaderMock() S3Reader {
	return &s3ReaderMock{}
}

func (m *s3ReaderMock) RetrieveArchivesFromS3() ([]ZipFile, error) {
	mock := new(s3ReaderMock)
	return nil, nil
}

func (m *s3ReaderMock) DownloadArchiveFromS3(fileName string) ([]byte, error) {
	return nil, nil
}

func TestHomepageHandlerSuccess(t *testing.T) {

}
