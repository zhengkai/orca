package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var baseName = `orca_`

func newCounter(name, help string) prometheus.Counter {
	return prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: baseName + name,
			Help: help,
		},
	)
}
func newSummary(name, help string) prometheus.Summary {
	return prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       baseName + name,
		Help:       help,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
}

func newGauge(name, help string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name: baseName + name,
		Help: help,
	})
}

func newCounterVec(name, help, field string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: baseName + name,
			Help: help,
		},
		[]string{field},
	)
}

func newGaugeVec(name, help, field string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: baseName + name,
			Help: help,
		},
		[]string{field},
	)
}
