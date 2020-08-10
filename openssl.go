package tool

import (
	"strings"
)

// OpenSSL ...
type OpenSSL struct {
	cmd  *Command
	name string
}

// NewOpenSSL ...
func NewOpenSSL() *OpenSSL {
	ssl := &OpenSSL{
		name: DefaultOpenSSLName,
	}
	ssl.cmd = NewCommand(ssl.name)
	return ssl
}

// Run ...
func (ssl *OpenSSL) Run(args string) (string, error) {
	return ssl.cmd.Run(args)
}

// Base64 ...
func (ssl *OpenSSL) Base64(size int) string {
	run, err := ssl.Run(formatArgs("rand,-base64,%d", size))
	if LogError(err) {
		panic(err)
	}
	return strings.TrimSpace(run)
}

// Hex ...
func (ssl *OpenSSL) Hex(size int) string {
	run, err := ssl.Run(formatArgs("rand,-hex,%d", size))
	if LogError(err) {
		panic(err)
	}
	return strings.TrimSpace(run)
}
