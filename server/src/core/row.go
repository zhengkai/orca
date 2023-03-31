package core

import (
	"bytes"
	"fmt"
	"net/http"
	"project/pb"
	"project/util"
	"project/zj"
	"sync"
	"time"
)

type row struct {
	hash    [16]byte
	hr      *http.Request
	req     *pb.Req
	rsp     []byte
	err     error
	done    bool
	t       time.Time
	mux     sync.RWMutex
	failLog *bytes.Buffer
}

func (pr *row) run() {

	pr.t = time.Now()

	s := fmt.Sprintf(`%x, %s`, pr.hash, pr.t.Format(`2006-01-02 15:04:05`))
	zj.J(`new`, s)

	var ok bool
	pr.rsp, ok, pr.err = pr.fetchRemote()
	if pr.err == nil && ok {
		pr.failLog.Reset()
	} else {
		go writeFailLog(pr.hash, pr.failLog.Bytes())
	}

	go pr.saveFile()
	// go pr.metrics()

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
	rspFile := rspCacheFile(pr.req)
	util.WriteFile(rspFile, pr.rsp)
}
