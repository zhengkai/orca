package es

import (
	"fmt"
	"project/config"
	"project/zj"
	"strings"
	"time"

	_ "embed" //

	"github.com/zhengkai/zu"
)

//go:embed tpl/mapping.json
var indexMapping string

func indexName(ts uint32) string {

	index := `orca-metrics`
	if !config.Prod {
		index = `dev-` + index
	}

	index = fmt.Sprintf(`%s-%s`, index, time.Unix(int64(ts), 0).Format(`2006.01.02`))

	return index
}

func createIndex() {

	ts := zu.TS()

	zj.J(`index name:`, indexName(ts))

	mapping(ts)
	mapping(ts + 86400)
	go func() {
		for {
			time.Sleep(time.Hour * 3)
			mapping(zu.TS() + 86400)
		}
	}()
}

func mapping(ts uint32) {
	theClient.Indices.Create(
		indexName(ts),
		theClient.Indices.Create.WithBody(strings.NewReader(indexMapping)),
	)
}
