package st

import (
	"fmt"
	"project/es"
	"project/pb"
	"project/zj"

	_ "embed" //

	"github.com/zhengkai/zu"
)

//go:embed tpl/date-histogram.json
var queryDateHistogramTpl string

// DateHistogram ...
func DateHistogram() {

	query := fmt.Sprintf(
		queryDateHistogramTpl,
		`now-24h`,
		`now`,
	)

	d := &pb.EsResultDateHistogram{}
	err := es.Search(query, d)
	if err != nil {
		zj.W(`fail`, err)
		return
	}

	li := d.GetAggregations().GetBytesSum().GetBuckets()
	if len(li) == 0 {
		return
	}

	re := make([]*pb.EsDateHistogram, len(li))
	for i, v := range li {
		re[i] = &pb.EsDateHistogram{
			Ts:              uint32(v.GetKey() / 1000),
			ReqBytes:        v.GetReqBytes().GetValue(),
			RspBytes:        v.GetRspBytes().GetValue(),
			TokenTotal:      v.GetTokenTotal().GetValue(),
			TokenCompletion: v.GetTokenCompletion().GetValue(),
			TokenPrompt:     v.GetTokenPrompt().GetValue(),
		}
	}

	j := zu.JSONPretty(re)
	// util.WriteFile(`date-histogram.json`, []byte(j))
	zj.J(j)
}
