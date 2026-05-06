package api

import "net/http"

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	default:
		errorResponse(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
