package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := db.Tasks(50)
	if err != nil {
		errorResponse(w, "error receiving tasks "+err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = make([]*db.Task, 0)
	}

	writeJSON(w, TasksResp{
		Tasks: tasks,
	})
}
