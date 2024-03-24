package core

type CreateUserDto struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
type VerifyUserDto struct {
	Code  int32  `json:"code" validate:"required,gte=100000,lte=999999"`
	Email string `json:"email" validate:"required"`
}

type CurrentUserResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}
type CreateUserResponse struct {
	Email string `json:"email"`
}
