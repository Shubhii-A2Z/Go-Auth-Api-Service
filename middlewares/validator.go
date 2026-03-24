package middlewares

import (
	"AuthInGo/dtos"
	"AuthInGo/utils"
	"context"
	"net/http"
)

func UserLoginRequestValidator(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dtos.LoginUserRequestDTO
		
		// Read and decode JSON body
		if err:=utils.ReadJsonBody(r,&payload); err!=nil{
			utils.WriteJsonErrorResponse(w,http.StatusBadRequest,"Invalid request body")
			return
		}

		// Original context coming for current request
		reqContext:=r.Context()

		// Create a new context with payload
		ctx:=context.WithValue(reqContext,"payload",payload)

		// Call next handler/controller in chain
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

func UserCreateRequestValidator(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload dtos.CreateUserRequestDTO
		
		// Read and decode JSON body
		if err:=utils.ReadJsonBody(r,&payload); err!=nil{
			utils.WriteJsonErrorResponse(w,http.StatusBadRequest,"Invalid request body")
			return
		}

		// Original context coming for current request
		reqContext:=r.Context()

		// Create a new context with payload
		ctx:=context.WithValue(reqContext,"payload",payload)

		// Call next handler/controller in chain
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}