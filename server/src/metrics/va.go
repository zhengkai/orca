package metrics

var (
	vaTime = newSummary(`va_time`, `vertexai api time`)
)

// VaTime ...
func VaTime(ms uint32) {
	vaTime.Observe(float64(ms))
}
