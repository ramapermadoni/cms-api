package media

import (
	"time"

	"gorm.io/gorm"
)

// Media represents a media file in the database
type Media struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`    // ID untuk media
	FileName   string         `gorm:"not null" json:"file_name"`             // Nama file
	FilePath   string         `gorm:"not null" json:"file_path"`             // Path penyimpanan file
	PostID     uint           `json:"post_id"`                               // Foreign key ke posts
	CreatedBy  string         `gorm:"column:created_by" json:"created_by"`   // Pengguna yang membuat
	ModifiedBy string         `gorm:"column:modified_by" json:"modified_by"` // Pengguna yang mengubah
	DeletedBy  string         `gorm:"column:deleted_by" json:"deleted_by"`   // Pengguna yang menghapus
	CreatedAt  time.Time      `gorm:"column:created_at" json:"created_at"`   // Waktu dibuat
	UpdatedAt  time.Time      `gorm:"column:updated_at" json:"updated_at"`   // Waktu diperbarui
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`   // Waktu dihapus
}
