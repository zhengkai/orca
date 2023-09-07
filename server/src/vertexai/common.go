package vertexai

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"project/metrics"
	"project/pb"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/zhengkai/life-go"
	"google.golang.org/protobuf/types/known/structpb"
)

var theClient *aiplatform.PredictionClient
var errEmptyAnswer = errors.New(`empty answer`)
var errBlocked = errors.New(`blocked by google`)

type canDebug interface {
	GetDebug() bool
}

// SpbMap ...
func SpbMap(o *structpb.Value, key string) *structpb.Value {
	m := o.GetStructValue().GetFields()
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	return v
}

func cacheFile[T chatKey | textKey](t string, k T) string {
	ab, _ := json.Marshal(k)
	h := md5.Sum(ab)
	file := fmt.Sprintf(`vertexai/%s/%02x/%02x/%02x/%x`, t, h[0], h[1], h[2], h[3:])
	return file
}

func buildReq(k paramKey, m map[string]any, model string) (*aiplatformpb.PredictRequest, error) {

	inst, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	p, err := structpb.NewStruct(map[string]any{
		`temperature`:     k.Temperature,
		`maxOutputTokens`: k.MaxOutputTokens,
		`topP`:            k.TopP,
		`topK`:            k.TopK,
	})
	if err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf(
		`projects/aigc-llm/locations/us-central1/publishers/google/models/%s`,
		model,
	)

	req := &aiplatformpb.PredictRequest{
		Endpoint: endpoint,
		Instances: []*structpb.Value{
			structpb.NewStructValue(inst),
		},
		Parameters: structpb.NewStructValue(p),
	}

	return req, nil
}

func getToken(d *pb.VaDebug, raw *aiplatformpb.PredictResponse) {
	tm := SpbMap(raw.GetMetadata(), `tokenMetadata`)
	if tm == nil {
		return
	}
	input := SpbMap(tm, `inputTokenCount`)
	if input != nil {
		d.InputChar = uint32(SpbMap(input, `totalBillableCharacters`).GetNumberValue())
		d.InputToken = uint32(SpbMap(input, `totalTokens`).GetNumberValue())
	}
	output := SpbMap(tm, `outputTokenCount`)
	if output != nil {
		d.OutputChar = uint32(SpbMap(output, `totalBillableCharacters`).GetNumberValue())
		d.OutputToken = uint32(SpbMap(output, `totalTokens`).GetNumberValue())
	}
}
func getSafety(d *pb.VaDebug, raw *aiplatformpb.PredictResponse) {
	p := raw.GetPredictions()
	if len(p) == 0 {
		return
	}
	sa := SpbMap(p[0], `safetyAttributes`).GetListValue().GetValues()

	for _, v := range sa {
		c := SpbMap(v, `categories`).GetListValue().GetValues()
		if len(c) == 0 {
			continue
		}
		s := SpbMap(v, `scores`).GetListValue().GetValues()
		if len(s) != len(c) {
			continue
		}

		for i, cv := range c {
			row := &pb.VaSafety{
				Category: cv.GetStringValue(),
				Score:    float32(s[i].GetNumberValue()),
			}
			d.Safety = append(d.Safety, row)
		}
	}
}

func isBlocked(o *structpb.Value) bool {
	sa := SpbMap(o, `safetyAttributes`).GetListValue().GetValues()
	for _, v := range sa {
		if SpbMap(v, `blocked`).GetBoolValue() {
			return true
		}
	}
	return false
}

func getVal(rsp *aiplatformpb.PredictResponse) (string, error) {

	p := rsp.GetPredictions()
	if len(p) == 0 {
		return ``, errEmptyAnswer
	}

	p0 := p[0]

	if isBlocked(p0) {
		return ``, errBlocked
	}

	// for text
	c := SpbMap(p0, `content`).GetStringValue()
	if c != `` {
		return c, nil
	}

	// for chat
	ca := SpbMap(p0, `candidates`).GetListValue().GetValues()
	if len(ca) == 0 {
		return ``, errEmptyAnswer
	}

	s := SpbMap(ca[0], `content`).GetStringValue()
	if s == `` {
		return ``, errEmptyAnswer
	}

	return s, nil
}

func doRequest(req *aiplatformpb.PredictRequest) (*Rsp, error) {

	ctx, cancel := life.CTXTimeout(10 * time.Second)
	t := time.Now()
	rsp, err := theClient.Predict(ctx, req)

	cancel()
	if err != nil {
		return nil, err
	}

	r := &Rsp{
		Raw:    rsp,
		CostMs: uint32(time.Since(t) / time.Millisecond),
	}
	metrics.VaTime(r.CostMs)

	answer := &pb.VaRsp{}
	answer.Content, err = getVal(rsp)
	if err == errBlocked {
		err = nil
		answer.Blocked = true
	}
	if err != nil {
		return nil, err
	}

	r.Answer = answer

	return r, nil
}
