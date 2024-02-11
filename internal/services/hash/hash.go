package hash

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func (h *Hasher) HashPassword(plainPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), err
}

func (h *Hasher) VerifyHashPassword(passwordHash, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword))
}
