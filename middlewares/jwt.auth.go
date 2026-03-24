package middlewares

import (
	env "AuthInGo/config/env"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader:=r.Header.Get("Authorization")

		if authHeader==""{
			http.Error(w,"Authorization header required",http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(authHeader,"Bearer "){
			http.Error(w,"Authorization header must start with Bearer",http.StatusUnauthorized)
			return
		}

		token:=strings.TrimPrefix(authHeader,"Bearer ")
		if token==""{
			http.Error(w,"Token required",http.StatusUnauthorized)
			return
		}

		claims:=jwt.MapClaims{}

		_,err:=jwt.ParseWithClaims(token,claims,func(t *jwt.Token) (any, error) {
			return []byte(env.GetString("JWT_SECRET","SECRET")),nil
		})

		if err!=nil{
			http.Error(w,"Invalid Token",http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		userId,okId:=claims["id"].(float64)
		userEmail,okEmail:=claims["email"].(string)

		if !okId || !okEmail{
			http.Error(w,"Invalid token claims",http.StatusUnauthorized)
			return
		}

		fmt.Println("UserId:",userId,"Email:",userEmail)

		next.ServeHTTP(w,r)
	})
}