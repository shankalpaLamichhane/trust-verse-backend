package service

import (
	"trust-verse-backend/app/module/auth/dto"
)

type AuthService struct {
	//userService *service.UserService
}

type IAuthService interface {
	Login(request dto.LoginRequest) (string, error)
	Register(request dto.RegisterRequest) (string, error)
}

func NewAuthService(
// userService *service.UserService
) *AuthService {
	return &AuthService{
		//userService: userService,
	}
}

func (s *AuthService) Login(request dto.LoginRequest) (string, error) {
	//user, err := s.userService.GetUser()
	//if err != nil {
	//	return "", err
	//}
	//fmt.Print(user)
	return "token", nil
}

func (s *AuthService) Register(request dto.RegisterRequest) (string, error) {
	//user, err := s.userService.GetUser()
	//if err != nil {
	//	return "", err
	//}
	//fmt.Print(user)
	return "token", nil
}
