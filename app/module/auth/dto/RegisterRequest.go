package dto

type RegisterRequest struct {
	Email    string `json:"email" form:"email" validate:"required,max=45"`
	UserName string `json:"userName" form:"username" validate:"required,max=25"`
	Password string `json:"password" form:"password" validate:"required,max=25"`
	Name     string `json:"name" form:"name" validate:"required,max=45"`
	Phone    string `json:"phone" form:"phone" validate:"required,max=16"`
}
