package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(userInput RegistrasiUser) (User, error)
	LoginUser(userInput LoginUser) (User, error)
	CheckEmailAvailability(userInput CheckEmailAvailable) (bool, error)
	SaveAvatar(id int, fileLocation string) (User, error)
	GetUserByID(id int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) *service {
	return &service{repo}
}

func (s *service) RegisterUser(userInput RegistrasiUser) (User, error) {
	user := User{
		Name:       userInput.Name,
		Email:      userInput.Email,
		Occupation: userInput.Occupation,
		Role:       "user",
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	userNew, err := s.repository.Save(user)
	if err != nil {

		return user, err
	}
	return userNew, nil

}

func (s *service) LoginUser(userInput LoginUser) (User, error) {
	email := userInput.Email
	password := userInput.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}

	return user, nil

}

func (s *service) CheckEmailAvailability(userInput CheckEmailAvailable) (bool, error) {
	email := userInput.Email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	if user.ID == 0 {
		return true, nil
	}
	return false, err
}

func (s *service) SaveAvatar(id int, fileLocation string) (User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	user.AvatarFileName = fileLocation
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return user, err
	}
	return updatedUser, nil
}

func (s *service) GetUserByID(id int) (User, error) {
	user, err := s.repository.FindByID(id)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("No user found on with that ID")
	}

	return user, nil
}
