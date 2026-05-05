package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = 7540
	const webDir = "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на http://localhost:%d", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
