package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"project/config"
	"project/util"
	"time"
)

func (pr *row) fetchRemote() (ab []byte, ok bool, err error) {

	r := pr.req
	var b bytes.Buffer
	pr.failLog = &b

	u, _ := url.Parse(config.OpenAIBase)
	u.Path = r.Path

	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader(r.Body))
	if err != nil {
		return
	}

	b.WriteString(pr.hr.URL.String())
	b.WriteString("\n\nreq header:\n\n")
	for k, v := range pr.hr.Header {
		fmt.Fprintf(&b, "\t%s: %v\n", k, v)
	}
	b.WriteString("\n")
	fmt.Fprintf(&b, "req body: %d\n\n", len(r.Body))
	if len(r.Body) > 0 {
		fmt.Fprintf(&b, "%s\n\n", r.Body)
	}

	req.Header.Set(`Content-Type`, `application/json`)
	req.Header.Set(`Authorization`, `Bearer `+config.OpenAIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(&b, "client.Do fail: %s\n", err.Error())
		return
	}

	if rsp.StatusCode >= 200 || rsp.StatusCode < 300 {
		// ok = true
	} else {
		err = fmt.Errorf(`status code fail: %d`, rsp.StatusCode)
		b.WriteString(err.Error())
	}

	b.WriteString("req header:\n\n")
	for k, v := range rsp.Header {
		fmt.Fprintf(&b, "\t%s: %v\n", k, v)
	}
	b.WriteString("\n")

	ab, err = io.ReadAll(rsp.Body)
	if err != nil {
		fmt.Fprintf(&b, "rsp body: %d %v\n\n", len(ab), err)
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
