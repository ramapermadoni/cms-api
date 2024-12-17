package users

import (
	"cms-api/pkg/utility/common"
	"errors"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserService interface {
	GetAllUserService(ctx *gin.Context) ([]User, int64, int, int, error)
	GetUserByIDService(ctx *gin.Context) (User, error)
	GetProfileByIDService(ctx *gin.Context) (User, error)
	CreateUserService(ctx *gin.Context) (User, error)
	UpdateUserService(ctx *gin.Context) (User, error)
	DeleteUserService(ctx *gin.Context) error
}

type userService struct {
	repo      Repository
	validator *validator.Validate
}

func NewUserService(repo Repository) UserService {
	return &userService{
		repo:      repo,
		validator: validator.New(),
	}
}

// GetAllUserService retrieves all users with pagination and search
func (s *userService) GetAllUserService(ctx *gin.Context) ([]User, int64, int, int, error) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	search := ctx.Query("search")
	orderBy := ctx.DefaultQuery("orderBy", "updated_at") // Default ke 'updated_at'
	sort := ctx.DefaultQuery("sort", "desc")
	users, total, err := s.repo.SelectAllUser(page, limit, search, orderBy, sort)
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, 0, 0, 0, err
	}

	if len(users) == 0 {
		log.Println("No users found.")
	}

	for i := range users {
		users[i].Password = "*****" // Mask password field
	}

	return users, total, page, limit, nil
}

func (s *userService) GetUserByIDService(ctx *gin.Context) (User, error) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return User{}, errors.New("invalid ID")
	}

	user, err := s.repo.SelectUserByID(id)
	user.Password = "*****" // Mask password field
	if err != nil {
		log.Printf("User not found with ID: %d", id)
		return User{}, err
	}

	return user, nil
}
func (s *userService) GetProfileByIDService(ctx *gin.Context) (User, error) {
	// Ambil id_user dari konteks yang di-set oleh JwtMiddleware
	idUser, exists := ctx.Get("id_user")
	if !exists {
		log.Println("id_user not found in context")
		return User{}, errors.New("id_user not found in context")
	}

	// Konversi idUser ke int64 jika awalnya berupa int atau jenis lain yang dapat dikonversi
	var id int64
	switch v := idUser.(type) {
	case int:
		id = int64(v)
	case int64:
		id = v
	default:
		log.Println("Invalid id_user type in context")
		return User{}, errors.New("invalid id_user type in context")
	}

	log.Printf("Fetching profile for id_user: %d", id)

	// Gunakan id untuk mengambil data user dari repository
	user, err := s.repo.SelectUserByID(id)
	user.Password = "*****" // Mask password field

	if err != nil {
		log.Printf("User not found with ID: %d", id)
		return User{}, err
	}

	return user, nil
}

func (s *userService) CreateUserService(ctx *gin.Context) (User, error) {
	var user User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		return User{}, errors.New("invalid input data")
	}

	// Validasi username
	if !isValidUsername(user.Username) {
		return User{}, errors.New("username can only contain letters and numbers")
	}
	if err := s.validator.Struct(user); err != nil {
		return User{}, err
	}

	hashedPassword, err := common.HashPassword(user.Password)
	if err != nil {
		return User{}, errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return User{}, errors.New("username not found in context")
	}

	user.CreatedBy = username
	user.ModifiedBy = username

	if err := s.repo.InsertUser(&user); err != nil {
		return User{}, errors.New("failed to add new user. " + err.Error())
	}

	return user, nil
}

func (s *userService) UpdateUserService(ctx *gin.Context) (User, error) {
	var user User
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return User{}, errors.New("invalid ID format")
	}

	existingUser, err := s.repo.SelectUserByID(int64(id))
	if err != nil {
		return User{}, errors.New("user not found")
	}

	if err := ctx.ShouldBindJSON(&user); err != nil {
		return User{}, errors.New("invalid input data")
	}

	// Validasi username
	if !isValidUsername(user.Username) {
		return User{}, errors.New("username can only contain letters and numbers")
	}
	if user.Password != "" {
		hashedPassword, err := common.HashPassword(user.Password)
		if err != nil {
			return User{}, errors.New("failed to hash password")
		}
		user.Password = hashedPassword
	}

	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return User{}, errors.New("username not found in context")
	}
	user.ID = existingUser.ID
	user.CreatedBy = existingUser.CreatedBy
	user.CreatedAt = existingUser.CreatedAt
	user.ModifiedBy = username
	user.UpdatedAt = time.Now()

	if err := s.repo.UpdateUser(user); err != nil {
		return User{}, errors.New("failed to update user. " + err.Error())
	}

	return user, nil
}

func (s *userService) DeleteUserService(ctx *gin.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return errors.New("invalid ID format")
	}

	existingUser, err := s.repo.SelectUserByID(int64(id))
	if err != nil {
		return errors.New("user not found")
	}

	username, ok := ctx.Value("username").(string)
	if !ok || username == "" {
		return errors.New("username not found in context")
	}

	existingUser.DeletedBy = username

	if err := s.repo.DeleteUser(existingUser); err != nil {
		return errors.New("failed to delete user: " + err.Error())
	}

	return nil
}

// isValidUsername checks if the username contains only letters and numbers
func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return re.MatchString(username)
}
