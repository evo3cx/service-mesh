package services

import (
	"log"
	"net/http"
	"time"

	"github.com/openzipkin/zipkin-go/idgenerator"
	"github.com/prometheus/client_golang/prometheus"
	uuid "github.com/satori/go.uuid"
)

var (
	// Histogram is histogram metric.
	Histogram *prometheus.HistogramVec
	// Counter is counter metric.
	Counter *prometheus.CounterVec
)

func init() {
	Histogram = newHistogramVec()
	Counter = newCounterVec()

	prometheus.MustRegister(Histogram, Counter)
}

func newHistogramVec() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "latency",
		Help: "latency in seconds",
	}, []string{"service", "status"})
}

func newCounterVec() *prometheus.CounterVec {
	return prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "counter",
		Help: "counter of entity of service",
	}, []string{"service", "status"})
}

func CreateRequest(serviceName, address string, headers http.Header) {
	start := time.Now()
	status := "ok"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, address, nil)
	if err != nil {
		log.Printf("[Error] New request %s \n", err)
		return
	}

	req.Header = headerTracer(headers, serviceName)
	resp, err := client.Do(req)
	if err != nil {
		status = "fail"

	} else {
		if resp.StatusCode != 200 {
			status = "fail"
		}
	}

	Histogram.WithLabelValues(serviceName, status).Observe(time.Since(start).Seconds())
	Counter.WithLabelValues(serviceName, status).Inc()
}

func headerTracer(headers http.Header, serviceName string) http.Header {
	clientTraceID, _ := uuid.NewV4()
	gen := idgenerator.NewRandom64()
	tracID := gen.TraceID()
	spanID := gen.SpanID(tracID)

	if headers.Get("X-Client-Trace-Id") == "" {
		headers.Add("X-Client-Trace-Id", clientTraceID.String())
	}
	if headers.Get("X-B3-Spanid") != "" {
		headers.Add("X-B3-Parentspanid", headers.Get("X-B3-Spanid"))
	}
	headers.Add("X-B3-Spanid", spanID.String())

	if headers.Get("X-B3-TraceId") == "" {
		headers.Add("X-B3-TraceId", tracID.String())
	}

	headers.Add("Service-Name", serviceName)

	return headers
}
