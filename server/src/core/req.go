package core

import (
	"io"
	"net/http"
	"project/pb"
	"project/util"
	"project/zj"
)

func (c *Core) getAB(p *pb.Req, r *http.Request) (ab []byte, cached bool, err error) {
	ab, ok := tryCache(p)
	if ok {
		cached = true
		return
	}

	pr, cached := c.add(p, r)

	go func() {
		reqFile := util.CacheName(p.Hash()) + `-req.json`
		if !util.FileExists(reqFile) {
			util.WriteFile(reqFile, p.Body)
		}
	}()

	pr.wait()

	ab = pr.rsp
	err = pr.err
	return
}

func req(w http.ResponseWriter, r *http.Request) (p *pb.Req, err error) {

	path := r.URL.Path
	method := r.Method

	if path == `/favicon.ico` {
		err = errSkip
		return
	}

	ab, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024*1024))
	if err != nil {
		return
	}

	go util.WriteFile(`last-req.json`, ab)

	p = &pb.Req{
		Path:   path,
		Method: method,
		Body:   ab,
	}
	zj.F(`%x %s %s`, p.Hash(), method, path)
	return
}

func err500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
