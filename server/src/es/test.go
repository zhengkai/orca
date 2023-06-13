package es

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"project/pb"
	"time"

	"github.com/zhengkai/zu"
)

// Test ...
func Test() {

	for {

		ts := zu.MS()

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

		ipChoose := []string{
			`10.0.32.43`,
			`10.0.84.231`,
			`10.0.84.49`,
			`10.2.197.197`,
		}
		ip := ipChoose[rand.Intn(len(ipChoose))]

		keyChoose := []string{
			`zhengkai.orca`,
			`zhengkai.debug`,
		}
		key := keyChoose[rand.Intn(len(keyChoose))]

		hash := md5.Sum([]byte(fmt.Sprintf(`%s-%s-%s-%d`, model, ip, key, ts)))

		d := &pb.EsMetrics{
			ID: fmt.Sprintf(`chatcmpl-%d`, ts),
			Token: &pb.EsMetricsToken{
				Total:      uint32(rand.Intn(minute * 10)),
				Completion: uint32(rand.Intn(minute * 10)),
				Prompt:     uint32(rand.Intn(minute * 10)),
			},
			Cached:   rand.Intn(2) == 1,
			Ip:       ip,
			Model:    model,
			Key:      key,
			ReqBytes: uint32(rand.Intn(minute * 57)),
			RspBytes: uint32(rand.Intn(minute * 21)),
			Ts:       ts,
			Hash:     fmt.Sprintf(`%x`, hash),
		}
		Insert(d)

		time.Sleep(time.Second / 2)
	}
}
