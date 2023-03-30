package core

import (
	"io"
	"net/http"
	"project/pb"
	"project/util"
	"project/zj"
)

func req(w http.ResponseWriter, r *http.Request) (p *pb.Req, err error) {

	path := r.URL.Path
	method := r.Method

	zj.J(method, path)

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
