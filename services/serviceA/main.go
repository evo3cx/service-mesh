package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/evo3cx/service-mesh/services"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus"
	uuid "github.com/satori/go.uuid"
)

func main() {

	router := httprouter.New()
	router.GET("/healthz", handler)
	router.Handler("GET", "/metrics", prometheus.Handler())

	// spawn goroutine create a http request to dependency service
	go clientRequest()

	log.Println("Start Service-A at port :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func clientRequest() {
	for {
		reqId, _ := uuid.NewV4()
		headers := http.Header{}
		headers.Add("X-Request-Id", reqId.String())
		services.CreateRequest("service-B", "http://service_a_envoy:8881/", headers)
		time.Sleep(9 * time.Second)
	}
}

func handler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "OK")
}
