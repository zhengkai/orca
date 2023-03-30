package core

import (
	"net/http"
	"project/pb"
	"project/zj"
	"time"
)

func (c *Core) add(req *pb.Req, hr *http.Request) (pr *row) {

	hash := req.Hash()

	c.mux.Lock()
	pr, ok := c.pool[hash]
	if ok {
		zj.F(`hit %x`, hash)
		c.mux.Unlock()
		return
	}

	pr = &row{
		hr:   hr,
		hash: hash,
		req:  req,
		t:    time.Now(),
	}
	pr.mux.Lock()
	pr.run()
	c.pool[hash] = pr
	c.mux.Unlock()
	return
}
