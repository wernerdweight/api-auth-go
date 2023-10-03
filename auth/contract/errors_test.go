package contract

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthError_Error(t *testing.T) {
	assertion := assert.New(t)

	err := NewAuthError(InvalidCredentials, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid credentials", err.Err.Error())
	assertion.Nil(err.Payload)

	err = NewAuthError(InvalidCredentials, "string")
	assertion.Equal("string", err.Payload)

	err = NewAuthError(InvalidCredentials, 123)
	assertion.Equal(123, err.Payload)

	err = NewAuthError(InvalidCredentials, map[string]float32{"test": 123})
	assertion.Equal(map[string]float32{"test": 123}, err.Payload)
}
