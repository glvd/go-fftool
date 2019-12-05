package fftool

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

// HLSCrypto ...
func (ssl *OpenSSL) HLSCrypto() *Crypto {
	ssl.init()
	run, err := ssl.cmd.Run("-base64,32")
	if err != nil {
		return nil
	}
	return &Crypto{
		KeyInfoPath: "",
		Key:         run,
		IV:          "",
		URL:         "",
	}
}
