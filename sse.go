package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func monitor(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	flusher, _ := w.(http.Flusher)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		fmt.Fprintf(w, "data: %v\n\n", time.Now())

		flusher.Flush()

		time.Sleep(5 * time.Second)
	}
}

func main() {
	fmt.Println("sse")

	router := httprouter.New()
	router.GET("/monitor", monitor)

	router.ServeFiles("/public/*filepath", http.Dir("./public"))

	log.Fatal(http.ListenAndServe(":8080", router))
}
