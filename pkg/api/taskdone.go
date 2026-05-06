package api

import (
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = db.UpdateDate(id, next)
		if err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, struct{}{})
}
