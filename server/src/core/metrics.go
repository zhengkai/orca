package core

import (
	"encoding/json"
	"net/http"
	"project/metrics"
	"project/pb"
	"project/util"
)

func doMetrics(ab []byte, cached bool, r *http.Request) {

	metrics.RspBytes(len(ab))

	o := &pb.Rsp{}
	json.Unmarshal(ab, o)

	u := o.GetUsage()
	if u == nil {
		metrics.RspJSONFail()
		return
	}

	metrics.RspToken(u.PromptTokens, u.TotalTokens, cached)

	ip, err := util.GetIP(r)
	sip := ip.String()
	if err != nil {
		sip = `unknown`
	}
	metrics.RspTokenByIP(sip, u.TotalTokens)
}
