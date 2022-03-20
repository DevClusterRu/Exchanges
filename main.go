package main

import (
	"Exchange/internal"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	metricsObject := internal.MetricsStructure{}
	metricsObject.MChannel = make(chan internal.Metric, 1000)
	metricsObject.AChannel = make(chan internal.Metric, 1000)
	metricsObject.Metrics = make(map[string]float64)
	metricsObject.Accumulator = make(map[string]float64)
	metricsObject.M = &sync.Mutex{}

	metricsObject.MetricsProcessor()

	go func() {
		http.HandleFunc("/metrics", metricsObject.ShowMetrics)
		log.Println("Starting webserver...")
		err := http.ListenAndServe(":9100", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		metricsObject.GetRequest("exchangesBuy", "https://www.bestchange.ru/visa-mastercard-rub-to-tether-trc20.html")
		time.Sleep(5 * time.Second)
		metricsObject.GetRequest("exchangeSell", "https://www.bestchange.ru/tether-trc20-to-tinkoff.html")
		time.Sleep(5 * time.Second)
	}

}
