package users

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the database
type User struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`                               // ID untuk pengguna
	Fullname   string         `gorm:"not null" json:"fullname" validate:"required,min=6"`               // Nama lengkap pengguna
	Username   string         `gorm:"unique;not null" json:"username" validate:"required,min=3,max=50"` // Username unik
	Password   string         `gorm:"not null" json:"password" validate:"required,min=6"`               // Kata sandi terenkripsi
	Email      string         `gorm:"unique;not null" json:"email" validate:"email"`                    // Email unik
	Role       string         `gorm:"not null" json:"role" validate:"required"`                         // Peran pengguna (admin/editor/author)
	CreatedBy  string         `gorm:"column:created_by" json:"created_by"`                              // Pengguna yang membuat
	ModifiedBy string         `gorm:"column:modified_by" json:"modified_by"`                            // Pengguna yang mengubah
	DeletedBy  string         `gorm:"column:deleted_by" json:"deleted_by"`                              // Pengguna yang menghapus
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`                              // Waktu dibuat
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`                              // Waktu diperbarui
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`                              // Waktu dihapus
}
