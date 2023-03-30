package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(reqCount)
	prometheus.MustRegister(reqFailCount)
	prometheus.MustRegister(reqBytes)
	prometheus.MustRegister(errorCount)
}
