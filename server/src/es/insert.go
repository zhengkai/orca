package es

import (
	"bytes"
	"project/pb"

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

	// zj.J(string(ab))

	// theClient.Create(`orca-metrics`, d.ID, bytes.NewReader(ab))

	index := indexName(uint32(d.Ts / 1000))

	re, err := theClient.Index(index, bytes.NewReader(ab))
	if err != nil {
		return
	}
	defer re.Body.Close()
	// zj.J(re.String())
}
