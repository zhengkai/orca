package core

import (
	"errors"
	"net/http"
	"project/metrics"
)

var errSkip = errors.New(`skip`)

// WebHandle ...
func (c *Core) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	metrics.ReqConcurrentInc()
	defer metrics.ReqConcurrentDec()

	p, err := req(w, r)
	if err != nil {
		if err != errSkip {
			metrics.ReqFailCount()
		}
		return
	}

	metrics.ReqBytes(len(p.Body))

	ab, cached, pr, err := c.getAB(p, r)
	if err != nil {
		if pr.httpCode != 0 {
			pr.httpCode = http.StatusInternalServerError
		}
		w.WriteHeader(pr.httpCode)
	}
	// zj.J(`cached`, cached)

	if len(ab) > 0 {
		w.Header().Add(`Content-Type`, `application/json`)
		w.Write(ab)
	}

	go doMetrics(ab, cached, r, p)
}
