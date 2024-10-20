package categories

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a category in the database
type Category struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"not null" json:"name" validate:"required"` // Nama Users
	Description string         `json:"description"`                              // Deskripsi Users
	CreatedBy   string         `gorm:"column:created_by" json:"created_by"`      // Pengguna yang membuat
	ModifiedBy  string         `gorm:"column:modified_by" json:"modified_by"`    // Pengguna yang mengubah
	DeletedBy   string         `gorm:"column:deleted_by" json:"deleted_by"`      // Pengguna yang menghapus
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`      // Waktu dibuat
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`      // Waktu diperbarui
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`      // Waktu dihapus
}
