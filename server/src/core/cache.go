package core

import (
	"fmt"
	"project/pb"
	"project/util"
)

func tryCache(p *pb.Req) ([]byte, bool) {

	file := rspCacheFile(p)

	ab, err := util.ReadFile(file)
	if err != nil {
		return nil, false
	}

	return ab, true
}

func rspCacheFile(r *pb.Req) string {
	return util.CacheName(r.Hash()) + `-rsp.json`
}

func cacheFile(hash [16]byte) string {
	s := fmt.Sprintf(`cache/%02x/%02x/%02x/%x`, hash[0], hash[1], hash[2], hash[3:])
	return s
}
