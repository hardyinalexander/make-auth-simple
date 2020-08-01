package authentication

import (
	"context"

	"github.com/jinzhu/gorm"
)

type Repository interface {
	GetUserIDByEmail(ctx context.Context, email string) (string, error)
	CreateUser(ctx context.Context, c *User) error
	UpdateProfile(ctx context.Context, id string, updateMap map[string]interface{}) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func InitRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	user := new(User)
	err := r.db.Select("id").Where("email = ?", email).First(&user).Error
	if user == nil || gorm.IsRecordNotFoundError(err) {
		return "", nil
	}

	return user.ID, nil
}

func (r *repository) CreateUser(ctx context.Context, user *User) (err error) {
	err = r.db.Create(&user).Error
	return
}

func (r *repository) UpdateProfile(ctx context.Context, id string, updateMap map[string]interface{}) (user *User, err error) {
	err = r.db.Model(&user).Where("id = ?", id).Updates(updateMap).Error
	return
}
