package handler

import (
	"agungstartup/helper"
	"agungstartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	//tangkap input dari user
	//map input dari user ke struct RegisterUserInput
	//struct diatas kita passing sebagai parameter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// var errors []string

		// for _, e := range err.(validator.ValidationErrors) {
		// 	errors = append(errors, e.Error())
		// }
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register Account Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register Account Failed", http.StatusBadRequest, "Failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "tokentokentoken")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context) {
	//user memasukkan input (email & password)
	//input ditangkap handler
	//mapping dari input user ke input struct
	//input struct passing ke service
	//di service mencari dengan bantuan repository user dengan email x
	//mencocokan password

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "Failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, "tokentokentoken")

	response := helper.APIResponse("Login Success", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email is not available", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		response := helper.APIResponse("Email is not available", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{
		"is_available": IsEmailAvailable,
	}

	metaMessage := "Email has been registered"
	if IsEmailAvailable {
		metaMessage = "Email is Available"
	}
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	//tangkap input dari user
	//simpan gambarnya di folder "images/"
	//service panggil repository
	//JWT (sementara hardcode user ID nya, seakann user yang login adalah user ID = 1)
	//repository ambil data user dengan ID = 1
	//update data user dengan filepath yang baru

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	//harusnya userID didapatkan dari JWT, tapi karena belum belajar JWT, untuk sementara userID di hardcode
	userID := 1

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar Successfully Uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
