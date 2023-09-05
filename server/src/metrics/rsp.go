package metrics

var (
	rspBytes            = newCounter(`rsp_bytes`, `rsp bytes`)
	rspPromptTokenCount = newCounter(`rsp_prompt_token`, `prompt token`)
	rspTokenCount       = newCounter(`rsp_token`, `token`)
	rspTokenCachedCount = newCounter(`rsp_token_cached`, `token cached`)
	rspJSONFailCount    = newCounter(`rsp_json_fail`, `json fail`)
	rspTokenByModel     = newCounterVec(
		`token_by_model`,
		`token by model`,
		`model`,
	)
	rspTokenByKey = newCounterVec(
		`token_by_key`,
		`openai key`,
		`key`,
	)
	rspTokenByIP = newCounterVec(
		`token_by_ip`,
		`token by ip`,
		`ip`,
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
