package posts

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// Repository interface for Post operations
type Repository interface {
	SelectAllPost(page, limit int, search, orderBy, sort string) ([]Post, int64, error)
	SelectPostByID(id int64) (Post, error)
	InsertPost(post *Post) error
	UpdatePost(post Post) error
	DeletePost(post Post) error
	CheckCategoryExists(categoryID uint) (bool, error)
}

// Repository struct implementing Repository interface
type postRepository struct {
	db *gorm.DB
}

// NewRepository initializes the repository with a DB connection
func NewRepository(db *gorm.DB) Repository {
	return &postRepository{db: db}
}

// SelectAllPost retrieves all post from the database
func (r *postRepository) SelectAllPost(page, limit int, search, orderBy, sort string) ([]Post, int64, error) {
	var post []Post
	var total int64
	query := r.db.Model(&Post{})
	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
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
	if err := query.Debug().Order(orderClause).Offset(offset).Limit(limit).Find(&post).Error; err != nil {
		return nil, 0, err
	}
	return post, total, nil
}

func (r *postRepository) SelectPostByID(id int64) (Post, error) {
	var post Post
	err := r.db.First(&post, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Post{}, errors.New("post not found")
	}
	return post, err
}

func (r *postRepository) InsertPost(post *Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) UpdatePost(post Post) error {
	return r.db.Save(&post).Error
}

func (r *postRepository) DeletePost(post Post) error {
	return r.db.Model(&post).Update("deleted_by", post.DeletedBy).Delete(&post).Error
}

// Check if CategoryID exists in the categories table
func (r *postRepository) CheckCategoryExists(categoryID uint) (bool, error) {
	var count int64
	err := r.db.Table("categories").Where("id = ?", categoryID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
