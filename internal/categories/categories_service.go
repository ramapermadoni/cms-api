package categories

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryService interface {
	GetAllCategoryService(ctx *gin.Context) ([]Category, int64, int, int, error)
	GetCategoryByIDService(ctx *gin.Context) (Category, error)
	CreateCategoryService(ctx *gin.Context) (Category, error)
	UpdateCategoryService(ctx *gin.Context) (Category, error)
	DeleteCategoryService(ctx *gin.Context) error
}
type categoryService struct {
	repo      Repository
	validator *validator.Validate
}

func NewCategoryService(repo Repository) CategoryService {
	return &categoryService{
		repo:      repo,
		validator: validator.New(),
	}
}

// GetAllCategoryService retrieves all category with pagination and search
func (s *categoryService) GetAllCategoryService(ctx *gin.Context) ([]Category, int64, int, int, error) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	orderBy := ctx.DefaultQuery("orderBy", "updated_at") // Default ke 'updated_at'
	sort := ctx.DefaultQuery("sort", "desc")

	categories, total, err := s.repo.SelectAllCategory(page, limit, search, orderBy, sort)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return categories, total, page, limit, nil
}
func (s *categoryService) GetCategoryByIDService(ctx *gin.Context) (Category, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return Category{}, errors.New("invalid ID")
	}

	return s.repo.SelectCategoryByID(id)
}

func (s *categoryService) CreateCategoryService(ctx *gin.Context) (Category, error) {
	var category Category
	if err := ctx.ShouldBindJSON(&category); err != nil {
		return Category{}, errors.New("invalid input data")
	}

	if err := s.validator.Struct(category); err != nil {
		return Category{}, err
	}
	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Category{}, errors.New("username not found in context")
	}

	// Set CreatedBy field
	category.CreatedBy = username
	category.ModifiedBy = username

	if err := s.repo.InsertCategory(&category); err != nil {
		// return Category{}, err
		return Category{}, errors.New("failed to add new category. " + err.Error())
	}

	return category, nil
}

// UpdateCategory updates an existing category
func (s *categoryService) UpdateCategoryService(ctx *gin.Context) (Category, error) {
	var category Category

	// Parse the ID from the URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return Category{}, errors.New("invalid ID format")
	}

	// Check if the category exists in the database
	existingCategory, err := s.repo.SelectCategoryByID(int64(id))
	if err != nil {
		return Category{}, errors.New("category not found")
	}

	// Bind the JSON request to the category struct
	if err := ctx.ShouldBindJSON(&category); err != nil {
		return Category{}, errors.New("invalid input data")
	}

	if err := s.validator.Struct(category); err != nil {
		return Category{}, err
	}

	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return Category{}, errors.New("username not found in context")
	}
	category.ID = existingCategory.ID
	category.CreatedBy = existingCategory.CreatedBy
	category.CreatedAt = existingCategory.CreatedAt
	// Set Modified field
	category.ModifiedBy = username
	category.UpdatedAt = time.Now()

	// Update the category in the database
	if err := s.repo.UpdateCategory(category); err != nil {
		return Category{}, errors.New("failed to update category. " + err.Error())
	}

	return category, nil
}

// DeleteCategory deletes a category
func (s *categoryService) DeleteCategoryService(ctx *gin.Context) error {
	// Parse the category ID from the URL parameter
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.New("invalid ID format")
	}

	// Check if the category exists
	existingCategory, err := s.repo.SelectCategoryByID(int64(id))
	if err != nil {
		return errors.New("category not found")
	}
	// Check if the category exists
	_, err = s.repo.SelectCategoryByID(int64(id))
	if err != nil {
		return errors.New("category not found")
	}

	// Retrieve the username from context and check if it exists
	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return errors.New("username not found in context")
	}
	// Set DeletedBy field
	existingCategory.DeletedBy = username

	// Proceed to delete the category by updating DeletedBy and soft-deleting
	if err := s.repo.DeleteCategory(existingCategory); err != nil {
		return errors.New("failed to delete category: " + err.Error())
	}

	return nil
}
