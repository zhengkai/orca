package core

import (
	"bytes"
	"fmt"
	"net/http"
	"project/pb"
	"project/util"
	"sync"
	"time"
)

type row struct {
	serial int
	hash   [16]byte
	hr     *http.Request
	req    *pb.Req
	rsp    []byte
	err    error
	done   bool
	t      time.Time
	mux    sync.RWMutex
	log    *bytes.Buffer
	core   *Core
}

func (pr *row) suicide() {
	pr.core.delete(pr)
	pr.core = nil
}

func (pr *row) run() {

	pr.t = time.Now()

	pr.startLog()

	if pr.req.Method != http.MethodGet {
		// pr.mux.Unlock()
		// return
	}

	pr.rsp, pr.err = pr.fetchRemote()
	pr.done = true
	pr.mux.Unlock()

	// go pr.metrics()

	if pr.err == nil {
		// pr.failLog.Reset()
		// go writeFailLog(pr.hash, pr.log.Bytes())
		go pr.saveFile()
	} else {
		go writeFailLog(pr.hash, pr.log.Bytes())
		go pr.suicide()
	}
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

func (pr *row) startLog() {

	var b bytes.Buffer
	pr.log = &b

	b.WriteString(time.Now().Format("2006-01-02 15:04:05.000\n"))
	b.WriteString(pr.hr.Method + ` ` + pr.hr.URL.String())
	b.WriteString("\n\nreq header:\n\n")
	for k, v := range pr.hr.Header {
		fmt.Fprintf(&b, "\t%s: %v\n", k, v)
	}
	b.WriteString("\n")

	body := pr.req.Body
	fmt.Fprintf(&b, "req body: %d\n\n", len(body))
	if len(body) > 0 {
		b.Write(body)
		b.WriteString("\n\n")
	}
}
