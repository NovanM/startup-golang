package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	userJwt     auth.Service
}

func NewUserHandler(userService user.Service, userJwt auth.Service) *userHandler {
	return &userHandler{userService, userJwt}
}

func (uh *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegistrasiUser
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorsMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register user failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	newUser, err := uh.userService.RegisterUser(input)
	//token jwt
	if err != nil {
		response := helper.APIResponse("Register user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	token, err := uh.userJwt.GenerateToken(newUser.ID)

	if err != nil {
		response := helper.APIResponse("Register user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formater := user.FormaterUser(newUser, token)
	response := helper.APIResponse("Your account has been created", http.StatusOK, `success`, formater)

	c.JSON(http.StatusOK, response)

}

func (uh *userHandler) Login(c *gin.Context) {
	var input user.LoginUser
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorsMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedin, err := uh.userService.LoginUser(input)
	if err != nil {
		errorsMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	token, err := uh.userJwt.GenerateToken(loggedin.ID)
	if err != nil {
		response := helper.APIResponse("Register user failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	formater := user.FormaterUser(loggedin, token)
	response := helper.APIResponse("Succesfuly loggedin", http.StatusOK, `success`, formater)

	c.JSON(http.StatusOK, response)
}

func (uh *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailAvailable

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorsMessage := gin.H{"errors": errors}

		response := helper.APIResponse("email is not valid", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	emailCheck, err := uh.userService.CheckEmailAvailability(input)
	if err != nil {
		errorsMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	metaMessage := "Email has been registered"
	data := gin.H{
		"is_available": emailCheck,
	}
	if emailCheck {
		metaMessage = "Email is available"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, `success`, data)
	c.JSON(http.StatusOK, response)

}

func (uh *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {

		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	//token jwt user
	currentUser := c.MustGet("currentUser").(user.User)
	userId := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userId, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {

		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = uh.userService.SaveAvatar(userId, path)
	if err != nil {
		fmt.Println("error 3")
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
