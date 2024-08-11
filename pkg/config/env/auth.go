package env

import (
	"fmt"
	"os"

	"github.com/andredubov/golibs/pkg/config"
)

const (
	passwordSaltEnvName = "PASSWORD_SALT"
)

type authConfig struct {
	passwordSalt string
}

// NewAuthConfig returns an intance of authConfig struct
func NewAuthConfig() (config.AuthConfing, error) {
	passwordSalt := os.Getenv(passwordSaltEnvName)
	if len(passwordSalt) == 0 {
		return nil, fmt.Errorf("%s", "password salt not found")
	}

	return &authConfig{
		passwordSalt,
	}, nil
}

// PasswordSalt returns password salt
func (cfg *authConfig) PasswordSalt() string {
	return cfg.passwordSalt
}
