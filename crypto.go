package fftool

import (
	"fmt"
	"github.com/goextension/tool"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Crypto ...
type Crypto struct {
	KeyInfoPath string
	Key         string
	IV          string
	URL         string
}

// GenerateCrypto ...
func GenerateCrypto(path string) Crypto {
	key := tool.GenerateRandomString(16)
	iv := tool.GenerateRandomString(16)
	split, file := filepath.Split(path)
	stat, err := os.Stat(split)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(split, 0755)
		}
	}
	if err == nil && !stat.IsDir() {
		panic(fmt.Sprintf("wrong target path:%s", split))
	}

	ioutil.WriteFile(path)
}

func outputKeyInfoString(url, key, iv string) []byte {

}
