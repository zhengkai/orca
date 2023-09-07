package vertexai

import (
	"project/pb"
	"project/util"

	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
)

// Rsp ...
type Rsp struct {
	Answer *pb.VaRsp                     `json:"answer,omitempty"`
	Raw    *aiplatformpb.PredictResponse `json:"raw"`
	CostMs uint32                        `json:"costMs"`
}

// Debug ...
func (rsp *Rsp) Debug() *pb.VaDebug {
	d := &pb.VaDebug{
		CostMs: rsp.CostMs,
	}

	getToken(d, rsp.Raw)
	getSafety(d, rsp.Raw)
	return d
}

func (rsp *Rsp) save(file string) {
	util.Mkdir(file)
	util.WriteJSON(file+`.json`, rsp)
}
