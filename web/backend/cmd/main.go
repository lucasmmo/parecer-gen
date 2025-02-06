package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"parecer-gen/pkg/file"
	"parecer-gen/pkg/parecer"
	"parecer-gen/pkg/storage"

	_ "github.com/lib/pq"
)

type ParecerPostInput struct {
	User    string `json:"user"`
	Creci   string `json:"creci"`
	Content string `json:"content"`
}

var (
	dbClient storage.SQLClient
)

func init() {
	var awsEndpoint string
	var dbHost string
	awsRegion := "us-east-1"

	if os.Getenv("IN_DOCKER") == "true" {
		awsEndpoint = "http://localstack:4566"
		dbHost = "postgres"
	} else {
		awsEndpoint = "http://localhost:4566"
		os.Setenv("AWS_ACCESS_KEY_ID", "YOUR_AKID")
		os.Setenv("AWS_SECRET_ACCESS_KEY", "YOUR_SECRET_KEY")
		os.Setenv("AWS_SESSION_TOKEN", "TOKEN")
	}

	connStr := fmt.Sprintf("host=%s port=5432 user=postgres password=postgres dbname=postgres sslmode=disable", dbHost)

	storage.NewS3Client(awsEndpoint, awsRegion)
	dbClient = storage.NewSQLClient(connStr)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /parecer", CreateParecerHandler)
	mux.HandleFunc("GET /parecer", GetParecerHandler)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func CreateParecerHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var input ParecerPostInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("Error decoding request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := parecer.NewData(input.User, input.Creci, input.Content)
	if err != nil {
		log.Println("Error creating parecer data", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := dbClient.SaveParecer(data); err != nil {
		log.Println("Error saving parecer to database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetParecerHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	id := r.URL.Query().Get("id")
	if id == "" {

		res, err := dbClient.GetAllParecer()
		if err != nil {
			log.Println("Error getting all pareceres from database", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	data, err := dbClient.GetParecer(id)
	if err != nil {
		log.Println("Error getting parecer from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch w.Header().Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	default:
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)

		HTMLfile, err := file.GenerateParecerHTML(data.User, data.Creci, data.Date, data.Content)
		if err != nil {
			log.Println("Error generating parecer HTML", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		PDFfile, err := file.GeneratePDF(HTMLfile)
		if err != nil {
			log.Println("Error generating parecer PDF", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := io.Copy(w, PDFfile.Reader); err != nil {
			log.Println("Error copying PDF file to response writer", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
