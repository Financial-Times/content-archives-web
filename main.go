package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/joho/godotenv"
)

var sessionStore *sessions.CookieStore

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load .env file, falling back to global env variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		log.Fatal("$AWS_REGION must be set")
	}

	awsBucketName := os.Getenv("AWS_BUCKET_NAME")
	if awsBucketName == "" {
		log.Fatal("$AWS_BUCKET_NAME must be set")
	}

	awsBucketPrefix := os.Getenv("AWS_BUCKET_PREFIX")
	if awsBucketPrefix == "" {
		log.Fatal("$AWS_BUCKET_PREFIX must be set")
	}

	oktaClientID := os.Getenv("OKTA_CLIENT_ID")
	if oktaClientID == "" {
		log.Fatal("$OKTA_CLIENT_ID must be set")
	}

	oktaClientSecret := os.Getenv("OKTA_CLIENT_SECRET")
	if oktaClientSecret == "" {
		log.Fatal("$OKTA_CLIENT_SECRET must be set")
	}

	oktaScope := os.Getenv("OKTA_SCOPE")
	if oktaScope == "" {
		oktaScope = "openid offline_access"
	}

	issuer := os.Getenv("ISSUER")
	if issuer == "" {
		log.Fatal("ISSUER must be set")
	}

	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		log.Fatal("$SESSION_KEY must be set")
	}

	callbackURL := os.Getenv("CALLBACK_URL")
	if callbackURL == "" {
		log.Fatal("$CALLBACK_URL must be set")
	}

	sessionStore = sessions.NewCookieStore([]byte(sessionKey))
	s3Service := NewS3Service(awsRegion, awsBucketName, awsBucketPrefix)
	healthCheck := HealthCheck{awsBucketName, awsBucketPrefix}
	appHandler := NewHandler(s3Service, HandlerConfig{
		oktaClientID,
		oktaClientSecret,
		oktaScope,
		issuer,
		sessionKey,
		callbackURL,
	})

	// Set default Max-Age for the session cookie for one hour
	sessionStore.MaxAge(3600)

	r := mux.NewRouter()
	// load static files
	staticH := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(staticH)

	// okta reserved routes
	r.HandleFunc("/login", appHandler.LoginHandler)
	r.HandleFunc("/authorization-code/callback", appHandler.AuthCodeCallbackHandler)

	// use middlewares to restrict access to FT members only
	r.Handle("/", appHandler.AuthHandler(appHandler.HomepageHandler))
	// r.HandleFunc("/", appHandler.HomepageHandler)
	r.Handle("/download/{prefix}/{name}", appHandler.AuthHandler(appHandler.DownloadHandler))

	// health should be accessible for anyone
	r.HandleFunc("/__health", healthCheck.Health())

	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
