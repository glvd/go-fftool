package fftool

import (
	"bytes"
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
	KeyPath     string
	UseIV       bool
	IV          string
	URL         string
}

// GenerateCrypto ...
func GenerateCrypto(ssl *OpenSSL, useIV bool) *Crypto {
	c := Crypto{
		err:         nil,
		KeyInfoPath: "",
		Key:         ssl.Base64(32),
		UseIV:       useIV,
		IV:          "",
		URL:         "",
	}

	if useIV {
		c.IV = ssl.Hex(16)
	}

	if c.Key == "" || (useIV && c.IV == "") {
		c.err = fmt.Errorf("generate crypto error(key:%v,useIV:%v,iv:%v", c.Key, useIV, c.IV)
	}

	return &c
}

// SaveKey ...
func (c *Crypto) SaveKey() error {
	split, _ := filepath.Split(c.KeyPath)
	stat, err := os.Stat(split)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(split, 0755)
		}
	}
	if err == nil && !stat.IsDir() {
		return fmt.Errorf("wrong target path:%s", split)
	}
	return ioutil.WriteFile(c.KeyPath, []byte(c.Key), 0755)
}

// SaveKeyInfo ...
func (c *Crypto) SaveKeyInfo() error {
	split, _ := filepath.Split(c.KeyPath)
	stat, err := os.Stat(split)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(split, 0755)
		}
	}
	if err == nil && !stat.IsDir() {
		return fmt.Errorf("wrong target path:(%v)", split)
	}

	if c.URL == "" {
		return fmt.Errorf("wrong url address:(%v)", c.URL)
	}

	buff := bytes.NewBufferString(c.URL)
	buff.WriteString("\n")
	buff.WriteString(c.KeyPath)
	buff.WriteString("\n")
	if c.UseIV {
		buff.WriteString(c.IV)
	}
	return ioutil.WriteFile(c.KeyInfoPath, buff.Bytes(), 0755)
}

// Error ...
func (c *Crypto) Error() error {
	return c.err
}

func outputKeyInfoString(url, key, iv string) []byte {
	return nil
}
