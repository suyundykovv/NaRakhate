package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProtectedEndpoint(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Welcome authorized user", "user_id": userID})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my house body!"))
}

// {
//     "email":"test@example.com",
//     "password":"password123"
// }
