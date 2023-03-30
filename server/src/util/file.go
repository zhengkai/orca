package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"project/config"
	"project/zj"
	"strings"

	"github.com/zhengkai/zu"
	"google.golang.org/protobuf/proto"
)

// DownloadFunc ...
type DownloadFunc func(url string) (ab []byte, err error)

// CacheName ...
func CacheName(hash [16]byte) string {
	s := fmt.Sprintf(`cache/%x/%x/%x`, hash[0:2], hash[2:4], hash[4:])
	os.MkdirAll(StaticFile(filepath.Dir(s)), 0755)
	return s
}

// FileExists ...
func FileExists(filename string) bool {
	filename = fmt.Sprintf(`%s/%s`, config.StaticDir, filename)
	return zu.FileExists(filename)
}

// IsURL ...
func IsURL(s string) bool {
	return strings.HasPrefix(s, `https://`) || strings.HasPrefix(s, `http://`)
}

// ReadFile ...
func ReadFile(file string) (ab []byte, err error) {
	ab, err = os.ReadFile(StaticFile(file))
	zj.Watch(&err)
	return
}

// StaticFile ...
func StaticFile(file string) string {
	file = strings.TrimPrefix(file, config.StaticDir+`/`)
	return fmt.Sprintf(`%s/%s`, config.StaticDir, file)
}

// SaveData ...
func SaveData(name string, p proto.Message) (err error) {

	defer zj.Watch(&err)

	ab, err := proto.Marshal(p)
	if err == nil {
		WriteFile(name+`.pb`, ab)
	}

	ab, err = json.Marshal(p)
	if err == nil {
		WriteFile(name+`.json`, ab)
	}

	return
}

// WriteFile ...
func WriteFile(file string, ab []byte) (err error) {

	file = StaticFile(file)

	defer zj.Watch(&err)

	f, err := os.CreateTemp(config.StaticDir+`/tmp`, `wr-*`)
	if err != nil {
		return
	}
	tmpName := f.Name()
	_, err = f.Write(ab)
	if err != nil {
		return
	}

	err = os.Chmod(tmpName, 0644)
	if err != nil {
		return
	}

	err = os.Rename(tmpName, file)
	if err != nil {
		return
	}

	return
}
