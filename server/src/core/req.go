package core

import (
	"crypto/md5"
	"io/ioutil"
	"net/http"
	"project/pb"
	"project/util"

	"google.golang.org/protobuf/proto"
)

func req(w http.ResponseWriter, r *http.Request) (p *pb.Req, hash [16]byte, err error) {

	path := r.URL.Path
	method := r.Method

	ab, err := ioutil.ReadAll(http.MaxBytesReader(w, r.Body, 1024*1024))
	if err != nil {
		return
	}

	go util.WriteFile(`last-req.json`, ab)

	p = &pb.Req{
		Path:   path,
		Method: method,
		Body:   ab,
	}

	pab, err := proto.Marshal(p)
	if err != nil {
		err500(w)
		return
	}

	hash = md5.Sum(pab)
	return
}

func err500(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}
