package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/horlathunbhosun/reducing-food-waste/models"
	"github.com/horlathunbhosun/reducing-food-waste/pkg/response"
	"github.com/horlathunbhosun/reducing-food-waste/validator"
	"net/http"
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
		responseBody.Error = true
		responseBody.Message = "something went wrong"
		responseBody.Status = false
		ctx.JSON(http.StatusBadRequest, responseBody)
		return
	}

	//err = user.Save()
	//if err != nil {
	//	fmt.Println(err)
	//	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user. Try again"})
	//	return
	//}
	//ctx.JSON(http.StatusCreated, gin.H{"message": "Registration successful", "data": user})
}
