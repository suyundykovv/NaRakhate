package api

import (
	"encoding/json"
	"net/http"

	"Aitu-Bet/utils"

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.Email == "test@example.com" && request.Password == "password123" {
		token, err := utils.GenerateJWT(1, request.Email)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my house body!"))
}
