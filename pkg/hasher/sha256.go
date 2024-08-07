package hasher

import (
	"crypto/sha256"
	"fmt"
)

// sha256Hasher uses SHA1 to hash passwords with provided salt.
type sha256Hasher struct {
	salt string
}

// NewSHA256Hasher returns an instance of sha256 hasher.
func NewSHA256Hasher(salt string) PasswordHasher {
	return &sha256Hasher{salt: salt}
}

// HashAndSalt returns the sha256 hash of the password at the given salt.
func (h *sha256Hasher) HashAndSalt(password string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}

// ComparePasswords compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure.
func (h *sha256Hasher) ComparePasswords(hashedPassword string, plainPassword string) error {
	const op = "sha256Hasher.ComparePasswords"
	passwordHash, _ := h.HashAndSalt(plainPassword)

	if passwordHash != hashedPassword {
		return fmt.Errorf("%s: hashed password and plain password don't match", op)
	}

	return nil
}
