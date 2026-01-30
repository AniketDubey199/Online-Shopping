package utils

import "golang.org/x/crypto/bcrypt"

func Hashing(pass string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashed), err
}
func ComparePassword(userpass string, authpass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(userpass), []byte(authpass))
	return err
}
