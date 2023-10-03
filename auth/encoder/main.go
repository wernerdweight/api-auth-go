package encoder

import (
	"github.com/alexedwards/argon2id"
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(apiUser contract.ApiUserInterface, password string) *contract.AuthError {
	encryptedPassword := apiUser.GetPassword()
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if nil != err {
		match, err := argon2id.ComparePasswordAndHash(password, encryptedPassword)
		if !match || nil != err {
			return contract.NewAuthError(contract.InvalidCredentials, nil)
		}
	}
	return nil
}

func EncryptPassword(plainPassword string) (string, error) {
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if nil != err {
		return "", err
	}
	return string(encryptedPassword), nil
}
