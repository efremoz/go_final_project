package main

import (
	"fmt"
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"log"
	"net/http"
	"os"
	"strconv"
)

const portDefault = 7540
const webDir = "./web"

func main() {

	dbFile := "data/scheduler.db"

	port := getPort()

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

func getPort() int {
	portStr := os.Getenv("TODO_PORT")
	if portStr == "" {
		return portDefault
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Printf("error converting TODO_PORT='%s' in particular, the port is being used %d", portStr, port)
		return portDefault
	}

	return port
}
