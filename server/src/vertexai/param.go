package vertexai

import "project/pb"

var defaultParam = &pb.VaParam{
	Temperature:     0.2,
	MaxOutputTokens: 0,
	TopP:            1,
	TopK:            40,
}

// paramKey  ...
type paramKey struct {
	Temperature     float32 `json:"temperature"`
	MaxOutputTokens uint32  `json:"maxOutputTokens"`
	TopP            float32 `json:"topP"`
	TopK            uint32  `json:"topK"`
}

func (k *paramKey) load(p *pb.VaParam) {
	if p == nil {
		p = defaultParam
	}
	k.Temperature = p.Temperature
	k.MaxOutputTokens = p.MaxOutputTokens
	k.TopP = p.TopP
	k.TopK = p.TopK

	if k.MaxOutputTokens == 0 {
		k.MaxOutputTokens = 512
	}
}
