package posts

import (
	"cms-api/internal/media"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostService interface {
	GetAllPostService(ctx *gin.Context) ([]Post, int64, error)
	GetPostByIDService(ctx *gin.Context) (Post, error)
	CreatePostService(ctx *gin.Context) (Post, error)
	UpdatePostService(ctx *gin.Context) (Post, error)
	DeletePostService(ctx *gin.Context) error
}
type postService struct {
	repo      Repository
	validator *validator.Validate
	mediaRepo media.Repository // Inject media repository
}
type PostInput struct {
	Post  Post  `json:"post"`
	Media []int `json:"media"` // Media opsional
}

func NewPostService(repo Repository, mediaRepo media.Repository) PostService {
	return &postService{
		repo:      repo,
		validator: validator.New(),
		mediaRepo: mediaRepo,
	}
}

// GetAllPostService retrieves all post with pagination and search
func (s *postService) GetAllPostService(ctx *gin.Context) ([]Post, int64, error) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	orderBy := ctx.DefaultQuery("orderBy", "updated_at") // Default ke 'updated_at'
	sort := ctx.DefaultQuery("sort", "desc")

	posts, total, err := s.repo.SelectAllPost(page, limit, search, orderBy, sort)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
func (s *postService) GetPostByIDService(ctx *gin.Context) (Post, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Post{}, errors.New("invalid ID")
	}

	return s.repo.SelectPostByID(id)
}
func (s *postService) CreatePostService(ctx *gin.Context) (Post, error) {
	var input PostInput // Menggunakan struct baru

	// Bind JSON ke struct input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return Post{}, errors.New("invalid input data")
	}

	// Ambil post dari input
	post := input.Post

	if err := s.validator.Struct(post); err != nil {
		return Post{}, err
	}

	// Cek apakah CategoryID ada di tabel posts
	exists, err := s.repo.CheckCategoryExists(post.CategoryID)
	if err != nil {
		return Post{}, errors.New("failed to check category: " + err.Error())
	}
	if !exists {
		return Post{}, errors.New("invalid category ID")
	}

	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Post{}, errors.New("username not found in context")
	}

	// Set CreatedBy field
	post.Author = username
	post.CreatedBy = username
	post.ModifiedBy = username

	// Set status to draft for authors
	post.Status = "draft" // Set default status to draft

	if err := s.repo.InsertPost(&post); err != nil {
		return Post{}, errors.New("failed to add new post. " + err.Error())
	}
	log.Print(input.Media)
	// Jika media ada, cek validitasnya dan update post_id
	if len(input.Media) > 0 {
		validMedia, err := s.mediaRepo.CheckMediaExists(input.Media)
		if err != nil {
			return Post{}, errors.New("failed to check media: " + err.Error())
		}
		if len(validMedia) == 0 {
			return Post{}, errors.New("no valid media found")
		}

		// Update post_id di media valid
		if err := s.mediaRepo.UpdateMediaPostID(validMedia, int64(post.ID)); err != nil {
			return Post{}, errors.New("failed to update media post_id: " + err.Error())
		}
	}

	return post, nil
}

// UpdatePost updates an existing post
// UpdatePostService updates an existing post
func (s *postService) UpdatePostService(ctx *gin.Context) (Post, error) {
	var input PostInput

	// Parse ID dari URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return Post{}, errors.New("invalid ID format")
	}

	// Cek apakah post ada di database
	existingPost, err := s.repo.SelectPostByID(int64(id))
	if err != nil {
		return Post{}, errors.New("post not found")
	}

	// Ambil username dari context
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Post{}, errors.New("username not found in context")
	}

	// Ambil role user dari context
	role, _ := ctx.Value("role").(string)

	// Validasi izin: hanya penulis, admin, atau editor yang bisa update
	if existingPost.Author != username && role != "admin" && role != "editor" {
		return Post{}, errors.New("you do not have permission to edit this post")
	}

	// Bind JSON request ke struct input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return Post{}, errors.New("invalid input data")
	}

	// Validasi struktur post
	post := input.Post
	if err := s.validator.Struct(post); err != nil {
		return Post{}, err
	}

	// Cek apakah CategoryID valid
	exists, err := s.repo.CheckCategoryExists(post.CategoryID)
	if err != nil {
		return Post{}, errors.New("failed to check category: " + err.Error())
	}
	if !exists {
		return Post{}, errors.New("invalid category ID")
	}

	// Pertahankan CreatedBy dan CreatedAt asli
	post.ID = existingPost.ID
	post.CreatedBy = existingPost.CreatedBy
	post.CreatedAt = existingPost.CreatedAt

	// Update informasi yang dimodifikasi
	post.ModifiedBy = username
	post.UpdatedAt = time.Now()

	// Validasi izin: hanya penulis, admin, atau editor yang bisa update
	if existingPost.Author != username && role != "admin" && role != "editor" {
		return Post{}, errors.New("you do not have permission to edit this post")
	}

	if post.Status == "" {
		post.Status = existingPost.Status

	}
	// Update post di database
	if err := s.repo.UpdatePost(post); err != nil {
		return Post{}, errors.New("failed to update post: " + err.Error())
	}

	// Cek dan update media jika ada
	if len(input.Media) > 0 {
		validMedia, err := s.mediaRepo.CheckMediaExists(input.Media)
		if err != nil {
			return Post{}, errors.New("failed to check media: " + err.Error())
		}
		if len(validMedia) == 0 {
			return Post{}, errors.New("no valid media found")
		}

		// Update post_id di media valid
		if err := s.mediaRepo.UpdateMediaPostID(validMedia, int64(post.ID)); err != nil {
			return Post{}, errors.New("failed to update media post_id: " + err.Error())
		}
	}

	return post, nil
}

// DeletePost deletes a post
func (s *postService) DeletePostService(ctx *gin.Context) error {
	// Parse the post ID from the URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Check if the post exists
	existingPost, err := s.repo.SelectPostByID(int64(id))
	if err != nil {
		return errors.New("post not found")
	}
	// Check if the post exists
	_, err = s.repo.SelectPostByID(int64(id))
	if err != nil {
		return errors.New("post not found")
	}

	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return errors.New("username not found in context")
	}
	// Set DeletedBy field
	existingPost.DeletedBy = username

	// Proceed to delete the post by updating DeletedBy and soft-deleting
	if err := s.repo.DeletePost(existingPost); err != nil {
		return errors.New("failed to delete post: " + err.Error())
	}

	return nil
}
