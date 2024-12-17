package posts

import (
	"cms-api/internal/database/connection"
	"cms-api/internal/media"
	"cms-api/middlewares"
	"cms-api/pkg/utility/common"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/post")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.Logging())
	{
		// Rute untuk membuat post (semua role)
		api.POST("", middlewares.RoleMiddleware("admin", "editor", "author"), CreatePostRouter) // Create

		// Rute untuk membaca post (semua role)
		api.GET("", GetAllPostRouter)  // Read (List)
		api.GET("/:id", GetPostRouter) // Read (By ID)

		// Rute untuk mengedit post (admin, editor, author milik sendiri)
		api.PUT("/:id", middlewares.RoleMiddleware("admin", "editor", "author"), UpdatePostRouter) // Update

		// Rute untuk menghapus post (admin dan editor)
		api.DELETE("/:id", middlewares.RoleMiddleware("admin", "editor"), DeletePostRouter) // Delete

	}
}
func GetAllPostRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewPostService(repo, media.NewRepository(connection.DB))

	post, total, page, limit, err := svc.GetAllPostService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	// data := gin.H{"total": total, "data": post}
	common.GenerateSuccessResponseWithListData(ctx, "successfully retrieved all post data", total, post, page, limit)
}
func GetPostRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewPostService(repo, media.NewRepository(connection.DB))

	post, err := svc.GetPostByIDService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully retrieved post data", post)
}

func CreatePostRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB) // Gunakan variabel yang diekspor
	svc := NewPostService(repo, media.NewRepository(connection.DB))

	post, err := svc.CreatePostService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully added post data", post)
}

func UpdatePostRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewPostService(repo, media.NewRepository(connection.DB))
	post, err := svc.UpdatePostService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully updated post data", post)
}

func DeletePostRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewPostService(repo, media.NewRepository(connection.DB))

	if err := svc.DeletePostService(ctx); err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponse(ctx, "successfully deleted post data")
}
