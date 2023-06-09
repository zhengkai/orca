package core

import (
	"encoding/json"
	"net/http"
	"project/es"
	"project/metrics"
	"project/pb"
	"project/util"
	"project/zj"
	"strings"

	"github.com/zhengkai/zu"
)

func doMetrics(ab []byte, cached bool, r *http.Request, reqBytes int) {

	rspBytes := len(ab)
	metrics.RspBytes(rspBytes)

	o := &pb.Rsp{}
	err := json.Unmarshal(ab, o)
	if err != nil {
		zj.J(`unmarshal fail`, err)
		util.WriteFile(`metrics-json-fail`, ab)
		return
	}

	u := o.GetUsage()
	if u == nil {
		metrics.RspJSONFail()
		return
	}

	metrics.RspToken(u.PromptTokens, u.TotalTokens, cached)
	if !cached {
		zj.J(`token`, u.PromptTokens, u.TotalTokens)
	}

	metrics.RspTokenByModel(o.Model, u.TotalTokens)

	key := strings.TrimPrefix(r.Header.Get(`Authorization`), `Bearer `)
	metrics.RspTokenByKey(key, u.TotalTokens)

	ip, err := util.GetIP(r)
	sip := ip.String()
	if err != nil {
		sip = `unknown`
	}
	metrics.RspTokenByIP(sip, u.TotalTokens)

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
		ReqBytes: uint32(reqBytes),
		RspBytes: uint32(rspBytes),
		Ts:       zu.MS(),
	}
	go es.Insert(d)
}
