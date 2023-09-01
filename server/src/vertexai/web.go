package vertexai

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"project/config"
	"project/pb"
	"time"
)

// ChatHandle ...
func ChatHandle(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	data, debug, err := chatHandle(w, r)
	o := &pb.VaChatWebRsp{
		Data:  data,
		Debug: debug,
	}
	if err == nil {
		o.Ok = true
	} else {
		o.Error = err.Error()
	}
	if debug != nil {
		i := uint32(time.Since(t) / time.Millisecond)
		if i < 1 {
			i = 1
		}
		debug.TotalMs = i
	}
	ab, _ := json.Marshal(o)
	w.Write(ab)
}

func chatHandleInput(w http.ResponseWriter, r *http.Request) (*pb.VaChatReq, error) {

	if r.Method != `POST` {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(`method not allow`)
	}

	token := r.Header.Get(`VA-TOKEN`)

	if token == `` {
		w.WriteHeader(http.StatusNonAuthoritativeInfo)
		return nil, errors.New(`no token`)
	}
	if config.VAToken == `` {
		w.WriteHeader(http.StatusInternalServerError)
		err := errors.New(`no token in server`)
		return nil, err
	}
	if config.VAToken != token {
		w.WriteHeader(http.StatusUnauthorized)
		return nil, errors.New(`token not match`)
	}

	o := &pb.VaChatReq{}

	ab, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(`read body fail`)
	}
	err = json.Unmarshal(ab, o)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, err
	}

	return o, nil
}

func chatHandle(w http.ResponseWriter, r *http.Request) (*pb.VaChatRsp, *pb.VaDebug, error) {

	req, err := chatHandleInput(w, r)
	if err != nil {
		return nil, nil, err
	}

	rsp, err := Chat(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, err
	}

	var debug *pb.VaDebug
	if req.Debug {
		debug = rsp.Debug()
	}

	return rsp.Answer, debug, nil
}
