package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project/es"
	"project/metrics"
	"project/pb"
	"project/util"
	"project/zj"
	"strings"

	"github.com/zhengkai/zu"
)

func doMetrics(ab []byte, cached bool, r *http.Request, req *pb.Req) {

	rspBytes := len(ab)
	metrics.RspBytes(rspBytes)

	o := &pb.Rsp{}
	err := json.Unmarshal(ab, o)
	if err != nil {
		zj.W(`unmarshal fail`, err)
		util.WriteFile(`metrics-json-fail`, ab)
		return
	}

	u := o.GetUsage()
	if u == nil {
		metrics.RspJSONFail()
		return
	}

	ip, err := util.GetIP(r)
	sip := ip.String()
	if err != nil {
		sip = `unknown`
	}

	key := strings.TrimPrefix(r.Header.Get(`Authorization`), `Bearer `)

	metrics.RspToken(u.PromptTokens, u.TotalTokens, cached)
	if !cached {
		metrics.RspTokenByModel(o.Model, u.TotalTokens)
		metrics.RspTokenByKey(key, u.TotalTokens)
		metrics.RspTokenByIP(sip, u.TotalTokens)
	}

	d := &pb.EsMetrics{
		ID: o.Id,
		Token: &pb.EsMetricsToken{
			Total:      u.TotalTokens,
			Completion: u.CompletionTokens,
			Prompt:     u.PromptTokens,
		},
		Cached:   cached,
		Ip:       sip,
		Model:    o.Model,
		Key:      key,
		ReqBytes: uint32(len(req.Body)),
		RspBytes: uint32(rspBytes),
		Ts:       zu.MS(),
		Hash:     fmt.Sprintf(`%x`, req.Hash()),
	}
	go es.Insert(d)
}
