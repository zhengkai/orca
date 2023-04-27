package core

import (
	"encoding/json"
	"net/http"
	"project/metrics"
	"project/pb"
	"project/util"
	"project/zj"
)

func doMetrics(ab []byte, cached bool, r *http.Request) {

	metrics.RspBytes(len(ab))

	o := &pb.Rsp{}
	err := json.Unmarshal(ab, o)
	if err != nil {
		zj.J(`unmarshal fail`, err)
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

	ip, err := util.GetIP(r)
	sip := ip.String()
	if err != nil {
		sip = `unknown`
	}
	metrics.RspTokenByIP(sip, u.TotalTokens)
}
