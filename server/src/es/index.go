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

func indexNameBase() string {
	index := `orca-metrics`
	if !config.Prod {
		index = `dev-` + index
	}
	return index
}

func indexName(ts uint32) string {
	return fmt.Sprintf(`%s-%s`, indexNameBase(), time.Unix(int64(ts), 0).Format(`2006.01.02`))
}

func indexNameAll() string {
	return indexNameBase() + `-*`
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
