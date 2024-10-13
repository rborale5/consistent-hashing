package monitoring

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"host"},
)

func totalRequestInc(host string) {
	totalRequests.WithLabelValues(host).Inc()
}

func init() {
	router := mux.NewRouter()
	router.Path("/prometheus").Handler(promhttp.Handler())
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
