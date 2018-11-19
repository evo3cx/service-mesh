package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	fmt.Fprintf(w, "Hello from service C")
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Start Service-C at port :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
