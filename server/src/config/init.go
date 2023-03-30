package config

import (
	"os"
	"path/filepath"
)

func init() {

	Dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	list := map[string]*string{
		`OPENAI_API_KEY`: &OpenAIKey,
		`STATIC_DIR`:     &StaticDir,
		`WEB_ADDR`:       &WebAddr,
	}
	for k, v := range list {
		s := os.Getenv(k)
		if len(s) > 1 {
			*v = s
		}
	}
}
