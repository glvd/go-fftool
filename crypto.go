package fftool

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//Crypto ...
type Crypto struct {
	err         error
	KeyInfoPath string
	Key         string
	IV          string
	URL         string
}

// GenerateCrypto ...
func GenerateCrypto() *Crypto {
	ssl := NewOpenSSL()
	c := ssl.HLSCrypto()

	run, err := ssl.Run("-hex,16")
	if err != nil {
		return nil
	}
	c.IV = run

	split, _ := filepath.Split(path)
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

// SaveToFile ...
func (c *Crypto) SaveToFile(path string) {

}

func outputKeyInfoString(url, key, iv string) []byte {

}
