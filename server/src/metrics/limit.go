package metrics

import (
	"net/http"
	"strconv"
)

var (
	limitReq = newGaugeVec(
		`limit_req_by_model`,
		`limit request by model`,
		`model`,
	)
	limitToken = newGaugeVec(
		`limit_token_by_model`,
		`limit token by model`,
		`model`,
	)
)

// Limit ...
func Limit(h http.Header) {

	model := h.Get(`openai-model`)
	if model == `` {
		return
	}

	if h.Get(`x-ratelimit-limit-requests`) != `` {
		req := h.Get(`x-ratelimit-remaining-requests`)
		limitReq.WithLabelValues(model).Set(strToFloat(req))
	}
	if h.Get(`x-ratelimit-limit-tokens`) != `` {
		token := h.Get(`x-ratelimit-remaining-tokens`)
		limitToken.WithLabelValues(model).Set(strToFloat(token))
		// zj.J(`limit time`, model, h.Get(`x-ratelimit-reset-requests`), h.Get(`x-ratelimit-reset-tokens`))
	}
}

func strToFloat(s string) float64 {
	if s == `` {
		return 0
	}
	f, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return float64(f)
}
