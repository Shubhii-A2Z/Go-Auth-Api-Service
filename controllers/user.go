package controllers

import (
	"AuthInGo/dtos"
	"AuthInGo/services"
	"AuthInGo/utils"
	"fmt"
	"net/http"
)

// // UserController depending on UserService Interface: (Controllerlayer->ServiceLayer)
type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController{
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterUser Controller Called")
	uc.UserService.CreateUser()
}

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUserById Controller Called")
	uc.UserService.GetUserById()
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers Controller Called")
	uc.UserService.GetAllUsers()
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginUser Controller Called")

	var payload dtos.LoginUserRequestDTO

	if jsonErr:=utils.ReadJsonBody(r,&payload); jsonErr!=nil {
		utils.WriteJsonErrorResponse(w,http.StatusBadRequest,"Invalid input data")
		return
	}

	jwtToken,err:=uc.UserService.LoginUser(&payload)

	if err!=nil {
		utils.WriteJsonErrorResponse(w,http.StatusInternalServerError,"Failed to login user")
		return
	}

	utils.WriteJsonSuccessResponse(w,http.StatusOK,"User loggedIn successfully",jwtToken)
}