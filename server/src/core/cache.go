package core

import (
	"project/pb"
	"project/util"
)

func tryCache(p *pb.Req) ([]byte, bool) {

	file := rspCacheFile(p)
	if !util.FileExists(file) {
		return nil, false
	}

	ab, err := util.ReadFile(file)
	if err != nil {
		return nil, false
	}

	return ab, true
}

func rspCacheFile(r *pb.Req) string {
	return util.CacheName(r.Hash()) + `-rsp.json`
}
