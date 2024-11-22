package handlers

import (
	"fmt"
	"net/http"
	"tender-backend/internal/http/token"
	"tender-backend/internal/pkg/config"
	request_model "tender-backend/model/request"
	response_model "tender-backend/model/response"

	"gorm.io/gorm/utils"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with email, username, and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body request_model.CreateUserReq true "User registration request"
// @Success 201 {object} response_model.LoginRes "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Router /register [post]
func (h *HTTPHandler) Register(c *gin.Context) {
	var req request_model.CreateUserReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	if req.Email == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "username or email cannot be empty"})
		return
	}

	if !config.IsValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid email format"})
		return
	}

	availableRoles := []string{"client", "contractor"}

	if !utils.Contains(availableRoles, req.Role) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid role"})
		return
	}

	hashedPassword, err := config.HashPassword(req.Password)
	if err != nil {
		fmt.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error", "err": err.Error()})
		return
	}

	req.Password = hashedPassword
	user, err := h.UserService.CreateUser(&req)
	if err != nil {
		fmt.Printf("Error creating user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	tkn, err := token.GenerateJWT(user.ID, user.Role)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := response_model.LoginRes{
		Token: tkn,
		Role:  user.Role,
	}

	c.JSON(201, res)
}

// Login godoc
// @Summary Login a user
// @Description Authenticate user with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body request_model.LoginUserReq true "User login credentials"
// @Success 200 {object} response_model.LoginRes "JWT tokens"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Invalid email or password"
// @Router /login [post]
func (h *HTTPHandler) Login(c *gin.Context) {
	req := request_model.LoginUserReq{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid request payload": err.Error()})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Username and password are required"})
		return
	}

	user, err := h.UserService.GetByUsername(req.Username)
	if err != nil {
		fmt.Printf("Error fetching user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}

	if !config.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	tkn, err := token.GenerateJWT(user.ID, user.Role)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res := response_model.LoginRes{
		Token: tkn,
		Role:  user.Role,
	}

	c.JSON(http.StatusOK, res)
}
