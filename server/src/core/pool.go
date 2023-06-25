package core

import (
	"net/http"
	"project/pb"
	"time"
)

func (c *Core) add(req *pb.Req, hr *http.Request) (pr *row, cached bool) {

	hash := req.Hash()

	c.mux.Lock()
	pr, ok := c.pool[hash]
	if ok {
		// zj.F(`hit %x`, hash)
		c.mux.Unlock()
		cached = true
		return
	}

	c.serial++
	pr = &row{
		serial: c.serial,
		hr:     hr,
		hash:   hash,
		req:    req,
		t:      time.Now(),
		core:   c,
	}
	pr.mux.Lock()
	go pr.run()
	c.pool[hash] = pr
	c.mux.Unlock()
	return
}

func (c *Core) delete(r *row) {
	c.mux.Lock()
	pr, ok := c.pool[r.hash]
	if ok && pr == r {
		delete(c.pool, r.hash)
	}
	c.mux.Unlock()
}
