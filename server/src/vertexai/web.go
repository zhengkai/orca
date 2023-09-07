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

// web type
const (
	WebText Web = iota + 1
	WebChat
)

// Web ...
type Web int

func (web Web) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	t := time.Now()

	data, debug, err := handle(web, w, r)
	o := &pb.VaWebRsp{
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

func chatHandleInput(w http.ResponseWriter, r *http.Request) ([]byte, error) {

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

	ab, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return nil, errors.New(`read body fail`)
	}

	return ab, nil
}

func handle(web Web, w http.ResponseWriter, r *http.Request) (*pb.VaRsp, *pb.VaDebug, error) {

	ab, err := chatHandleInput(w, r)
	if err != nil {
		return nil, nil, err
	}

	var rsp *Rsp
	isDebug := false

	if web == WebChat {

		req := &pb.VaChatReq{}
		isDebug, err = unmarshalWebReq(w, ab, req)
		if err != nil {
			return nil, nil, err
		}

		rsp, err = Chat(req)

	} else {

		req := &pb.VaTextReq{}
		isDebug, err = unmarshalWebReq(w, ab, req)
		if err != nil {
			return nil, nil, err
		}

		rsp, err = Text(req)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return nil, nil, err
	}

	var debug *pb.VaDebug
	if isDebug {
		debug = rsp.Debug()
	}

	return rsp.Answer, debug, nil
}

func unmarshalWebReq(w http.ResponseWriter, ab []byte, req canDebug) (bool, error) {

	err := json.Unmarshal(ab, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return false, err
	}

	return req.GetDebug(), nil
}
