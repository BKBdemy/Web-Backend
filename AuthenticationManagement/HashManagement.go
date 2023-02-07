package AuthenticationManagement

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

var argon2settings = argon2Params{
	memory:      256 * 1024, // 256 MB memory cost
	iterations:  12,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

// TODO: RATELIMITING
func (am AuthenticationService) ComparePasswords(hashedPassword string, password string) (bool, error) {
	params, salt, hash, err := decodeHash(hashedPassword)
	if err != nil {
		return false, err
	}

	dbHash := argon2.IDKey([]byte(password), salt, params.iterations, params.memory, params.parallelism, params.keyLength)

	if subtle.ConstantTimeCompare(hash, dbHash) == 1 {
		return true, nil
	}

	return false, nil
}

func (am AuthenticationService) HashPassword(password string) (string, error) {
	// generate salt
	salt := make([]byte, argon2settings.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// hash password
	hashedPassword := argon2.IDKey([]byte(password), salt, argon2settings.iterations, argon2settings.memory, argon2settings.parallelism, argon2settings.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hashedPassword)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, argon2settings.memory, argon2settings.iterations, argon2settings.parallelism, b64Salt, b64Hash)
	return encodedHash, nil
}

var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func decodeHash(encodedHash string) (p *argon2Params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &argon2Params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
