package repo

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/server/handlers/users"
)

type Users struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Users {
	db.AutoMigrate(&users.User{})
	return &Users{db: db}
}

func (r *Users) CreateUser(user *users.User) error {
	return r.db.Create(user).Error
}

func (r *Users) DeleteUser(u *users.User) error {
	return r.db.Delete(u).Error
}

func (r *Users) UpdateUser(u *users.User) error {
	return r.db.Save(u).Error
}

func (r *Users) GetUserByID(id uint) (*users.User, error) {
	var u users.User
	if err := r.db.Preload("Building").Preload("Entry").
		Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Users) GetUserByPhone(phone string) (*users.User, error) {
	var u users.User
	if err := r.db.Where("phone = ?", phone).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
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
		return errs.RegCodeInvalid
	}

	return nil
}

func (r *Users) UseRegCode(code string) error {
	return r.db.Exec("update reg_codes set used = true where code = ?", code).Error
}

func (r *Users) GetRegCode() (string, error) {
	var code []string
	if err := r.db.Table("reg_codes").Where("used = ?", false).Pluck("code", &code).Error; err != nil {
		return "", err
	}

	if len(code) == 0 {
		return "", errors.New("no reg code available")
	}

	return code[0], nil
}

func (r *Users) GetFamilyMembers(ownerID uint) ([]*users.User, error) {
	var res []*users.User
	if err := r.db.Where("parent_id = ?", ownerID).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Users) FindUserByApartment(building uint, apt uint) (*users.User, error) {
	var u users.User
	if err := r.db.Where("building_id = ? AND apartment = ? AND parent_id IS NULL", building, apt).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
