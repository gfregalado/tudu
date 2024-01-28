package handlers

import (
	"database/sql"
	"errors"
	"github.com/gfregalado/todo/db"
	"github.com/gfregalado/todo/models"
	"log"
	"net/http"
	"strconv"
)

func Create(w http.ResponseWriter, r *http.Request) {
	dbClient, err := db.Init("./config/.env")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	} // Limit size of upload, e.g., 32MB
	name := r.FormValue("name")
	price, err := strconv.ParseFloat(r.FormValue("price"), 64)
	available, err := strconv.ParseBool(r.FormValue("available"))

	product := models.Product{
		Name:      name,
		Price:     price,
		Available: available,
	}

	p, err := product.Create(dbClient)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully added: %+v", p)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product created successfully"))
}

func GetById(w http.ResponseWriter, r *http.Request) {
	dbClient, err := db.Init("./config/.env")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var product models.Product

	p, err := product.GetById(dbClient, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "No product found with the given ID", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("Successfully fetched: %+v", p)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Fetched product successfully"))
}
