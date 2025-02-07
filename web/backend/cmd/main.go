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

type ParecerInput struct {
	User    string `json:"user"`
	Creci   string `json:"creci"`
	Content string `json:"content"`
}

var (
	dbClient storage.SQLClient
)

func init() {

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)

	dbClient = storage.NewSQLClient(connStr)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/parecer", handleParecer)

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func handleParecer(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodPost:
		CreateParecer(w, r)
	case http.MethodGet:
		ReadParecer(w, r)
	case http.MethodPut:
		UpdateParecer(w, r)
	case http.MethodDelete:
		DeleteParecer(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func CreateParecer(w http.ResponseWriter, r *http.Request) {
	var input ParecerInput

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

func ReadParecer(w http.ResponseWriter, r *http.Request) {
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

func UpdateParecer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input ParecerInput
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

	if err := dbClient.UpdateParecer(id, data); err != nil {
		log.Println("Error updating parecer in database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func DeleteParecer(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := dbClient.DeleteParecer(id); err != nil {
		log.Println("Error deleting parecer from database", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
