package routes

import "github.com/gin-gonic/gin"

func resettingRequestHandler(c *gin.Context) {
	// TODO: get login information (email) from request body
	// TODO: validate login information (user exists)
	// TODO: check for recent requests (prevent spam)
	// TODO: generate reset token and set reset token date
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: call external service to send reset email (event)
}

func resettingResetHandler(c *gin.Context) {
	// TODO: get token from request params
	// TODO: validate token
	// TODO: get password from request body
	// TODO: validate password
	// TODO: set new password and reset token and request date
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: issue resetting done event
}
