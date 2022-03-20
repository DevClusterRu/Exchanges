package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Metric struct {
	Name  string  //4 example temperature{"device=1"}
	Value float64 //Float value
}

type MetricsStructure struct {
	MChannel    chan Metric
	AChannel    chan Metric
	Metrics     map[string]float64
	Accumulator map[string]float64
	M           *sync.Mutex
}

func (mt *MetricsStructure) ShowMetrics(w http.ResponseWriter, r *http.Request) {
	mt.M.Lock()
	for k, v := range mt.Metrics {
		fmt.Fprintf(w, fmt.Sprintf("%s %f\n", k, v))
	}
	mt.M.Unlock()
}

func (mt *MetricsStructure) MetricsProcessor() {

	go func() {
		for {
			mt.M.Lock()
			for k, v := range mt.Accumulator {
				mt.Metrics[k] = v
				mt.Accumulator[k] = 0
			}
			mt.M.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case m := <-mt.MChannel:
				mt.M.Lock()
				mt.Metrics[m.Name] = m.Value
				mt.M.Unlock()
			case a := <-mt.AChannel:
				mt.M.Lock()
				mt.Accumulator[a.Name] += a.Value
				mt.M.Unlock()
			}
		}
	}()

}
