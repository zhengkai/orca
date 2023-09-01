package vertexai

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"project/pb"
	"project/util"
	"time"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"github.com/zhengkai/coral/v2"
	"github.com/zhengkai/life-go"
	"google.golang.org/protobuf/types/known/structpb"
)

var chatClient *aiplatform.PredictionClient

var errEmptyAnswer = errors.New(`empty answer`)
var errBlocked = errors.New(`blocked by google`)

var chatCache = coral.NewLRU(loadChatForCoral, 1000, 100)

type chatKey struct {
	System          string  `json:"system"`
	User            string  `json:"user"`
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens uint32  `json:"maxOutputTokens"`
	TopP            float32 `json:"topP"`
	TopK            uint32  `json:"topK"`
}

// ChatRsp ...
type ChatRsp struct {
	Answer *pb.VaChatRsp                 `json:"answer,omitempty"`
	Raw    *aiplatformpb.PredictResponse `json:"raw"`
	CostMs uint32                        `json:"costMs"`
}

// Chat ...
func Chat(req *pb.VaChatReq) (*ChatRsp, error) {

	p := req.Param
	if p == nil {
		p = defaultParam
	}

	k := chatKey{
		System:          req.System,
		User:            req.User,
		Temperature:     p.Temperature,
		MaxOutputTokens: p.MaxOutputTokens,
		TopP:            p.TopP,
		TopK:            p.TopK,
	}

	if req.NoCache {
		chatCache.Delete(k)
	} else {
		ab, err := util.ReadFile(chatCacheFile(k) + `.json`)
		if err == nil && len(ab) > 2 {
			rsp := &ChatRsp{}
			err = json.Unmarshal(ab, rsp)
			if err == nil {
				return rsp, nil
			}
		}
	}

	return chatCache.Get(k)
}

func buildChatReq(k chatKey) (*aiplatformpb.PredictRequest, error) {

	m := map[string]any{
		`context`: k.System,
	}
	if k.User != `` {
		m[`messages`] = []any{
			map[string]any{
				`author`:  `user`,
				`content`: k.User,
			},
		}
	}

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

	req := &aiplatformpb.PredictRequest{
		Endpoint: `projects/aigc-llm/locations/us-central1/publishers/google/models/chat-bison@001`,
		Instances: []*structpb.Value{
			structpb.NewStructValue(inst),
		},
		Parameters: structpb.NewStructValue(p),
	}

	return req, nil
}

func getChatVal(rsp *aiplatformpb.PredictResponse) (string, error) {

	p := rsp.GetPredictions()
	if len(p) == 0 {
		return ``, errEmptyAnswer
	}

	p0 := p[0]

	if isBlocked(p0) {
		return ``, errBlocked
	}

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

func isBlocked(o *structpb.Value) bool {
	sa := SpbMap(o, `safetyAttributes`).GetListValue().GetValues()
	for _, v := range sa {
		if SpbMap(v, `blocked`).GetBoolValue() {
			return true
		}
	}
	return false
}

func loadChat(k chatKey) (*ChatRsp, error) {

	req, err := buildChatReq(k)
	if err != nil {
		return nil, err
	}

	ctx, cancel := life.CTXTimeout(10 * time.Second)
	t := time.Now()
	rsp, err := chatClient.Predict(ctx, req)

	cancel()
	if err != nil {
		return nil, err
	}

	r := &ChatRsp{
		Raw:    rsp,
		CostMs: uint32(time.Since(t) / time.Millisecond),
	}

	answer := &pb.VaChatRsp{}
	answer.Content, err = getChatVal(rsp)
	if err == errBlocked {
		err = nil
		answer.Blocked = true
	}
	if err != nil {
		return nil, err
	}

	r.Answer = answer

	go chatSaveCache(k, r)
	return r, nil
}

func loadChatForCoral(k chatKey) (*ChatRsp, *time.Time, error) {
	r, err := loadChat(k)
	if err != nil {
		return nil, nil, err
	}
	return r, nil, nil
}

func chatCacheFile(k chatKey) string {
	ab, _ := json.Marshal(k)
	h := md5.Sum(ab)
	file := fmt.Sprintf(`vertexai/chat/%02x/%02x/%02x/%x`, h[0], h[1], h[2], h[3:])
	return file
}

func chatSaveCache(k chatKey, rsp *ChatRsp) {
	file := chatCacheFile(k)
	util.Mkdir(file)
	util.WriteJSON(file+`.json`, rsp)
}

// Debug ...
func (rsp *ChatRsp) Debug() *pb.VaDebug {
	d := &pb.VaDebug{
		CostMs: rsp.CostMs,
	}

	getToken(d, rsp.Raw)
	getSafety(d, rsp.Raw)

	return d
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
