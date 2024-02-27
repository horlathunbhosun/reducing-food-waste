package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/horlathunbhosun/reducing-food-waste/models"
	"github.com/horlathunbhosun/reducing-food-waste/pkg/response"
	"github.com/horlathunbhosun/reducing-food-waste/validator"
	"net/http"
	"strconv"
)

func Signup(ctx *gin.Context) {
	var user models.User
	var responseBody response.JsonResponse

	err := ctx.ShouldBindJSON(&user)

	v := validator.New()

	if models.ValidateUserData(v, &user); !v.Valid() {
		//app.failedValidationResponse(w, r, v.Errors)
		responseBody.Error = true
		responseBody.ErrorMessage = v.Errors
		responseBody.Status = false
		ctx.JSON(http.StatusBadRequest, responseBody)
		return
	}

	if err != nil {
		fmt.Println(err)

		ctx.JSON(http.StatusBadRequest, responseBody)
		return
	}

	err = user.Save()
	if err != nil {
		fmt.Println(err)
		responseBody.Error = true
		responseBody.Message = "Could not save user. Try again"
		responseBody.Status = false
		responseBody.ErrorMessage = err
		ctx.JSON(http.StatusInternalServerError, responseBody)
		return
	}
	responseBody.Error = false
	responseBody.Message = "Registration successful"
	responseBody.Status = true
	responseBody.Data = user

	ctx.JSON(http.StatusCreated, responseBody)
}

func VerificationToken(ctx *gin.Context) {
	var user models.User
	var responseBody response.JsonResponse

	tokenStr := ctx.Param("token")
	if tokenStr == "" {
		responseBody.Error = true
		responseBody.Message = "Token not provided"
		responseBody.Status = false
		ctx.JSON(http.StatusBadRequest, responseBody)
		return
	}

	token, err := strconv.Atoi(tokenStr)
	if err != nil {
		responseBody.Error = true
		responseBody.Message = "Invalid token format"
		responseBody.Status = false
		responseBody.ErrorMessage = err.Error()
		ctx.JSON(http.StatusBadRequest, responseBody)
		return
	}

	valid, err := user.VerifyToken(token)
	if err != nil || !valid {
		responseBody.Error = true
		responseBody.Message = "Could not verify token"
		responseBody.Status = false
		if err != nil {
			responseBody.ErrorMessage = err.Error()
		}
		ctx.JSON(http.StatusBadRequest, responseBody)
		return
	}

	responseBody.Error = false
	responseBody.Message = "Token verified"
	responseBody.Status = true
	ctx.JSON(http.StatusOK, responseBody)
}
