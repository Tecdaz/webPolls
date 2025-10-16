package handlers

import (
	sqlc "webpolls/db/sqlc"

	"webpolls/dataconvertion"

	"database/sql"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	queries *sqlc.Queries
}

func NewUserHandler(queries *sqlc.Queries) *userHandler {
	return &userHandler{queries: queries}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *userHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// checkeos para asegurar que no exista ni el usuario ni el mail
	_, err := h.queries.GetUserByUsername(c, req.Username)
	if err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	_, err = h.queries.GetUserByEmail(c, req.Email)
	if err == nil {
		c.JSON(400, gin.H{"error": "Email already exists"})
		return
	}

	user, err := h.queries.CreateUser(c, sqlc.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(200, gin.H{"user": user})
}

func (h *userHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	UserIDInt32, err := dataconvertion.ConvertTo32(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.queries.DeleteUser(c, UserIDInt32)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (h *userHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	UserIDInt32, err := dataconvertion.ConvertTo32(userID)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.queries.GetUserByID(c, UserIDInt32)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(200, gin.H{"user": user})
}

/*
DISCLAIMER: revisar esta parte por que ya estaba re quemado y la hice con gpt
*/

// struct para la actualizacion
type UpdateUserRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (h *userHandler) UpdateUser(c *gin.Context) {

	userIDStr := c.Param("id")
	userID, err := dataconvertion.ConvertTo32(userIDStr)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid user ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Verificar que exista el usuario
	_, err = h.queries.GetUserByID(c, userID)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	var username, email, password sql.NullString

	if req.Username != nil {
		if _, err := h.queries.GetUserByUsername(c, *req.Username); err == nil {
			c.JSON(400, gin.H{"error": "Username already exists"})
			return
		}
		username = sql.NullString{String: *req.Username, Valid: true}
	} else {
		username = sql.NullString{Valid: false}
	}

	if req.Email != nil {
		if _, err := h.queries.GetUserByEmail(c, *req.Email); err == nil {
			c.JSON(400, gin.H{"error": "Email already exists"})
			return
		}
		email = sql.NullString{String: *req.Email, Valid: true}
	} else {
		email = sql.NullString{Valid: false}
	}

	if req.Password != nil {
		password = sql.NullString{String: *req.Password, Valid: true}
	} else {
		password = sql.NullString{Valid: false}
	}

	user, err := h.queries.UpdateUser(c, sqlc.UpdateUserParams{
		ID:      userID,
		Column2: username,
		Column3: email,
		Column4: password,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(200, gin.H{"user": user})
}
