package login

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kp/pager/common"
	"github.com/kp/pager/login/models"
)

type AuthController struct {
	authService AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: *authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user type
	validTypes := map[string]bool{
		UserTypeAdmin:     true,
		UserTypeMarketing: true,
		UserTypeNormal:    true,
	}
	if !validTypes[req.UserType] {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":       "Invalid user type",
			"valid_types": []string{UserTypeAdmin, UserTypeMarketing, UserTypeNormal},
		})
		return
	}

	// Validate username ends with wealthy.in
	if !strings.HasSuffix(req.Username, "@wealthy.in") {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must end with wealthy.in",
		})
		return
	}

	encryptedPass := common.Encryptbase64(req.Password)
	user, permissions, err := c.authService.RegisterUser(ctx.Request.Context(), req.Username, encryptedPass, req.UserType, req.Name, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = req.Password
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data": gin.H{
			"user":        user,
			"permissions": permissions,
		},
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Username and password are required",
		})
		return
	}

	encryptedPass := common.Encryptbase64(req.Password)
	user, permissions, err := c.authService.Login(ctx.Request.Context(), req.Username, encryptedPass)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Authentication failed",
			"message": "Invalid username or password",
		})
		return
	}

	permNames := make([]string, len(permissions))
	for i, p := range permissions {
		permNames[i] = p.Name
	}

	token, expiryTime, err := GenerateToken(user, permNames)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Login processing error",
			"message": "Could not generate authentication token",
		})
		return
	}

	user.Password = ""
	isAdmin := false
	if user.UserType == UserTypeAdmin {
		isAdmin = true
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token":       token,
			"user":        user,
			"permissions": permNames,
			"expiry_time": expiryTime.Format(time.RFC3339),
			"expiry_secs": int(time.Until(expiryTime).Seconds()),
			"is_admin":    isAdmin,
		},
	})
}

func (c *AuthController) AddPermission(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.authService.AddPermission(ctx.Request.Context(), req.UserID, req.PermissionName); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to add permission",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Permission added successfully",
	})
}

func (c *AuthController) GetPermissions(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Missing user ID parameter",
			"error":   "user_id is required",
		})
		return
	}
	permissions, err := c.authService.GetUserPermissions(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get user permissions",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Permissions retrieved successfully",
		"data":    permissions,
	})
}

func (c *AuthController) GetUserDetails(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Missing user ID parameter",
			"error":   "user_id is required",
		})
		return
	}

	// Get user details
	user, err := c.authService.userRepo.GetByID(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get user details",
			"error":   err.Error(),
		})
		return
	}

	// Get user permissions
	permissions, err := c.authService.GetUserPermissions(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get user permissions",
			"error":   err.Error(),
		})
		return
	}

	// Decrypt password before returning
	decryptedPass, err := common.DecryptBase64(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to process user data",
			"error":   err.Error(),
		})
		return
	}
	user.Password = decryptedPass
	response := UserWithPermissionsResponse{
		Response: common.Response{
			Status: true,
		},
		Data: struct {
			User        models.User         `json:"user"`
			Permissions []models.Permission `json:"permissions"`
		}{
			User:        *user,
			Permissions: permissions,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) GetAllPermissions(ctx *gin.Context) {
	permissions, err := c.authService.GetAllPermissions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to get permissions",
			"error":   err.Error(),
		})
		return
	}

	response := AllPermissionsResponse{
		Response: common.Response{
			Status: true,
		},
		Data: permissions,
	}
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthController) AddUserPermission(ctx *gin.Context) {
	var req AddPermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.authService.AddUserPermission(ctx.Request.Context(), req.UserID, req.PermissionID, req.CreatedBy); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to add user permission",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User permission added successfully",
	})
}

func (c *AuthController) RemoveUserPermission(ctx *gin.Context) {
	var req RemovePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.authService.RemoveUserPermission(ctx.Request.Context(), req.UserID, req.PermissionID, req.CreatedBy); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to remove user permission",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User permission removed successfully",
	})
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	if err := c.authService.ChangePassword(ctx.Request.Context(), req.UserName, req.NewPassword); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Failed to change password",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Password changed successfully",
	})
}

func (c *AuthController) GetAllUsers(ctx *gin.Context) {
	usersWithPerms, err := c.authService.GetAllUsersWithPermissions(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := AllUsersResponse{
		Response: common.Response{
			Status: true,
		},
		Data: usersWithPerms,
	}
	ctx.JSON(http.StatusOK, response)
}
