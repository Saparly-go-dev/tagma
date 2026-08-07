package service

import (
	"errors"
	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(user tagma.CreateUser) error {
	if !isAllowedRole(user.Type) {
		return errors.New("invalid user role")
	}
	if user.Password != user.ConfirmPassword {
		return errors.New("password error")
	}

	password, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = password

	return s.repo.Create(user)
}

func (s *UserService) GetPage(pageSize, pageNumber int, name, tip, language string) (tagma.UserPage, error) {
	return s.repo.GetPage(pageSize, pageNumber, name, tip, language)
}

func (s *UserService) ChangePassword(data tagma.Change_Password) error {

	if data.Password != data.ConfirmPassword {
		return errors.New("password error")
	}

	password, err := hashPassword(data.Password)
	if err != nil {
		return err
	}
	data.Password = password

	return s.repo.ChangePassword(data)
}

func (s *UserService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *UserService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *UserService) Update(Id int, data tagma.CreateUser) error {
	return s.repo.Update(Id, data)
}
