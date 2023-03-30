package core

import (
	"fmt"
	"net/http"
	"project/pb"
	"project/util"
	"project/zj"
	"sync"
	"time"
)

type row struct {
	hash [16]byte
	hr   *http.Request
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

	pr.rsp, pr.err = fetchRemote(pr.req)
	go pr.saveFile()
	go pr.metrics()

	pr.done = true
	pr.mux.Unlock()
}

func (pr *row) wait() {
	if !pr.done {
		pr.mux.RLock()
		pr.mux.RUnlock()
	}
}

func (pr *row) saveFile() {
	rspFile := util.CacheName(pr.req.Hash()) + `-rsp.json`
	if !util.FileExists(rspFile) {
		util.WriteFile(rspFile, pr.rsp)
		zj.J(rspFile)
	}
}

func (pr *row) metrics() {
}
