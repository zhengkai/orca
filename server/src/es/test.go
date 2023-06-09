package es

import (
	"fmt"
	"math/rand"
	"project/pb"
	"time"

	"github.com/zhengkai/zu"
)

// Test ...
func Test() {

	for {

		minute := time.Now().Minute() - 30
		if minute < 0 {
			minute = -minute
		}
		minute += 3

		modelChoose := []string{
			`gpt-3.5-turbo-0301`,
			`gpt-4-0314`,
			`text-embedding-ada-002-v2`,
		}
		model := modelChoose[rand.Intn(len(modelChoose))]

		d := &pb.EsMetrics{
			ID: fmt.Sprintf(`chatcmpl-%d`, zu.MS()),
			Token: &pb.EsMetricsToken{
				Total:      uint32(rand.Intn(minute * 10)),
				Completion: uint32(rand.Intn(minute * 10)),
				Prompt:     uint32(rand.Intn(minute * 10)),
			},
			Cached:   true,
			Ip:       `127.0.0.1`,
			Model:    model,
			Key:      `zhengkai.orca`,
			ReqBytes: uint32(rand.Intn(minute * 57)),
			RspBytes: uint32(rand.Intn(minute * 21)),
			Ts:       zu.MS(),
		}
		Insert(d)

		time.Sleep(time.Second / 2)
	}
}
