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
	rspTokenByIP        = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: `orca_token_by_ip`,
			Help: `API 返回报错`,
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
