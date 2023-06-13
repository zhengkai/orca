package core

import (
	"errors"
	"net/http"
	"project/metrics"
	"project/zj"
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
		err500(w)
		return
	}

	metrics.ReqBytes(len(p.Body))

	ab, cached, err := c.getAB(p, r)
	if err != nil {
		err500(w)
		return
	}
	zj.J(`cached`, cached)

	w.Header().Add(`Content-Type`, `application/json`)
	w.Write(ab)

	go doMetrics(ab, cached, r, p)
}
