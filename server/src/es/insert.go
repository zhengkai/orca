package es

import (
	"bytes"
	"project/pb"
	"project/zj"

	"google.golang.org/protobuf/encoding/protojson"
)

// Insert ...
func Insert(d *pb.EsMetrics) {

	if theClient == nil {
		return
	}

	ab, err := protojson.MarshalOptions{
		EmitUnpopulated: true,
		AllowPartial:    true,
	}.Marshal(d)
	if err != nil {
		return
	}

	index := indexName(uint32(d.Ts / 1000))

	_, err = theClient.Index(index, bytes.NewReader(ab))
	if err != nil {
		zj.W(`insert fail:`, err)
		return
	}
}
