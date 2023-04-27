package core

import (
	"net/http"
	"project/pb"
	"project/zj"
	"time"
)

func (c *Core) add(req *pb.Req, hr *http.Request) (pr *row, cached bool) {

	hash := req.Hash()

	c.mux.Lock()
	pr, ok := c.pool[hash]
	if false && ok {
		zj.F(`hit %x`, hash)
		c.mux.Unlock()
		cached = true
		return
	}

	pr = &row{
		hr:   hr,
		hash: hash,
		req:  req,
		t:    time.Now(),
	}
	pr.mux.Lock()
	go pr.run()
	c.pool[hash] = pr
	c.mux.Unlock()
	return
}
