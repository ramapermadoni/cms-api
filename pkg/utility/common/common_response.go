package common

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
	AdditionalInfo
}

type Pagination struct {
	CurrentPage int   `json:"current_page"`
	TotalPages  int   `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
	PerPage     int   `json:"per_page"`
}

type AdditionalInfo struct {
	// TotalData int64  `json:"total_data"`
	TraceId string `json:"trace_id"`
}

func GenerateSuccessResponse(ctx *gin.Context, message string) {
	ctx.JSON(
		http.StatusOK,
		GenerateSuccessMessage(message),
	)
}

func GenerateSuccessResponseWithData(ctx *gin.Context, message string, data interface{}) {
	ctx.JSON(
		http.StatusOK,
		GenerateSuccessMessageWithData(message, data),
	)
}

func GenerateSuccessResponseWithListData(ctx *gin.Context, message string, total int64, data interface{}, page, limit int) {
	totalPages := int(math.Ceil(float64(total) / float64(limit))) // Hitung total halaman
	ctx.JSON(
		http.StatusOK,
		APIResponse{
			Success: true,
			Message: message,
			Data:    data,
			AdditionalInfo: AdditionalInfo{
				// TotalData: total,
				TraceId: ctx.GetString("trace_id"),
			},
			Pagination: &Pagination{
				CurrentPage: page,
				TotalPages:  totalPages,
				TotalItems:  total,
				PerPage:     limit,
			},
		},
	)
}

func GenerateErrorResponse(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(
		http.StatusBadRequest,
		GenerateErrorMessage(ctx, message),
	)
}

func GenerateSuccessMessage(message string) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    nil,
	}
}

func GenerateSuccessMessageWithData(message string, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func GenerateSuccessMessageWithListData(message string, total int64, data interface{}, page, limit int) APIResponse {
	totalPages := int((total + int64(limit) - 1) / int64(limit)) // Hitung total halaman dengan pembulatan ke atas

	return APIResponse{
		Success:        true,
		Message:        message,
		Data:           data,
		AdditionalInfo: AdditionalInfo{
			// TotalData: total,
		},
		Pagination: &Pagination{
			CurrentPage: page,
			TotalPages:  totalPages,
			TotalItems:  total,
			PerPage:     limit,
		},
	}
}

func GenerateErrorMessage(ctx *gin.Context, message string) APIResponse {
	return APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
		AdditionalInfo: AdditionalInfo{
			TraceId: ctx.GetString("trace_id"),
		},
	}
}
