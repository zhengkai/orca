package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(reqCount)
	prometheus.MustRegister(reqFailCount)
	prometheus.MustRegister(reqBytes)
	prometheus.MustRegister(errorCount)

	prometheus.MustRegister(rspBytes)
	prometheus.MustRegister(rspPromptTokenCount)
	prometheus.MustRegister(rspTokenCount)
	prometheus.MustRegister(rspTokenCachedCount)
	prometheus.MustRegister(rspJSONFailCount)
	prometheus.MustRegister(rspTokenByIP)
}
