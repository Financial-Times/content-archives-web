package main

import "testing"

var s3Reader S3Reader = NewS3Reader("REGION", "BUCKET", "BUCKET_PREFIX")

func TestRetrieveArchivesFromS3Success(t *testing.T) {
	// zipFiles := s3Reader.RetrieveArchivesFromS3()
}

func TestRetrieveArchivesFromS3WithError(t *testing.T) {

}

func TestDownloadArchiveFromS3Success(t *testing.T) {

}

func TestDownloadArchiveFromS3WithError(t *testing.T) {

}
