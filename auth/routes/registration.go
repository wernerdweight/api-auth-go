package routes

import "github.com/gin-gonic/gin"

func registrationRequestHandler(c *gin.Context) {
	// TODO: get login information (email) from request body
	// TODO: validate login information
	// TODO: check for duplicates
	// TODO: get password from request body
	// TODO: validate password
	// TODO: create new user (make sure to set active to false, set confirmation token)
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: call external registration service to handle confirm callback (event)
}

func registrationConfirmHandler(c *gin.Context) {
	// TODO: get token from request params
	// TODO: validate token
	// TODO: activate user
	// TODO: call external service to set user details and other fields (event)
	// TODO: save user
	// TODO: issue registration confirmed event
}
