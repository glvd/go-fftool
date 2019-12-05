package fftool

import (
	"strings"
)

// OpenSSL ...
type OpenSSL struct {
	cmd  *Command
	Name string
}

// NewOpenSSL ...
func NewOpenSSL() *OpenSSL {
	return &OpenSSL{Name: "openssl"}
}

func (ssl *OpenSSL) init() {
	if ssl.cmd == nil {
		ssl.cmd = New(ssl.Name)
	}
}

// Run ...
func (ssl *OpenSSL) Run(args string) (string, error) {
	ssl.init()
	return ssl.cmd.Run(args)
}

// Base64 ...
func (ssl *OpenSSL) Base64(size int) string {
	run, err := ssl.Run(formatArgs("rand,-base64,%d", size))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(run)
}

// Hex ...
func (ssl *OpenSSL) Hex(size int) string {
	run, err := ssl.Run(formatArgs("rand,-hex,%d", size))
	if LogError(err) {
		return ""
	}
	return strings.TrimSpace(run)
}

// HLSCrypto ...
func (ssl *OpenSSL) HLSCrypto() *Crypto {
	ssl.init()
	return &Crypto{
		KeyInfoPath: "",
		Key:         ssl.Base64(32),
		IV:          "",
		URL:         "",
	}
}
