package metrics

var (
	openaiTime = newSummary(`openai_time`, `openai api time`)
)

// OpenAITime ...
func OpenAITime(ms uint32) {
	openaiTime.Observe(float64(ms))
}
