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

	payload:= r.Context().Value("payload").(dtos.CreateUserRequestDTO)

	err:=uc.UserService.CreateUser(&payload)

	if err!=nil{
		utils.WriteJsonErrorResponse(w,http.StatusInternalServerError,"Failed to create user")
		return
	}

	utils.WriteJsonSuccessResponse(w,http.StatusCreated,"User successfully created",payload)
}

func (uc *UserController) GetUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUserById Controller Called")
	// extract userId from url params
	userId:=r.URL.Query().Get("id")

	if userId==""{
		utils.WriteJsonErrorResponse(w,http.StatusBadRequest,"User id required")
	}

	user,err:=uc.UserService.GetUserById(userId)

	if err!=nil{
		utils.WriteJsonErrorResponse(w,http.StatusInternalServerError,"Failed to fetch user")
	}

	if user==nil{
		utils.WriteJsonErrorResponse(w,http.StatusNotFound,"User not found")
	}

	fmt.Println("User found:",user)
}

func (uc *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetUsers Controller Called")
	uc.UserService.GetAllUsers()
}

func (uc *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LoginUser Controller Called")

	payload:=r.Context().Value("payload").(dtos.LoginUserRequestDTO)

	jwtToken,err:=uc.UserService.LoginUser(&payload)

	if err!=nil {
		utils.WriteJsonErrorResponse(w,http.StatusInternalServerError,"Failed to login user")
		return
	}

	utils.WriteJsonSuccessResponse(w,http.StatusOK,"User loggedIn successfully",jwtToken)
}