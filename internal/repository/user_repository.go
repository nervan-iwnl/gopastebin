package repository

import (
	"gopastebin/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *domain.User) error
	Update(u *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmailOrUsername(login string) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	ExistsByEmail(email string) (bool, error)
	ExistsByUsername(username string) (bool, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *userRepo) Update(u *domain.User) error {
	return r.db.Save(u).Error
}

func (r *userRepo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByEmailOrUsername(login string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ? OR username = ?", login, login).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) FindByUsername(username string) (*domain.User, error) {
	var u domain.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) ExistsByEmail(email string) (bool, error) {
	var c int64
	if err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&c).Error; err != nil {
		return false, err
	}
	return c > 0, nil
}

func (r *userRepo) ExistsByUsername(username string) (bool, error) {
	var c int64
	if err := r.db.Model(&domain.User{}).Where("username = ?", username).Count(&c).Error; err != nil {
		return false, err
	}
	return c > 0, nil
}
