package es

import (
	"context"
	"errors"
	"io"
	"project/pb"
	"project/util"
	"project/zj"
	"strings"
	"time"

	"google.golang.org/protobuf/proto"
)

// LastItem ...
func LastItem() {

	query := `{
		"size": 100,
		"sort": [
			{
				"ts": {
					"order": "desc"
				}
			}
		]
	}`
	zj.J(query)

	/*
		ab, err := Search(query)
		if err != nil {
			zj.W(`fail`, err)
			return
		}

		util.WriteFile(`search.json`, ab)

		zj.J(string(ab))
	*/
}

// Search ...
func Search(query string, m proto.Message) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	re, err := theClient.Search(
		theClient.Search.WithContext(ctx),
		theClient.Search.WithIndex(indexNameAll()),
		theClient.Search.WithBody(strings.NewReader(query)),
		theClient.Search.WithTrackTotalHits(false),
	)
	defer cancel()
	if err != nil {
		return err
	}
	defer re.Body.Close()

	ab, err := io.ReadAll(re.Body)
	if err != nil {
		return err
	}

	et := &pb.EsErrorTry{}
	util.JSONUnmarshal(ab, et)
	e := et.GetError()
	if e != nil {
		return errors.New(e.Reason)
	}

	util.WriteFile(`last-search.json`, ab)

	return util.JSONUnmarshal(ab, m)
}
