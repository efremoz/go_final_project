package main

import (
	"fmt"
	"log"
	"net/http"

	"go_final_project/pkg/db"
)

func main() {
	const port = 7540
	const webDir = "./web"

	dbFile := "scheduler.db"

	err := db.Init(dbFile)
	if err != nil {
		db.Close()
		log.Fatal("Ошибка инициализации БД: ", err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost:%d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
