package repo

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/ivch/dynasty/common/errs"
	"github.com/ivch/dynasty/server/handlers/users"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	db.AutoMigrate(&users.User{})
	return &Repo{db: db}
}

func (r *Repo) CreateUser(user *users.User) error {
	return r.db.Create(user).Error
}

func (r *Repo) DeleteUser(u *users.User) error {
	return r.db.Delete(u).Error
}

func (r *Repo) UpdateUser(req *users.UserUpdate) error {
	update := prepareUpdateQuery(req)
	return r.db.Table(users.User{}.TableName()).Where("id = ?", req.ID).Updates(update).Error
}

func (r *Repo) GetUserByID(id uint) (*users.User, error) {
	var u users.User
	if err := r.db.Preload("Building").Preload("Entry").
		Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repo) GetUserByPhone(phone string) (*users.User, error) {
	var u users.User
	if err := r.db.Where("phone = ?", phone).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repo) ValidateRegCode(code string) error {
	var c struct{ Exists bool }
	if err := r.db.Raw("select exists(select id from reg_codes where code = ? and not used)", c).Scan(&code).Error; err != nil {
		return err
	}

	if !c.Exists {
		return errs.RegCodeInvalid
	}

	return nil
}

func (r *Repo) UseRegCode(code string) error {
	return r.db.Exec("update reg_codes set used = true where code = ?", code).Error
}

func (r *Repo) GetRegCode() (string, error) {
	var code []string
	if err := r.db.Table("reg_codes").Where("used = ?", false).Pluck("code", &code).Error; err != nil {
		return "", err
	}

	if len(code) == 0 {
		return "", errors.New("no reg code available")
	}

	return code[0], nil
}

func (r *Repo) GetFamilyMembers(ownerID uint) ([]*users.User, error) {
	var res []*users.User
	if err := r.db.Where("parent_id = ?", ownerID).Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *Repo) FindUserByApartment(building uint, apt uint) (*users.User, error) {
	var u users.User
	if err := r.db.Where("building_id = ? AND apartment = ? AND parent_id IS NULL", building, apt).First(&u).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}

func (r *Repo) CreateRecoverCode(c *users.PasswordRecovery) error {
	return r.db.Create(c).Error
}

func (r *Repo) GetRecoveryCode(c *users.PasswordRecovery) (*users.PasswordRecovery, error) {
	var code users.PasswordRecovery
	if err := r.db.Where("code = ? AND active = ? ", c.Code, c.Active).
		First(&code).Error; err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *Repo) CountRecoveryCodesByUserIn24h(userID uint) (int, error) {
	var count int
	from := time.Now().Add(-24 * time.Hour)
	to := time.Now()
	if err := r.db.Model(&users.PasswordRecovery{}).
		Where("created_at >= ? AND created_at <= ?", from, to).
		Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repo) ResetPassword(codeID uint, req *users.UserUpdate) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		update := prepareUpdateQuery(req)
		if err := tx.Table(users.User{}.TableName()).Where("id = ?", req.ID).Updates(update).Error; err != nil {
			return err
		}
		return tx.Model(users.PasswordRecovery{}).Where("id = ?", codeID).Update("active", "false").Error
	})
}

func prepareUpdateQuery(req *users.UserUpdate) map[string]interface{} {
	update := make(map[string]interface{})
	if req.Email != nil {
		update["email"] = *req.Email
	}

	if req.Password != nil {
		update["password"] = *req.Password
	}

	if req.FirstName != nil {
		update["first_name"] = *req.FirstName
	}

	if req.LastName != nil {
		update["last_name"] = *req.LastName
	}

	if req.Active != nil {
		update["active"] = *req.Active
	}

	return update
}
