package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/constants"
	"go_final_project/pkg/db"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errorResponse(w, "error deserializing JSON: "+err.Error(), http.StatusBadRequest)
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

	id, err := db.AddTask(&task)
	if err != nil {
		errorResponse(w, "error adding a task to the database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id": fmt.Sprintf("%d", id),
	}
	writeJSON(w, response)
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(constants.DateFormat)
		return nil
	}

	t, err := time.Parse(constants.DateFormat, task.Date)
	if err != nil {
		return err
	}

	if !afterNow(t, now) {
		return nil
	}
	if task.Repeat == "" {
		task.Date = now.Format(constants.DateFormat)
		return nil
	}

	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = next
	return nil
}
