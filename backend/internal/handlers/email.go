package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"backend/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetEmails(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil || size < 1 {
		size = 10
	}

	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "@timestamp"
	}
	order := r.URL.Query().Get("order")
	if order == "" {
		order = "desc"
	}

	query := utils.BuildQuery(filter, page, size, sort, order)

	queryBytes, err := json.Marshal(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := database.GetIndexDocuments("emails", queryBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(results, &response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hits, ok := response["hits"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	totalHits, ok := hits["total"].(map[string]interface{})
	if !ok {
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	totalElements, ok := totalHits["value"].(float64)
	if !ok {
		http.Error(w, "Invalid response format", http.StatusInternalServerError)
		return
	}

	emails := response["hits"].(map[string]interface{})["hits"].([]interface{})

	var emailList []models.Email
	for _, e := range emails {
		var email models.Email
		hit := e.(map[string]interface{})
		source := hit["_source"].(map[string]interface{})
		sourceBytes, _ := json.Marshal(source)
		json.Unmarshal(sourceBytes, &email)

		email.Id = hit["_id"].(string)
		email.Timestamp = hit["@timestamp"].(string)

		emailList = append(emailList, email)
	}

	pagination := utils.Paginate(emailList, int(totalElements), page, size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pagination)
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
