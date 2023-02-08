package LicenseKeyManager

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
)

type LicenseKey struct {
	KeyID         string
	EncryptionKey EncryptionKey
	ProductID     int
}
type MultiFormatKey interface {
	GetBase64() (string, error)
	GetHex() string
	GetBytes() ([]byte, error)
}

type EncryptionKey struct {
	Hex string
}

func (e EncryptionKey) GetBase64() (string, error) {
	bytes, err := e.GetBytes()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (e EncryptionKey) GetHex() string {
	return e.Hex
}

func (e EncryptionKey) GetBytes() ([]byte, error) {
	dst := make([]byte, len(e.Hex))
	_, err := hex.Decode(dst, []byte(e.Hex))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func Base64ToHex(input string) (string, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(strings.Trim(input, " "))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
