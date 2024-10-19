package db

import (
	"fmt"
	"gopastebin/src/models"
)

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.Where("email =?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(nickname string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", nickname).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *models.User) error {
	if user == nil {
		return fmt.Errorf("received nil user")
	}

	fmt.Println("Creating user:", user.Email)
	return db.Create(user).Error
}

func UpdateUser(user *models.User) error {
	return db.Save(user).Error
}

func IsEmailUnique(email string) (bool, error) {
	var count int64
	err := db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func IsUsernameUnique(email string) (bool, error) {
	var count int64
	err := db.Model(&models.User{}).Where("username = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
