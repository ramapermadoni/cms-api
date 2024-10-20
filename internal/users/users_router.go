package users

import (
	"cms-api/internal/database/connection"
	"cms-api/middlewares"
	"cms-api/pkg/utility/common"

	"github.com/gin-gonic/gin"
)

func Initiator(router *gin.Engine) {
	api := router.Group("/api/user")
	api.Use(middlewares.JwtMiddleware())
	api.Use(middlewares.Logging())
	{
		// api.POST("", CreateUserRouter)       // Create
		// api.GET("", GetAllUserRouter)        // Read (List)
		// api.GET("/:id", GetUserRouter)       // Read (By ID)
		// api.PUT("/:id", UpdateUserRouter)    // Update
		// api.DELETE("/:id", DeleteUserRouter) // Delete
		// Rute untuk menambah user (hanya admin)
		api.POST("", middlewares.RoleMiddleware("admin"), CreateUserRouter) // Create

		// Rute untuk membaca data user (semua role)
		api.GET("", middlewares.RoleMiddleware("admin"), GetAllUserRouter)  // Read (List)
		api.GET("/:id", middlewares.RoleMiddleware("admin"), GetUserRouter) // Read (By ID)

		// Rute untuk mengubah user (hanya admin)
		api.PUT("/:id", middlewares.RoleMiddleware("admin"), UpdateUserRouter) // Update

		// Rute untuk menghapus user (hanya admin)
		api.DELETE("/:id", middlewares.RoleMiddleware("admin"), DeleteUserRouter) // Delete

		// // Rute untuk mengubah role pengguna (hanya admin)
		// api.PUT("/role/:id", middlewares.RoleMiddleware("admin"), UpdateUserRoleRouter) // Update Role

	}
}
func GetAllUserRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewUserService(repo)

	user, total, err := svc.GetAllUserService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}

	// data := gin.H{"total": total, "data": user}
	common.GenerateSuccessResponseWithListData(ctx, "successfully retrieved all user data", total, user)
}
func GetUserRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewUserService(repo)

	user, err := svc.GetUserByIDService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully retrieved user data", user)
}

func CreateUserRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB) // Gunakan variabel yang diekspor
	svc := NewUserService(repo)

	user, err := svc.CreateUserService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully added user data", user)
}

func UpdateUserRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewUserService(repo)
	user, err := svc.UpdateUserService(ctx)
	if err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponseWithData(ctx, "successfully updated user data", user)
}

func DeleteUserRouter(ctx *gin.Context) {
	repo := NewRepository(connection.DB)
	svc := NewUserService(repo)

	if err := svc.DeleteUserService(ctx); err != nil {
		common.GenerateErrorResponse(ctx, err.Error())
		return
	}
	common.GenerateSuccessResponse(ctx, "successfully deleted user data")
}
