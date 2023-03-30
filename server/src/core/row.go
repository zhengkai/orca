package core

import (
	"fmt"
	"project/pb"
	"project/zj"
	"sync"
	"time"
)

type row struct {
	hash [16]byte
	req  *pb.Req
	rsp  []byte
	err  error
	done bool
	t    time.Time
	mux  sync.RWMutex
}

func (pr *row) run() {

	pr.t = time.Now()

	s := fmt.Sprintf(`%x, %s`, pr.hash, pr.t.Format(`2006-01-02 15:04:05`))
	zj.J(`new`, s)

	pr.rsp = []byte(s)
	pr.done = true
	pr.mux.Unlock()
}

func (pr *row) wait() {
	if !pr.done {
		pr.mux.RLock()
		pr.mux.RUnlock()
	}
}
