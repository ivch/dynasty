package repository

import (
	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/models/entities"
)

type Users struct {
	db *gorm.DB
}

func NewUsers(db *gorm.DB) *Users {
	return &Users{db: db}
}

func (r *Users) CreateUser(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *Users) DeleteUser(u *entities.User) error {
	return r.db.Delete(u).Error
}

func (r *Users) GetUserByID(id uint) (*entities.User, error) {
	var u entities.User
	if err := r.db.Preload("Building").Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Users) GetUserByPhone(phone string) (*entities.User, error) {
	var u entities.User
	if err := r.db.Where("phone = ?", phone).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Users) ValidateRegCode(code string) error {
	var c struct{ Exists bool }
	if err := r.db.Raw("select exists(select id from reg_codes where code = ? and not used)", c).Scan(&code).Error; err != nil {
		return err
	}

	if !c.Exists {
		return entities.ErrInvalidRegCode
	}
	return nil
}

func (r *Users) UseRegCode(code string) error {
	return r.db.Exec("update reg_codes set used = true where code = ?", code).Error
}
