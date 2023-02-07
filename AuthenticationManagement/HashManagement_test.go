package AuthenticationManagement_test

import (
	"EntitlementServer/AuthenticationManagement"
	"strings"
	"testing"
)

func TestAuthenticationService_HashPassword(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "admin",
			input: "admin",
		},
		{
			name:  "verylonginput",
			input: "verylonginputverylonginputverylonginputverylonginputverylonginputverylonginputverylonginput",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			am := AuthenticationManagement.AuthenticationService{}
			got, err := am.HashPassword(test.input)
			if err != nil {
				t.Errorf("HashPassword() error = %v", err)
				return
			}

			// has to match argon2id specification
			if !strings.HasPrefix(got, "$argon2id$v=") {
				t.Errorf("HashPassword() = %v, want $argon2id$v=", got)
			}
		})
	}
}

func TestAuthenticationService_ComparePasswords(t *testing.T) {
	tests := []struct {
		name  string
		hash  string
		plain string
		want  bool
	}{
		// TODO: test cases
		{
			name:  "admin",
			hash:  "$argon2id$v=19$m=65536,t=6,p=1$dGVzdHRlc3Q$TvWngSaVheFLS5b6VzkFgQ21gscr2YBdH5tjum0FpqA",
			plain: "admin",
			want:  true,
		},
		{
			name:  "verylonginput",
			hash:  "$argon2id$v=19$m=65564,t=6,p=1$dGVzdHRlc3Q$+uNyHlrS4YsA576bBiyfP9PkJ9eorOUEdRUa7GkrW5M",
			plain: "verylonginputverylonginputverylonginputverylonginputverylonginputverylonginputverylonginput",
			want:  true,
		},
		{
			name:  "verydifferentparameters",
			hash:  "$argon2id$v=19$m=64566,t=3,p=1$dGVzdHRlc3Q$mYHs9BtWbXktdizV62GEqH82GgYKTo9XMayNRyKBt6w",
			plain: "admin123",
			want:  true,
		},
		{
			name:  "wrongpassword",
			hash:  "$argon2id$v=19$m=65536,t=6,p=1$dGVzdHRlc3Q$TvWngSaVheFLS5b6VzkFgQ21gscr2YBdH5tjum0FpqA",
			plain: "admin1254ewt5w45t3",
			want:  false,
		},
		{
			name:  "wronghash_wrongparams",
			hash:  "$argon2id$v=19$m=64566,t=3,p=1$dGVzdHRlc3Q$mYHs9BtWbXktdizV62GEqH82GgYKTo9XMayNRyKBt6w",
			plain: "verylonginputverylonginputverylonginputverylonginputverylonginputverylonginputverylonginput",
			want:  false,
		},
		{
			name:  "empty_hash",
			hash:  "",
			plain: "verylonginputverylonginputverylonginputverylonginputverylonginputverylonginputverylonginput",
			want:  false,
		},
		{
			name:  "empty_plain",
			hash:  "$argon2id$v=19$m=65536,t=6,p=1$dGVzdHRlc3Q$TvWngSaVheFLS5b6VzkFgQ21gscr2YBdH5tjum0FpqA",
			plain: "",
			want:  false,
		},
		{
			name:  "malformed_hash",
			hash:  "$argtwteon2id$v=19$m=655tgse36,t=6,p=1$dGVzdHRlc3Q$TvWngSaVheFLS5b6VzkFgQ21gscr2YBdH5tjum0FpqA",
			plain: "verylonginputverylonginputverylonginputverylonginputverylonginputverylonginputverylonginput",
			want:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			am := AuthenticationManagement.AuthenticationService{}
			got, err := am.ComparePasswords(test.hash, test.plain)
			if err != nil {
				if !test.want {
					return
				}
				t.Errorf("ComparePasswords() error = %v", err)
				return
			}
			if got != test.want {
				t.Errorf("ComparePasswords() = %v, want %v", got, test.want)
			}
		})
	}
}
