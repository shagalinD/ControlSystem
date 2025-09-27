package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
    DB *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
    return &Handler{DB: db}
}

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

func (h *Handler) success(c *gin.Context, data interface{}, message string) {
    c.JSON(http.StatusOK, Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func (h *Handler) error(c *gin.Context, status int, message string) {
    c.JSON(status, Response{
        Success: false,
        Error:   message,
    })
}

func (h *Handler) badRequest(c *gin.Context, message string) {
    h.error(c, http.StatusBadRequest, message)
}

func (h *Handler) notFound(c *gin.Context, message string) {
    h.error(c, http.StatusNotFound, message)
}

func (h *Handler) internalError(c *gin.Context, message string) {
    h.error(c, http.StatusInternalServerError, message)
}