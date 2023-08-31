package vertexai

import "google.golang.org/protobuf/types/known/structpb"

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
