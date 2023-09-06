package security

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) error {
	// TODO: Set example variable
	c.Set("example", "12345")

	// TODO: require config provider to retrieve config values by name

	// TODO: check if this endpoint should even be authenticated (regex?)
	// TODO: check if this endpoint should be authenticated by api client and user (1) or just by an api key (2)
	// TODO: 1
	// TODO: authenticate api client (id and secret)
	// TODO: if scope access model is enabled, check if client (and user) is allowed to access this endpoint
	// TODO: check client scope (and FUP)
	// TODO: - if client is not allowed to access this endpoint (omitted or false), return respective error
	// TODO: - if client is allowed to access this endpoint freely (true), return nil (done)
	// TODO: - if client is allowed to access this endpoint on-behalf ("on-behalf"), continue
	// TODO: authenticate api user (user token)
	// TODO: check user scope (and FUP)
	// TODO: - if user is not allowed to access this endpoint (omitted or false), return respective error
	// TODO: - if user is allowed to access this endpoint freely (true), return nil (done)

	// TODO: 2
	// TODO: authenticate user by api key (use some default api client)
	// TODO: check scope access (and FUP)

	return errors.New("not implemented")
}
