package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"project/config"
	"project/metrics"
	"project/pb"
	"project/util"
	"time"
)

func (pr *row) fetchRemote() (ab []byte, err error) {

	r := pr.req
	b := pr.log

	u, err := url.Parse(config.OpenAIBase + r.Url)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, u.String(), bytes.NewReader(r.Body))
	if err != nil {
		return
	}

	req.Header.Set(`Content-Type`, r.ContentType)
	req.Header.Set(`Authorization`, `Bearer `+config.OpenAIKey)

	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	t := time.Now()
	rsp, err := client.Do(req)
	costMs := uint32(time.Since(t) / time.Millisecond)
	metrics.OpenAITime(costMs)
	if err != nil {
		fmt.Fprintf(b, "client.Do fail: %s\n", err.Error())
		return
	}

	if rsp.StatusCode < 200 || rsp.StatusCode >= 300 {
		pr.httpCode = rsp.StatusCode
		err = fmt.Errorf(`status code fail: %d`, rsp.StatusCode)
		b.WriteString(err.Error())
		b.WriteString("\n\n")
	}

	b.WriteString("rsp header:\n\n")
	for k, v := range rsp.Header {
		fmt.Fprintf(b, "\t%s: %v\n", k, v)
	}
	b.WriteString("\n")

	ab, err = io.ReadAll(rsp.Body)
	fmt.Fprintf(b, "rsp body: %d %v\n\n", len(ab), err)
	b.Write(ab)

	rsp.Body.Close()

	metrics.Limit(rsp.Header)

	if err == nil {
		e := &pb.OpenAIError{}
		json.Unmarshal(ab, e)
		if e.GetError() != nil {
			err = fmt.Errorf(`openai error: %s`, e.GetError().GetMessage())
		}
	}

	return
}

func writeFailLog(hash [16]byte, ab []byte) {
	date := time.Now().Format(`0102/150405`)
	file := fmt.Sprintf(`fail/%s-%x.txt`, date, hash)
	os.MkdirAll(path.Dir(util.Static(file)), 0755)
	util.WriteFile(file, ab)
}
