package encoder

import (
	"github.com/wernerdweight/api-auth-go/auth/contract"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(apiUser contract.ApiUserInterface, password string) *contract.AuthError {
	encryptedPassword := apiUser.GetPassword()
	err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if nil != err {
		argonPassword := argon2.IDKey([]byte(password), []byte(""), 1, 64*1024, 4, 32)
		if string(argonPassword) != encryptedPassword {
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
