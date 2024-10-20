package users

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository interface for User operations
type Repository interface {
	SelectAllUser(page, limit int, search, orderBy, sort string) ([]User, int64, error)
	SelectUserByID(id int64) (User, error)
	InsertUser(user *User) error
	UpdateUser(user User) error
	DeleteUser(user User) error
}

// Repository struct implementing Repository interface
type userRepository struct {
	db *gorm.DB
}

// NewRepository initializes the repository with a DB connection
func NewRepository(db *gorm.DB) Repository {
	return &userRepository{db: db}
}

// SelectAllUser retrieves all user from the database
func (r *userRepository) SelectAllUser(page, limit int, search, orderBy, sort string) ([]User, int64, error) {
	var user []User
	var total int64
	query := r.db.Model(&User{})
	if search != "" {
		query = query.Where("username ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if orderBy == "" {
		orderBy = "updated_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}
	// Order secara dinamis
	orderClause := fmt.Sprintf("%s %s", orderBy, sort)

	offset := (page - 1) * limit
	if err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&user).Error; err != nil {
		return nil, 0, err
	}
	return user, total, nil
}

func (r *userRepository) SelectUserByID(id int64) (User, error) {
	var user User
	err := r.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errors.New("user not found")
	}
	return user, err
}

func (r *userRepository) InsertUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user User) error {
	return r.db.Debug().Model(&User{}).Where("id = ?", user.ID).Updates(user).Error
}

// func (r *userRepository) DeleteUser(user User) error {
// 	return r.db.Delete(&user).Error
// }

func (r *userRepository) DeleteUser(user User) error {
	return r.db.Model(&user).Update("deleted_by", user.DeletedBy).Delete(&user).Error
}
