package login

import (
	"github.com/kp/pager/common"
	"github.com/kp/pager/login/models"
)

type LoginRequest struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	UserType       string `json:"user_type"`
	UserID         string `json:"user_id"`
	PermissionName string `json:"permission_name,omitempty"`
	Name           string `json:"name"`
	Email          string `json:"email"`
}

type LoginResponse struct {
	common.Response
}

type UserWithPermissions struct {
	User        models.User         `json:"user"`
	Permissions []models.Permission `json:"permissions"`
}

type UserWithPermissionsResponse struct {
	common.Response
	Data struct {
		User        models.User         `json:"user"`
		Permissions []models.Permission `json:"permissions"`
	} `json:"data"`
}

type AllUsersResponse struct {
	common.Response
	Data []UserWithPermissions `json:"data"`
}

type AllPermissionsResponse struct {
	common.Response
	Data []models.Permission `json:"data"`
}

type ChangePasswordRequest struct {
	UserName    string `json:"user_name"`
	NewPassword string `json:"new_password"`
}

type AddPermissionRequest struct {
	UserID       int64  `json:"user_id"`
	PermissionID int64  `json:"permission_id"`
	CreatedBy    string `json:"created_by"`
}

type RemovePermissionRequest struct {
	UserID       int64  `json:"user_id"`
	PermissionID int64  `json:"permission_id"`
	CreatedBy    string `json:"created_by"`
}
