package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	rspBytes            = newCounter(`orca_rsp_bytes`, `rsp bytes`)
	rspPromptTokenCount = newCounter(`orca_rsp_prompt_token`, `prompt token`)
	rspTokenCount       = newCounter(`orca_rsp_token`, `token`)
	rspTokenCachedCount = newCounter(`orca_rsp_token_cached`, `token cached`)
	rspJSONFailCount    = newCounter(`orca_rsp_json_fail`, `json fail`)
	rspTokenByModel     = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: `orca_token_by_model`,
			Help: `token by model`,
		},
		[]string{`model`},
	)
	rspTokenByKey = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: `orca_token_by_key`,
			Help: `openai key`,
		},
		[]string{`key`},
	)
	rspTokenByIP = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: `orca_token_by_ip`,
			Help: `token by ip`,
		},
		[]string{`ip`},
	)
)

// RspToken ...
func RspToken(prompt, total uint32, cached bool) {
	if cached {
		rspTokenCachedCount.Add(float64(total))
		return
	}
	rspPromptTokenCount.Add(float64(prompt))
	rspTokenCount.Add(float64(total))
}

// RspBytes ...
func RspBytes(n int) {
	rspBytes.Add(float64(n))
}

// RspJSONFail ...
func RspJSONFail() {
	rspJSONFailCount.Inc()
}

// RspTokenByIP ...
func RspTokenByIP(ip string, token uint32) {
	rspTokenByIP.WithLabelValues(ip).Add(float64(token))
}

// RspTokenByKey ...
func RspTokenByKey(key string, token uint32) {
	if len(key) > 30 {
		key = key[:30]
	} else if key == `` {
		key = `<empty>`
	}
	rspTokenByKey.WithLabelValues(key).Add(float64(token))
}

// RspTokenByModel ...
func RspTokenByModel(model string, token uint32) {
	rspTokenByModel.WithLabelValues(model).Add(float64(token))
}
