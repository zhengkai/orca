package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(reqCount)
	prometheus.MustRegister(reqFailCount)
	prometheus.MustRegister(reqBytes)
	prometheus.MustRegister(reqConcurrent)

	prometheus.MustRegister(errorCount)

	prometheus.MustRegister(rspBytes)
	prometheus.MustRegister(rspPromptTokenCount)
	prometheus.MustRegister(rspTokenCount)
	prometheus.MustRegister(rspTokenCachedCount)
	prometheus.MustRegister(rspJSONFailCount)
	prometheus.MustRegister(rspTokenByIP)
	prometheus.MustRegister(rspTokenByModel)
	prometheus.MustRegister(rspTokenByKey)

	prometheus.MustRegister(limitReq)
	prometheus.MustRegister(limitToken)

	prometheus.MustRegister(openaiTime)
	prometheus.MustRegister(vaTime)
}
