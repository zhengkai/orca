package es

import (
	"bytes"
	"encoding/json"
	"project/pb"
)

// Insert ...
func Insert(d *pb.EsMetrics) {

	if theClient == nil {
		return
	}

	ab, err := json.Marshal(d)
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
