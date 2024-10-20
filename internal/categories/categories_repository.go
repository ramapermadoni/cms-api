package categories

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository interface for Category operations
type Repository interface {
	SelectAllCategory(page, limit int, search, orderBy, sort string) ([]Category, int64, error)
	SelectCategoryByID(id int64) (Category, error)
	InsertCategory(category *Category) error
	UpdateCategory(category Category) error
	DeleteCategory(category Category) error
}

// Repository struct implementing Repository interface
type categoryRepository struct {
	db *gorm.DB
}

// NewRepository initializes the repository with a DB connection
func NewRepository(db *gorm.DB) Repository {
	return &categoryRepository{db: db}
}

// SelectAllCategory retrieves all category from the database
func (r *categoryRepository) SelectAllCategory(page, limit int, search, orderBy, sort string) ([]Category, int64, error) {
	var category []Category
	var total int64
	query := r.db.Model(&Category{})
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
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
	if err := query.Order(orderClause).Offset(offset).Limit(limit).Find(&category).Error; err != nil {
		return nil, 0, err
	}
	return category, total, nil
}

func (r *categoryRepository) SelectCategoryByID(id int64) (Category, error) {
	var category Category
	err := r.db.First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Category{}, errors.New("category not found")
	}
	return category, err
}

func (r *categoryRepository) InsertCategory(category *Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) UpdateCategory(category Category) error {
	return r.db.Debug().Model(&Category{}).Where("id = ?", category.ID).Updates(category).Error
}

func (r *categoryRepository) DeleteCategory(category Category) error {
	return r.db.Model(&category).Update("deleted_by", category.DeletedBy).Delete(&category).Error
}
