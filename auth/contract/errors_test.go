package contract

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAuthError_Error(t *testing.T) {
	assertion := assert.New(t)

	err := NewAuthError(InvalidCredentials, nil)
	assertion.NotNil(err)
	assertion.Equal("invalid credentials", err.Err.Error())
	assertion.Nil(err.Payload)
	assertion.Equal(http.StatusUnauthorized, err.Status)

	err = NewAuthError(InvalidCredentials, "string")
	assertion.Equal("string", err.Payload)

	err = NewAuthError(InvalidCredentials, 123)
	assertion.Equal(123, err.Payload)

	err = NewAuthError(InvalidCredentials, map[string]float32{"test": 123})
	assertion.Equal(map[string]float32{"test": 123}, err.Payload)

	err = NewInternalError(DatabaseError, nil)
	assertion.Equal(http.StatusInternalServerError, err.Status)

	err = NewFUPError(DatabaseError, nil)
	assertion.Equal(http.StatusTooManyRequests, err.Status)
}
