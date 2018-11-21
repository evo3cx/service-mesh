package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/evo3cx/service-mesh/services"
	"github.com/prometheus/client_golang/prometheus"
)

func handler(w http.ResponseWriter, r *http.Request) {
	services.CreateRequest("service-B", "http://service_b_envoy:8882/", r.Header)
	fmt.Fprintf(w, "Hello from service B")
}

func main() {
	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", handler)

	log.Println("Start Service-B at port :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
