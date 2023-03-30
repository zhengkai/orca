package core

import (
	"net/http"
	"project/metrics"
	"project/util"
)

// WebHandle ...
func (c *Core) WebHandle(w http.ResponseWriter, r *http.Request) {

	p, err := req(w, r)
	if err != nil {
		metrics.ReqFailCount()
		return
	}
	metrics.ReqBytes(len(p.Body))

	pr := c.add(p, r)

	go func() {
		reqFile := util.CacheName(p.Hash()) + `-req.json`
		if !util.FileExists(reqFile) {
			util.WriteFile(reqFile, p.Body)
		}
	}()

	pr.wait()

	if pr.err != nil {
		err500(w)
		return
	}

	w.Write(pr.rsp)
}
