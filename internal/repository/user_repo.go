package repository

import (
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
			return nil, errors.New("usuário não encontrado")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("senha incorreta")
	}

	return &user, nil
}

func (r *UserRepo) GetAllUsersId() ([]uint, error) {
	var usersId []uint
	if err := r.DB.Find(&usersId).Error; err != nil {
		return nil, err
	}
	return usersId, nil
}
