package core

import "github.com/dgrijalva/jwt-go"

type CreateAuthCodeDto struct {
	UserId string
	Email  string
	Name   string
}

type AuthTokenClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}

type SignInDto struct {
	Code  int32  `json:"code" validate:"required,gte=100000,lte=999999"`
	Email string `json:"email" validate:"required,email"`
}
type GenerateAuthCodeDto struct {
	Email string `json:"email" validate:"required,email"`
}
type GenerateAuthCodeResponse struct {
	Email string `json:"email" validate:"required,email"`
}
