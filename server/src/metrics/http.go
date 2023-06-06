package metrics

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	reqCount      = newCounter(`orca_req_count`, `HTTP 请求次数`)
	reqFailCount  = newCounter(`orca_req_fail_count`, `无法响应的 HTTP 请求数`)
	reqBytes      = newCounter(`orca_req_bytes`, `文件总上传量`)
	reqConcurrent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: `orca_req_concurrent`,
			Help: `当前并发请求数`,
		},
	)
	errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: `orca_error_code`,
			Help: `API 返回报错`,
		},
		[]string{`code`},
	)
)

// ReqFailCount ...
func ReqFailCount() {
	reqFailCount.Inc()
}

// ReqConcurrentInc ...
func ReqConcurrentInc() {
	reqConcurrent.Inc()
}

// ReqConcurrentDec ...
func ReqConcurrentDec() {
	reqConcurrent.Dec()
}

// ReqBytes ...
func ReqBytes(n int) {
	reqCount.Inc()
	reqBytes.Add(float64(n))
}

// ErrorCount ...
func ErrorCount(code int32) {
	c := strconv.Itoa(int(code))
	errorCount.WithLabelValues(c).Inc()
}
