package categories

import (
	"cms-api/internal/database/connection"
	"cms-api/middlewares"
	"cms-api/pkg/utility/common"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/category")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.Logging())
	{
		// Rute untuk menambah kategori (admin dan editor)
		api.POST("", middlewares.RoleMiddleware("admin", "editor"), CreateCategoryRouter) // Create

		// Rute untuk membaca kategori (semua role)
		api.GET("", GetAllCategoryRouter)  // Read (List)
		api.GET("/:id", GetCategoryRouter) // Read (By ID)

		// Rute untuk mengedit kategori (admin dan editor)
		api.PUT("/:id", middlewares.RoleMiddleware("admin", "editor"), UpdateCategoryRouter) // Update

		// Rute untuk menghapus kategori (admin dan editor)
		api.DELETE("/:id", middlewares.RoleMiddleware("admin", "editor"), DeleteCategoryRouter) // Delete

	}
}
func GetAllCategoryRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewCategoryService(repo)

	category, total, page, limit, err := svc.GetAllCategoryService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	// data := gin.H{"total": total, "data": category}
	common.GenerateSuccessResponseWithListData(ctx, "successfully retrieved all category data", total, category, page, limit)
}
func GetCategoryRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewCategoryService(repo)

	category, err := svc.GetCategoryByIDService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully retrieved category data", category)
}

func CreateCategoryRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB) // Gunakan variabel yang diekspor
	svc := NewCategoryService(repo)

	category, err := svc.CreateCategoryService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully added category data", category)
}

func UpdateCategoryRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewCategoryService(repo)
	category, err := svc.UpdateCategoryService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully updated category data", category)
}

func DeleteCategoryRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewCategoryService(repo)

	if err := svc.DeleteCategoryService(ctx); err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponse(ctx, "successfully deleted category data")
}
