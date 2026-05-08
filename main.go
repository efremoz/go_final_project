package main

import (
	"fmt"
	"go_final_project/pkg/api"
	"log"
	"net/http"

	"go_final_project/pkg/db"
)

const port = 7540
const webDir = "./web"

func main() {

	dbFile := "scheduler.db"

	err := db.Init(dbFile)
	if err != nil {
		log.Fatal("db initialization error: ", err)
	}

	defer db.Close()

	api.Init()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("the server is running on http://localhost:%d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("server startup error:", err)
	}
}
