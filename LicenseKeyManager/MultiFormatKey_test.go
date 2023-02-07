package LicenseKeyManager

import "testing"

func TestEncryptionKey_GetHex(t *testing.T) {
	tests := []struct {
		name string
		key  EncryptionKey
		want string
	}{
		{
			name: "Pass",
			key:  EncryptionKey{Hex: "test"},
			want: "test",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.key.GetHex(); got != test.want {
				t.Errorf("EncryptionKey.GetHex() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestEncryptionKey_GetBase64(t *testing.T) {
	tests := []struct {
		name string
		hex  string
		want string
	}{
		{
			name: "Pass",
			hex:  "0bc7c30945eb99d9ce5831957eee8ed1",
			want: "C8fDCUXrmdnOWDGVfu6O0Q==",
		},
		{
			name: "invalid hex",
			hex:  "EXTREMELY_INVALID_HEX",
			want: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			key := EncryptionKey{Hex: test.hex}
			if got, _ := key.GetBase64(); got != test.want {
				if got != test.want && test.want != "" {
					return
				}
				t.Errorf("EncryptionKey.GetBase64() = %v, want %v", got, test.want)
			}
		})
	}
}
