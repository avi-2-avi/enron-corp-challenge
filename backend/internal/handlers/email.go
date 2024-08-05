package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetEmails(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetEmails now!!!"))
}

func GetEmail(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := database.GetIndexByID("emails", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(result, &response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	source := response["_source"].(map[string]interface{})
	sourceBytes, _ := json.Marshal(source)

	var email models.Email
	if err := json.Unmarshal(sourceBytes, &email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email.Id = response["_id"].(string)
	email.Timestamp = response["@timestamp"].(string)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(email)
}
