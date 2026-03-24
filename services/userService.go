package services

import (
	env "AuthInGo/config/env"
	db "AuthInGo/db/repositories"
	"AuthInGo/dtos"
	"AuthInGo/utils"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	CreateUser(*dtos.CreateUserRequestDTO) error
	GetUserById() error
	GetAllUsers() error
	LoginUser(*dtos.LoginUserRequestDTO) (string,error)
}

// UserService depending on UserRepo Interface: (ServiceLayer->RepositoryLayer)
type UserServiceImpl struct {
	userRepository db.UserRepository
}

func NewUserService(userRepo db.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: userRepo,
	}
}

func (u *UserServiceImpl) CreateUser(payload *dtos.CreateUserRequestDTO) error {
	// Hashing the password using utils.HashPassword
	hashPass,err:=utils.HashPassword(payload.Password)
	if err!=nil {
		fmt.Println("Error hashing password:",err)
		return err
	}

	// Calling the repository to create user
	err2:=u.userRepository.Create(payload.Username,payload.Email,hashPass)
	if err2!=nil {
		fmt.Println("Error creating user:",err)
		return err2
	}
	return nil
}

func (u *UserServiceImpl) GetUserById() error {
	u.userRepository.GetById()
	return nil
}

func (u *UserServiceImpl) GetAllUsers() error {
	u.userRepository.GetAll()
	return nil
}

func (u *UserServiceImpl) LoginUser(payload *dtos.LoginUserRequestDTO) (string,error) {
	email:=payload.Email
	password:=payload.Password

	user,err:=u.userRepository.GetByEmail(email)

	if err!=nil{
		fmt.Println("Error fetching user:",err)
		return "",err
	}

	// Check if user exists or not
	if user==nil{
		fmt.Println("User not found with given email")
		return "",fmt.Errorf("No user with email %s",email)
	}

	// If user exists, check is password is correct
	isPasswordValid:=utils.CheckPasswordHash(password,user.Password)
	if !isPasswordValid{
		fmt.Println("Incorrect Password")
		return "",nil
	}

	// Creating payload object
	jwtPayload:=jwt.MapClaims{
		"email": user.Email,
		"id": user.Id,
	}

	// Creating JWT token object
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwtPayload)

	// Converting to JWT token string, signed using secret key
	tokenString,err:=token.SignedString([]byte(env.GetString("JWT_SECRET","SECRET")))

	if err!=nil{
		fmt.Println("Error signing token:",err)
		return "",err
	}

	fmt.Println("JWT Token:",tokenString)
	return tokenString,nil
}
