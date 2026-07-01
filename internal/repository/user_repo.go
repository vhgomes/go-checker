package repository

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Name     string
	Password string
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(email, password, name string) error {
	var existing User
	if err := r.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return errors.New("email já cadastrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("falha ao gerar hash da senha")
	}

	user := User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}

	return r.DB.Create(&user).Error
}

func (r *UserRepo) Login(email, password string) (*User, error) {
	var user User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("credenciais inválidas ")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("credenciais inválidas")
	}

	return &user, nil
}

func (r *UserRepo) GetAllUsersId(ctx context.Context) ([]uint, error) {
	var ids []uint
	if err := r.DB.WithContext(ctx).
		Model(&User{}).
		Pluck("id", &ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}
