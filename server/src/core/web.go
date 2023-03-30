package core

import (
	"net/http"
	"project/metrics"
)

// WebHandle ...
func (c *Core) WebHandle(w http.ResponseWriter, r *http.Request) {

	p, hash, err := req(w, r)
	if err != nil {
		metrics.ReqFailCount()
		return
	}
	metrics.ReqBytes(len(p.Body))

	pr := c.add(hash, p)

	pr.wait()

	if pr.err != nil {
		err500(w)
		return
	}

	w.Write(pr.rsp)
}
