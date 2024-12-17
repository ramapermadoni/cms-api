package media

import (
	"errors"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type MediaService interface {
	GetAllMediaService(ctx *gin.Context) ([]Media, int64, int, int, error)
	GetMediaByIDService(ctx *gin.Context) (Media, error)
	CreateMediaService(ctx *gin.Context) (Media, error)
	UpdateMediaService(ctx *gin.Context) (Media, error)
	DeleteMediaService(ctx *gin.Context) error
	UploadMediaService(ctx *gin.Context) (Media, error)
}
type mediaService struct {
	repo      Repository
	validator *validator.Validate
}

func NewMediaService(repo Repository) MediaService {
	return &mediaService{
		repo:      repo,
		validator: validator.New(),
	}
}

// GetAllMediaService retrieves all media with pagination and search
func (s *mediaService) GetAllMediaService(ctx *gin.Context) ([]Media, int64, int, int, error) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")

	medias, total, err := s.repo.SelectAllMedia(page, limit, search)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	// Add base URL to FilePath if it's not a valid URL
	for i := range medias {
		medias[i].FilePath = addBaseURLIfMissing(medias[i].FilePath)
	}
	return medias, total, page, limit, nil
}
func (s *mediaService) GetMediaByIDService(ctx *gin.Context) (Media, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Media{}, errors.New("invalid ID")
	}

	media, err := s.repo.SelectMediaByID(id)
	if err != nil {
		return Media{}, err
	}

	// Add base URL to FilePath if it's not a valid URL
	media.FilePath = addBaseURLIfMissing(media.FilePath)

	return media, nil
}

func (s *mediaService) CreateMediaService(ctx *gin.Context) (Media, error) {
	var media Media
	if err := ctx.ShouldBindJSON(&media); err != nil {
		return Media{}, errors.New("invalid input data")
	}

	if err := s.validator.Struct(media); err != nil {
		return Media{}, err
	}
	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Media{}, errors.New("username not found in context")
	}

	// Set CreatedBy field
	media.CreatedBy = username
	media.ModifiedBy = username

	if err := s.repo.InsertMedia(&media); err != nil {
		// return Media{}, err
		return Media{}, errors.New("failed to add new media. " + err.Error())
	}

	return media, nil
}

// UpdateMedia updates an existing media
func (s *mediaService) UpdateMediaService(ctx *gin.Context) (Media, error) {
	var media Media

	// Parse the ID from the URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return Media{}, errors.New("invalid ID format")
	}

	// Check if the media exists in the database
	existingMedia, err := s.repo.SelectMediaByID(int64(id))
	if err != nil {
		return Media{}, errors.New("media not found")
	}

	// Bind the JSON request to the media struct
	if err := ctx.ShouldBindJSON(&media); err != nil {
		return Media{}, errors.New("invalid input data")
	}

	if err := s.validator.Struct(media); err != nil {
		return Media{}, err
	}

	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Media{}, errors.New("username not found in context")
	}

	media.CreatedBy = existingMedia.CreatedBy
	media.CreatedAt = existingMedia.CreatedAt
	// Set Modified field
	media.ModifiedBy = username
	media.UpdatedAt = time.Now()
	// Assign the existing media ID to the new media object
	media.ID = existingMedia.ID

	// Update the media in the database
	if err := s.repo.UpdateMedia(media); err != nil {
		return Media{}, errors.New("failed to update media. " + err.Error())
	}

	return media, nil
}

// DeleteMedia deletes a media
func (s *mediaService) DeleteMediaService(ctx *gin.Context) error {
	// Parse the media ID from the URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Check if the media exists
	existingMedia, err := s.repo.SelectMediaByID(int64(id))
	if err != nil {
		return errors.New("media not found")
	}
	// Check if the media exists
	_, err = s.repo.SelectMediaByID(int64(id))
	if err != nil {
		return errors.New("media not found")
	}

	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return errors.New("username not found in context")
	}
	// Set DeletedBy field
	existingMedia.DeletedBy = username

	// Proceed to delete the media by updating DeletedBy and soft-deleting
	if err := s.repo.DeleteMedia(existingMedia); err != nil {
		return errors.New("failed to delete media: " + err.Error())
	}

	return nil
}
func (s *mediaService) UploadMediaService(ctx *gin.Context) (Media, error) {
	// Ambil file dari request
	file, err := ctx.FormFile("file")
	if err != nil {
		return Media{}, errors.New("failed to retrieve file: " + err.Error())
	}
	// Dapatkan ekstensi file (contoh: .jpg, .png, .mp4, .mp3)
	extension := strings.ToLower(filepath.Ext(file.Filename))

	// Cek tipe MIME file
	mimeType := file.Header.Get("Content-Type")
	if !isValidMediaType(mimeType, extension) {
		return Media{}, errors.New("invalid file type: must be image, video, or audio")
	}

	// Generate filename unik dengan UUID
	uniqueFilename := uuid.New().String() + extension

	// Tentukan path penyimpanan file
	destination := "./uploads/" + uniqueFilename

	// Simpan file ke server
	if err := ctx.SaveUploadedFile(file, destination); err != nil {
		return Media{}, errors.New("failed to save file: " + err.Error())
	}

	// Ambil username dari context
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Media{}, errors.New("username not found in context")
	}

	// Generate data media untuk disimpan ke database
	media := Media{
		FileName:   uniqueFilename,
		FilePath:   destination,
		CreatedBy:  username,
		ModifiedBy: username,
	}

	// Simpan media ke database
	if err := s.repo.UploadMedia(&media); err != nil {
		return Media{}, errors.New("failed to save media data: " + err.Error())
	}

	return media, nil
}

// addBaseURLIfMissing checks if the URL is valid, if not it adds the base URL
func addBaseURLIfMissing(filePath string) string {
	baseURL := viper.GetString("base_url") // Ganti dengan URL produksi
	// Check if filePath is a valid URL
	_, err := url.ParseRequestURI(filePath)
	if err != nil || (filePath[:4] != "http" && filePath[:5] != "https") {
		// Remove leading "./" if present
		filePath = strings.TrimPrefix(filePath, "./")
		return baseURL + filePath // Add base URL
	}
	return filePath // Return the original file path if it's valid
}

// isValidMediaType checks if the uploaded file type is valid
func isValidMediaType(mimeType string, extension string) bool {
	validImageTypes := []string{"image/jpeg", "image/png", "image/gif"}
	validVideoTypes := []string{"video/mp4", "video/mpeg", "video/quicktime"}
	validAudioTypes := []string{"audio/mpeg", "audio/wav", "audio/ogg"}

	// Check MIME type
	for _, v := range validImageTypes {
		if mimeType == v {
			return true
		}
	}
	for _, v := range validVideoTypes {
		if mimeType == v {
			return true
		}
	}
	for _, v := range validAudioTypes {
		if mimeType == v {
			return true
		}
	}

	// Optionally check extension if MIME type is unknown
	if extension == ".jpg" || extension == ".jpeg" ||
		extension == ".png" || extension == ".gif" ||
		extension == ".mp4" || extension == ".mpeg" ||
		extension == ".mp3" || extension == ".wav" ||
		extension == ".ogg" {
		return true
	}

	return false
}
