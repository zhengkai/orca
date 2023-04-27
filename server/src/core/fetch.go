package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"project/config"
	"project/util"
	"project/zj"
	"time"
)

func (pr *row) fetchRemote() (ab []byte, ok bool, err error) {

	r := pr.req
	b := pr.log

	u, err := url.Parse(config.OpenAIBase + r.Url)
	if err != nil {
		return nil, false, err
	}
	zj.J(`real url`, u.String())

	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader(r.Body))
	if err != nil {
		return
	}

	req.Header.Set(`Content-Type`, r.ContentType)
	req.Header.Set(`Authorization`, `Bearer `+config.OpenAIKey)

	client := &http.Client{
		// Timeout: 30 * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(b, "client.Do fail: %s\n", err.Error())
		return
	}

	if rsp.StatusCode >= 200 || rsp.StatusCode < 300 {
		ok = true
	} else {
		err = fmt.Errorf(`status code fail: %d`, rsp.StatusCode)
		b.WriteString(err.Error())
	}

	b.WriteString("rsp header:\n\n")
	for k, v := range rsp.Header {
		fmt.Fprintf(b, "\t%s: %v\n", k, v)
	}
	b.WriteString("\n")

	ab, err = io.ReadAll(rsp.Body)
	fmt.Fprintf(b, "rsp body: %d %v\n\n", len(ab), err)
	if err == nil {
		b.Write(ab)
	}

	rsp.Body.Close()
	return
}

func writeFailLog(hash [16]byte, ab []byte) {
	date := time.Now().Format(`0102-150405`)
	file := fmt.Sprintf(`fail/%s-%x.txt`, date, hash)
	util.WriteFile(file, ab)
}
