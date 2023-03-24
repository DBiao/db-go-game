package service

type IUserService interface {
}

type userService struct {
}

func NewUserService() IUserService {
	return &userService{}
}
