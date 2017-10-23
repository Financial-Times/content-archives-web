package main

import (
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

//HealthCheck service
type HealthCheck struct {
	bucketName   string
	bucketPrefix string
}

// Health check of the app
func (h *HealthCheck) Health() func(w http.ResponseWriter, r *http.Request) {
	checks := []fthealth.Check{h.s3ConnectivityCheck()}

	healthCheck := &fthealth.HealthCheck{
		SystemCode:  "upp-exports",
		Name:        "UPP Daily Exports",
		Description: "Downloadable Content and Concept archives",
		Checks:      checks,
	}
	return fthealth.Handler(healthCheck)
}

func (h *HealthCheck) s3ConnectivityCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "check-connectivity-to-s3",
		Name:             "Check connectivity to AWS S3",
		Severity:         1,
		BusinessImpact:   "Content and Concept archives won't be available for download",
		TechnicalSummary: "The service is unable to connect to AWS S3",
		PanicGuide:       "https://dewey.ft.com/upp-exports.html",
		Checker:          h.awsS3ConnectionChecker,
	}
}

func (h *HealthCheck) awsS3ConnectionChecker() (string, error) {
	input := &s3.ListObjectsInput{
		Bucket: aws.String(h.bucketName),
		Prefix: aws.String(h.bucketPrefix),
	}

	_, err := NewS3Reader().Read(input)
	if err == nil {
		return "Connectivity to AWS S3 is ok", err
	}
	return "AWS S3 connection couldn't be established. ", err
}
