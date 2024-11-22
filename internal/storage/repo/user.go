package server

import (
	"errors"
	"tender-backend/model"
	request_model "tender-backend/model/request"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) CreateUser(user *request_model.CreateUserReq) (*model.User, error) {
	newUser := model.User{
		FullName: user.FullName,
		Password: user.Password,
		Email:    user.Email,
		Role:     user.Role,
		Username: user.Username,
	}

	if _, err := s.GetByUsername(user.Email); err == nil {
		return nil, errors.New("Email already exists")
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (s *UserService) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetByUsername(email string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(user *request_model.UpdateUserReq, id int64) (*model.User, error) {
	var existingUser model.User
	if err := s.db.First(&existingUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	existingUser.FullName = user.FullName
	existingUser.Email = user.Email

	if err := s.db.Save(&existingUser).Error; err != nil {
		return nil, err
	}

	return &existingUser, nil
}

func (s *UserService) DeleteUser(id int64) error {
	if err := s.db.Delete(&model.User{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}
