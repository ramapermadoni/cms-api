package media

import (
	"errors"

	"gorm.io/gorm"
)

// Repository interface for Media operations
type Repository interface {
	SelectAllMedia(page, limit int, search string) ([]Media, int64, error)
	SelectMediaByID(id int64) (Media, error)
	InsertMedia(media *Media) error
	UpdateMedia(media Media) error
	DeleteMedia(media Media) error
	UploadMedia(media *Media) error
	CheckMediaExists(mediaID []int) ([]int, error)
	UpdateMediaPostID(mediaIDs []int, postID int64) error
}

// Repository struct implementing Repository interface
type mediaRepository struct {
	db *gorm.DB
}

// NewRepository initializes the repository with a DB connection
func NewRepository(db *gorm.DB) Repository {
	return &mediaRepository{db: db}
}

// SelectAllMedia retrieves all media from the database
func (r *mediaRepository) SelectAllMedia(page, limit int, search string) ([]Media, int64, error) {
	var media []Media
	var total int64
	query := r.db.Model(&Media{})
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&media).Error; err != nil {
		return nil, 0, err
	}
	return media, total, nil
}

func (r *mediaRepository) SelectMediaByID(id int64) (Media, error) {
	var media Media
	err := r.db.First(&media, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Media{}, errors.New("media not found")
	}
	return media, err
}

func (r *mediaRepository) InsertMedia(media *Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) UpdateMedia(media Media) error {
	return r.db.Debug().Model(&Media{}).Where("id = ?", media.ID).Updates(media).Error
}

func (r *mediaRepository) DeleteMedia(media Media) error {
	return r.db.Model(&media).Update("deleted_by", media.DeletedBy).Delete(&media).Error
}

// UploadMedia menyimpan metadata media setelah upload berhasil.
func (r *mediaRepository) UploadMedia(media *Media) error {
	return r.db.Create(media).Error
}

func (r *mediaRepository) CheckMediaExists(mediaIDs []int) ([]int, error) {
	var validMedia []int

	// Query untuk mengecek keberadaan media berdasarkan ID
	if err := r.db.
		Table("media").
		Where("id IN ?", mediaIDs).
		Pluck("id", &validMedia).Error; err != nil {
		return nil, err
	}

	return validMedia, nil
}
func (r *mediaRepository) UpdateMediaPostID(mediaIDs []int, postID int64) error {
	// Use GORM to update the PostID for the specified media IDs
	if err := r.db.Model(&Media{}).Where("id IN ?", mediaIDs).Update("post_id", postID).Error; err != nil {
		return err
	}
	return nil
}
