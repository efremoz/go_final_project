package api

import (
	"encoding/json"
	"net/http"

	"go_final_project/pkg/db"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	default:
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errorResponse(w, "ID not specified", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, task)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errorResponse(w, "error deserializing JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		errorResponse(w, "the issue ID is not specified", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		errorResponse(w, "the issue title is not specified", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := db.UpdateTask(&task)
	if err != nil {
		errorResponse(w, "error updating the task in the database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, struct{}{})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		errorResponse(w, "ID not specified", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, struct{}{})
}
