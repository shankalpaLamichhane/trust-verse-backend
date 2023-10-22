package dto

type LoginRequest struct {
	UserName string `json:"userName" form:"username" validate:"required,max=25"`
	Password string `json:"password" form:"password" validate:"required,max=25"`
}
