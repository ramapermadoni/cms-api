package posts

import (
	"time"

	"gorm.io/gorm"
)

// Post represents a post in the database
type Post struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`                                                                            // ID untuk post
	Title      string         `gorm:"not null" json:"title" validate:"required"`                                                                     // Judul post
	Content    string         `gorm:"not null" json:"content" validate:"required"`                                                                   // Konten post
	CategoryID uint           `gorm:"not null;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CategoryID;references:ID" json:"category_id"` // Foreign key ke categories
	Author     string         `gorm:"not null" json:"author" validate:"required"`                                                                    // Penulis post
	Status     string         `gorm:"default:'draft'" json:"status"`                                                                                 // Status post (draft/published)
	CreatedBy  string         `gorm:"column:created_by" json:"created_by"`                                                                           // Pengguna yang membuat
	ModifiedBy string         `gorm:"column:modified_by" json:"modified_by"`                                                                         // Pengguna yang mengubah
	DeletedBy  string         `gorm:"column:deleted_by" json:"deleted_by"`                                                                           // Pengguna yang menghapus
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`                                                                           // Waktu dibuat
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`                                                                           // Waktu diperbarui
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`                                                                           // Waktu dihapus
}
