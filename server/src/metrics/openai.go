package metrics

var (
	openaiTime = newSummary(`orca_openai_time`, `openai api time`)
)

// OpenAITime ...
func OpenAITime(ms uint32) {
	openaiTime.Observe(float64(ms))
}
