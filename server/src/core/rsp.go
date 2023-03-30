package core

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"project/config"
	"project/pb"
	"project/zj"
	"time"
)

func fetchRemote(r *pb.Req) (ab []byte, err error) {

	u, err := url.Parse(config.RemoteAPI)
	if err != nil {
		zj.W(`url fail`, config.RemoteAPI, err)
		return
	}
	u.Path = r.Path

	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader(r.Body))
	if err != nil {
		return
	}

	req.Header.Set(`Content-Type`, `application/json`)
	req.Header.Set(`Authorization`, `Bearer `+config.OpenAIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {
		return
	}

	defer rsp.Body.Close()

	return io.ReadAll(rsp.Body)
}
