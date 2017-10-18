package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const (
	indexTemplate = "templates/index.tmpl.html"
)

// Handler struct containig handlers configuration
type Handler struct {
	s3Service S3Service
	tmpl      *template.Template
}

// NewHandler returns a new Handler instance
func NewHandler(s3Service S3Service) Handler {
	tmpl := template.Must(template.ParseFiles(indexTemplate))
	return Handler{s3Service, tmpl}
}

// HomepageHandler serves the homepage content
func (h *Handler) HomepageHandler(w http.ResponseWriter, r *http.Request) {
	zipFiles, err := h.s3Service.RetrieveArchivesFromS3()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to get archives list from S3"))
	}
	h.tmpl.Execute(w, zipFiles)
}

// DownloadHandler starts the download of the specified file in request
func (h *Handler) DownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" || len(strings.TrimSpace(name)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please specify the name of the file"))
	}

	bytes, err := h.s3Service.DownloadArchiveFromS3(name)
	if err != nil {
		log.Println("Unable to download archive from S3", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Unable to download archive from S3"))
	}

	w.Header().Add("Content-Disposition", "attachment; filename="+name)
	_, err = w.Write(bytes)
	if err != nil {
		log.Println("Could not start file download", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not start file download"))
	}
}

// S3AutHandler middleware handler that adds authentication for the initial handler
func (h *Handler) S3AutHandler(f http.HandlerFunc) http.Handler {
	var handler http.Handler = http.HandlerFunc(f)
	return handler
}
