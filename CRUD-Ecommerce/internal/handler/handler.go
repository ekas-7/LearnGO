package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	User     *UserHandler
	Product  *ProductHandler
	Category *CategoryHandler
	Order    *OrderHandler
}

func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, fmt.Errorf("user_id not found in context")
	}

	userID, ok := userIDValue.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user_id type")
	}

	return userID, nil
}

func isAdminUser(c *gin.Context) bool {
	role, exists := c.Get("role")
	if !exists {
		return false
	}

	roleStr, ok := role.(string)
	if !ok {
		return false
	}

	return roleStr == "admin"
}
