package core

import (
	"io"
	"net/http"
	"project/pb"
	"project/util"
	"project/zj"
)

func (c *Core) getAB(p *pb.Req, r *http.Request) (ab []byte, cached bool, err error) {

	canCache := p.Method != http.MethodGet && p.Method != http.MethodDelete

	canCache = false

	if canCache {
		var ok bool
		ab, ok = tryCache(p)
		if ok {
			cached = true
			return
		}
	}

	pr, cached := c.add(p, r)

	if canCache {
		go func() {
			reqFile := util.CacheName(p.Hash()) + `-req.json`
			if !util.FileExists(reqFile) {
				util.WriteFile(reqFile, p.Body)
			}
		}()
	}

	pr.wait()

	ab = pr.rsp
	err = pr.err
	return
}

func req(w http.ResponseWriter, r *http.Request) (p *pb.Req, err error) {

	url := r.URL.String()
	method := r.Method
	contentType := r.Header.Get(`Content-Type`)
	if contentType == `` {
		contentType = `application/json`
	}

	if url == `/favicon.ico` {
		err = errSkip
		return
	}

	ab, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1024*1024*10))
	go util.WriteFile(`last-req.json`, ab)
	if err != nil {
		return
	}

	p = &pb.Req{
		Url:         url,
		Method:      method,
		ContentType: contentType,
		Body:        ab,
	}
	zj.F(`%x %s %s %s`, p.Hash(), method, url, contentType)
	return
}

func err500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
