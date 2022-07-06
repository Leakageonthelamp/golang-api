package service

import (
	"log"

	"github.com/Leakageonthelamp/golang-api/dto"
	"github.com/Leakageonthelamp/golang-api/entity"
	"github.com/Leakageonthelamp/golang-api/repository"
	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user dto.UserUpdateDTO) entity.User
	Profile(UserID string) entity.User
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepo,
	}
}

func (service *userService) Update(user dto.UserUpdateDTO) entity.User {
	userToUpdate := entity.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Error mapping fields: %v", err)
	}
	return service.userRepository.UpdateUser(userToUpdate)
}

func (service *userService) Profile(UserID string) entity.User {
	return service.userRepository.ProfileUser(UserID)
}
