package metrics

var (
	vaTime = newSummary(`orca_va_time`, `vertexai api time`)
)

// VaTime ...
func VaTime(ms uint32) {
	vaTime.Observe(float64(ms))
}
