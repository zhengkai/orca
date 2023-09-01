package vertexai

import (
	"project/pb"

	"google.golang.org/protobuf/types/known/structpb"
)

var defaultParam = &pb.VaParam{
	Temperature:     0.2,
	MaxOutputTokens: 0,
	TopP:            1,
	TopK:            40,
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
