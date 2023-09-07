package vertexai

import (
	"encoding/json"
	"project/pb"
	"project/util"
	"time"

	"github.com/zhengkai/coral/v2"
)

var textCache = coral.NewLRU(loadTextForCoral, 1000, 100)

const textModel = `text-bison@001`

type textKey struct {
	paramKey
	Prompt string `json:"system"`
}

// Text ...
func Text(req *pb.VaTextReq) (*Rsp, error) {

	k := textKey{
		Prompt: req.Prompt,
	}
	k.load(req.Param)

	if req.NoCache {
		textCache.Delete(k)
	} else {
		ab, err := util.ReadFile(cacheFile(`text`, k) + `.json`)
		if err == nil && len(ab) > 2 {
			rsp := &Rsp{}
			err = json.Unmarshal(ab, rsp)
			if err == nil {
				return rsp, nil
			}
		}
	}

	return textCache.Get(k)
}

func loadText(k textKey) (*Rsp, error) {

	m := map[string]any{
		`prompt`: k.Prompt,
	}

	req, err := buildReq(k.paramKey, m, textModel)
	if err != nil {
		return nil, err
	}

	rsp, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	go func() {
		rsp.save(cacheFile(`text`, k))
	}()
	return rsp, nil
}

func loadTextForCoral(k textKey) (*Rsp, *time.Time, error) {
	r, err := loadText(k)
	if err != nil {
		return nil, nil, err
	}
	return r, nil, nil
}
