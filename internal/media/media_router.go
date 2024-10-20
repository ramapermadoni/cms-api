package media

import (
	"cms-api/internal/database/connection"
	"cms-api/middlewares"
	"cms-api/pkg/utility/common"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/media")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.Logging())
	{
		// Rute untuk mengunggah media (semua role)
		api.POST("", middlewares.RoleMiddleware("admin", "editor", "author"), CreateMediaRouter)        // Create
		api.POST("/upload", middlewares.RoleMiddleware("admin", "editor", "author"), UploadMediaRouter) // Upload

		// Rute untuk membaca media (semua role)
		api.GET("", GetAllMediaRouter)  // Read (List)
		api.GET("/:id", GetMediaRouter) // Read (By ID)

		// Rute untuk mengedit media (admin dan editor)
		api.PUT("/:id", middlewares.RoleMiddleware("admin", "editor"), UpdateMediaRouter) // Update

		// Rute untuk menghapus media (admin dan editor)
		api.DELETE("/:id", middlewares.RoleMiddleware("admin", "editor"), DeleteMediaRouter) // Delete

	}
}
func GetAllMediaRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewMediaService(repo)

	media, total, err := svc.GetAllMediaService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	// data := gin.H{"total": total, "data": media}
	common.GenerateSuccessResponseWithListData(ctx, "successfully retrieved all media data", total, media)
}
func GetMediaRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewMediaService(repo)

	media, err := svc.GetMediaByIDService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully retrieved media data", media)
}

func CreateMediaRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB) // Gunakan variabel yang diekspor
	svc := NewMediaService(repo)

	media, err := svc.CreateMediaService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully added media data", media)
}

func UpdateMediaRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewMediaService(repo)
	media, err := svc.UpdateMediaService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully updated media data", media)
}

func DeleteMediaRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewMediaService(repo)

	if err := svc.DeleteMediaService(ctx); err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponse(ctx, "successfully deleted media data")
}
func UploadMediaRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewMediaService(repo)

	media, err := svc.UploadMediaService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	// Kembalikan data media untuk disimpan ke tabel posts oleh user
	common.GenerateSuccessResponseWithData(ctx, "successfully uploaded media", media)
}
