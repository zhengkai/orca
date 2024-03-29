package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

func init() {

	Dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	LogDir = Dir + `/log`

	list := map[string]*string{
		`OPENAI_API_KEY`: &OpenAIKey,
		`ORCA_DIR`:       &StaticDir,
		`ORCA_WEB`:       &WebAddr,
		`ORCA_LOG`:       &LogDir,
		`ORCA_ES_ADDR`:   &ESAddr,
		`ORCA_ES_USER`:   &ESUser,
		`ORCA_ES_PASS`:   &ESPass,
		`ORCA_VA_TOKEN`:  &VAToken,
	}
	for k, v := range list {
		s := os.Getenv(k)
		if len(s) > 1 {
			*v = s
		}
	}

	_, err := url.Parse(OpenAIBase)
	if err != nil {
		fmt.Println(`OpenAI base URL is invalid.`)
		panic(err)
	}

	os.MkdirAll(StaticDir+`/tmp`, 0777)
}
