package fftool

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/goextension/log"
	"io"
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

// LoadCrypto ...
func LoadCrypto(path string) (c *Crypto) {
	c = &Crypto{}
	path = abs(path)
	open, err := os.Open(path)
	if err != nil {
		c.err = Err(err, "crypto open")
	}
	reader := bufio.NewReader(open)

	c.KeyInfoPath = path
	line, _, err := reader.ReadLine()
	for i := 0; err == nil; i++ {
		switch i {
		case 0:
			c.KeyPath = string(line)

			if !filepath.IsAbs(c.KeyPath) {
				c.KeyPath = filepath.Join(filepath.Dir(path), c.KeyPath)
			}

			key, err := ioutil.ReadFile(c.KeyPath)
			if err != nil {
				c.err = Err(err, "crypto read key")
				return
			}
			c.Key = string(key)

		case 1:
			c.URL = string(line)
		case 2:
			c.IV = string(line)
			if c.IV != "" {
				c.UseIV = true
			}
		default:
			log.Infow("crypto", "read", string(line))
		}
		line, _, err = reader.ReadLine()
	}

	if err != io.EOF {
		c.err = Err(err, "crypto read info")
	}
	return
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
