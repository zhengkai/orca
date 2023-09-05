package metrics

import (
	"strconv"
)

var (
	reqCount      = newCounter(`req_count`, `HTTP 请求次数`)
	reqFailCount  = newCounter(`req_fail_count`, `无法响应的 HTTP 请求数`)
	reqBytes      = newCounter(`req_bytes`, `文件总上传量`)
	reqConcurrent = newGauge(`req_concurrent`, `当前并发请求数`)
	errorCount    = newCounterVec(`error_code`, `API 返回报错`, `code`)
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
