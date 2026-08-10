package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	token, err := h.service.Login(
		c.Request.Context(),
		req.Username,
		req.Password,
	)

	if err != nil {

		log.Println(err)
		if errors.Is(err, ErrInvalidCredentials) {

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
	})
}

func (h *Handler) Logout(c *gin.Context) {

	value, exists := c.Get("sessionID")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "session not found",
		})
		return
	}

	sessionID, ok := value.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid session id",
		})
		return
	}

	err := h.service.Logout(
		c.Request.Context(),
		sessionID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) LogoutAll(c *gin.Context) {
	value, exists := c.Get(ContextUserID)

	if !exists {

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})

		return
	}

	userID, ok := value.(uuid.UUID)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id",
		})
		return
	}

	err := h.service.LogoutAll(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})

		return
	}

	c.Status(http.StatusNoContent)
}
