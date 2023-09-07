package vertexai

import (
	"encoding/json"
	"project/pb"
	"project/util"
	"time"

	"github.com/zhengkai/coral/v2"
)

var chatCache = coral.NewLRU(loadChatForCoral, 1000, 100)

const chatModel = `chat-bison@001`

type chatKey struct {
	paramKey
	System string `json:"system"`
	User   string `json:"user"`
}

// Chat ...
func Chat(req *pb.VaChatReq) (*Rsp, error) {

	k := chatKey{
		System: req.System,
		User:   req.User,
	}
	k.load(req.Param)

	if req.NoCache {
		chatCache.Delete(k)
	} else {
		ab, err := util.ReadFile(cacheFile(`chat`, k) + `.json`)
		if err == nil && len(ab) > 2 {
			rsp := &Rsp{}
			err = json.Unmarshal(ab, rsp)
			if err == nil {
				return rsp, nil
			}
		}
	}

	return chatCache.Get(k)
}

func loadChat(k chatKey) (*Rsp, error) {

	m := map[string]any{
		`context`: k.System,
	}
	if k.User != `` {
		m[`messages`] = []any{
			map[string]any{
				`author`:  `user`,
				`content`: k.User,
			},
		}
	}

	req, err := buildReq(k.paramKey, m, chatModel)
	if err != nil {
		return nil, err
	}

	rsp, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	go func() {
		rsp.save(cacheFile(`chat`, k))
	}()
	return rsp, nil
}

func loadChatForCoral(k chatKey) (*Rsp, *time.Time, error) {
	r, err := loadChat(k)
	if err != nil {
		return nil, nil, err
	}
	return r, nil, nil
}
