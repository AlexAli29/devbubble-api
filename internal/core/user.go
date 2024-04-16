package core

type CreateUserDto struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}
type VerifyUserDto struct {
	Code  int32  `json:"code" validate:"required,gte=100000,lte=999999"`
	Email string `json:"email" validate:"required"`
}

type UpdateUserDto struct {
	Id          string `json:"id" validate:"required,uuid"`
	Name        string `json:"name" validate:"required,min=2"`
	Description string `json:"description"`
}

type CurrentUserResponse struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags"`
}

type CreateUserResponse struct {
	Email string `json:"email"`
}
