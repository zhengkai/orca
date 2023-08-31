package vertexai

import (
	"errors"
	"project/pb"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	"cloud.google.com/go/aiplatform/apiv1/aiplatformpb"
	"google.golang.org/protobuf/types/known/structpb"
)

var chatClient *aiplatform.PredictionClient

var errEmptyAnswer = errors.New(`empty answer`)
var errBlocked = errors.New(`blocked by google`)

// Chat ...
func Chat(req *pb.VaChatReq) {
}

func buildChatReq(system string, user []string, param *pb.VaParam) (*aiplatformpb.PredictRequest, error) {

	m := map[string]any{
		`context`: system,
	}
	if len(user) > 0 {
		var li []any
		for _, v := range user {
			li = append(li, map[string]any{
				`author`:  `user`,
				`content`: v,
			})
		}
		m[`messages`] = li
	}

	inst, err := structpb.NewStruct(m)
	if err != nil {
		return nil, err
	}

	p, err := structpb.NewStruct(map[string]any{
		`temperature`:     param.Temperature,
		`maxOutputTokens`: param.MaxOutputTokens,
		`topP`:            param.TopP,
		`topK`:            param.TopK,
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
