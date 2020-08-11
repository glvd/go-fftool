package tool

import "testing"

// TestOpenSSL_Hex ...
func TestOpenSSL_Hex(t *testing.T) {
	type fields struct {
		cmd  *Command
		Name string
	}
	type args struct {
		size int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "hex",
			fields: fields{
				cmd:  nil,
				Name: "openssl",
			},
			args: args{
				size: 16,
			},
			want: "5ffdbf238234107e7e53d864a44f27ca",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssl := NewOpenSSL()
			got := ssl.Hex(tt.args.size)
			if len(got) != len(tt.want) {
				t.Errorf("Hex() = %v, want %v", len(got), len(tt.want))
			}
			t.Log(got)
		})
	}
}

// TestOpenSSL_Base64 ...
func TestOpenSSL_Base64(t *testing.T) {
	type fields struct {
		cmd  *Command
		Name string
	}
	type args struct {
		size int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "base64",
			fields: fields{
				cmd:  nil,
				Name: "openssl",
			},
			args: args{
				size: 32,
			},
			want: "+CrOytKp+DtxK5pZcOIQbfBPzyxCSuXNwFop5X8ZacI=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssl := &OpenSSL{
				cmd:  tt.fields.cmd,
				name: tt.fields.Name,
			}
			if got := ssl.Base64(tt.args.size); len(got) != len(tt.want) {
				t.Errorf("Base64() = %v, want %v", got, tt.want)
			}
		})
	}
}
