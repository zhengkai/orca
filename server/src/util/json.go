package util

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSONMarshal ...
func JSONMarshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		AllowPartial: true,
	}.Marshal(m)
}

// JSONUnmarshal ...
func JSONUnmarshal(ab []byte, m proto.Message) error {
	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}.Unmarshal(ab, m)
}
