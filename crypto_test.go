package fftool

import (
	"testing"
)

// TestGenerateCrypto ...
func TestGenerateCrypto(t *testing.T) {
	type args struct {
		ssl   *OpenSSL
		useIV bool
	}
	tests := []struct {
		name string
		args args
		want *Crypto
	}{
		{
			name: "crypto1",
			args: args{
				ssl:   NewOpenSSL(),
				useIV: true,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateCrypto(tt.args.ssl, tt.args.useIV); got.Error() != nil {
				t.Errorf("GenerateCrypto() = %v", got.err)
			}
		})
	}
}

// TestCrypto_SaveKey ...
func TestCrypto_SaveKey(t *testing.T) {
	type fields struct {
		err         error
		KeyInfoPath string
		Key         string
		KeyPath     string
		UseIV       bool
		IV          string
		URL         string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "savekey",
			fields: fields{
				err:         nil,
				KeyInfoPath: DefaultKeyInfoName,
				Key:         "",
				KeyPath:     DefaultKeyName,
				UseIV:       false,
				IV:          "",
				URL:         "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GenerateCrypto(NewOpenSSL(), tt.fields.UseIV)
			c.KeyPath = tt.fields.KeyPath
			if err := c.SaveKey(); (err != nil) != tt.wantErr {
				t.Errorf("SaveKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCrypto_SaveKeyInfo ...
func TestCrypto_SaveKeyInfo(t *testing.T) {
	type fields struct {
		err         error
		KeyInfoPath string
		Key         string
		KeyPath     string
		UseIV       bool
		IV          string
		URL         string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "savekeyinfo",
			fields: fields{
				err:         nil,
				KeyInfoPath: DefaultKeyInfoName,
				Key:         "",
				KeyPath:     DefaultKeyName,
				UseIV:       true,
				IV:          "",
				URL:         "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GenerateCrypto(NewOpenSSL(), tt.fields.UseIV)
			c.KeyPath = tt.fields.KeyPath
			c.KeyInfoPath = tt.fields.KeyInfoPath
			if err := c.SaveKeyInfo(); (err != nil) != tt.wantErr {
				t.Errorf("SaveKeyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestLoadCrypto ...
func TestLoadCrypto(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Crypto
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				path: `D:\workspace\golang\project\go-fftool\output_key\5f2e2a10-f6df-450f-834d-97c17cc70396\m3u8_key_info`,
			},
			wantC:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC := LoadCrypto(tt.args.path)
			if (gotC.Error() != nil) != tt.wantErr {
				t.Errorf("LoadCrypto() error = %v, wantErr %v", gotC.Error(), tt.wantErr)
				return
			}
			if gotC != nil {
				t.Logf("LoadCrypto() gotC = %+v", gotC)
			}
		})
	}
}
