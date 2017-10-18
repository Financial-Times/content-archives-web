package main

import (
	"net/http"

	fthealth "github.com/Financial-Times/go-fthealth/v1_1"
)

type HealthCheck struct {
}

func (h *HealthCheck) Health() func(w http.ResponseWriter, r *http.Request) {
	checks := []fthealth.Check{s3ConnectivityCheck()}

	healthCheck := &fthealth.HealthCheck{
		SystemCode:  "upp-exports",
		Name:        "UPP Daily Exports",
		Description: "Downloadable Content and Concept archives",
		Checks:      checks,
	}
	return fthealth.Handler(healthCheck)
}

func s3ConnectivityCheck() fthealth.Check {
	return fthealth.Check{
		ID:               "check-connectivity-to-s3",
		Name:             "Check connectivity to AWS S3",
		Severity:         1,
		BusinessImpact:   "Content and Concept archives won't be available for download",
		TechnicalSummary: "The service is unable to connect to AWS S3",
		PanicGuide:       "https://dewey.ft.com/upp-exports.html",
		Checker:          awsS3ConnectionChecker,
	}
}

func awsS3ConnectionChecker() (string, error) {
	err := func() error { return nil }() // DO A PROPER CHECK HERE
	if err == nil {
		return "Connectivity to AWS S3 is ok", err
	}
	return "Error connecting to AWS S3", err
}
