package login

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kp/pager/common"
	"github.com/kp/pager/login/models"

	log "github.com/sirupsen/logrus"
)

type AuthService struct {
	userRepo       *UserRepository
	permissionRepo *PermissionRepository
	userPermRepo   *UserPermissionRepository
}

func NewAuthService(
	userRepo *UserRepository,
	permissionRepo *PermissionRepository,
	userPermRepo *UserPermissionRepository,
) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
		userPermRepo:   userPermRepo,
	}
}

func (s *AuthService) RegisterUser(ctx context.Context, username, password, userType, name, email string) (*models.User, []models.Permission, error) {
	// Business logic and validation here
	user, err := s.userRepo.Create(ctx, username, password, userType, name, email)
	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
			"userType": userType,
			"name":     name,
			"email":    email,
		}).WithError(err).Error("Failed to create user")
		return nil, nil, err
	}

	// Assign default permissions based on user type
	var permissions []string
	switch userType {
	case UserTypeAdmin:
		permissions = DefaultAdminPermissions
	case UserTypeMarketing:
		permissions = DefaultMarketingPermissions
	default:
		permissions = DefaultUserPermissions
	}

	// Add all permissions
	var assignedPerms []models.Permission
	for _, perm := range permissions {
		if err := s.AddPermission(ctx, strconv.FormatInt(user.ID, 10), perm); err != nil {
			log.WithFields(log.Fields{
				"permission": perm,
				"userId":     user.ID,
			}).WithError(err).Error("Failed to add permission")
			return nil, nil, err
		}
		// Get the permission details to return
		p, err := s.permissionRepo.GetByName(ctx, perm)
		if err != nil {
			log.WithFields(log.Fields{
				"permission": perm,
			}).WithError(err).Error("Failed to get permission details")
			return nil, nil, err
		}
		assignedPerms = append(assignedPerms, *p)
	}

	return user, assignedPerms, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, []models.Permission, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
		}).WithError(err).Error("Failed to get user")
		return nil, nil, err
	}
	if user.Password != password {
		log.WithFields(log.Fields{
			"username": username,
		}).Error("Invalid password attempt")
		return nil, nil, fmt.Errorf("invalid password")
	}

	perms, err := s.GetUserPermissions(ctx, strconv.FormatInt(user.ID, 10))
	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
			"userId":   user.ID,
		}).WithError(err).Error("Failed to get user permissions")
		return nil, nil, fmt.Errorf("failed to get permissions: %v", err)
	}

	return user, perms, nil
}

func (s *AuthService) AddPermission(ctx context.Context, userID string, permissionName string) error {
	perm, err := s.permissionRepo.GetByName(ctx, permissionName)
	if err != nil {
		log.WithFields(log.Fields{
			"userId":     userID,
			"permission": permissionName,
		}).WithError(err).Error("Failed to get permission")
		return err
	}
	err = s.userPermRepo.Add(ctx, userID, perm.ID)
	if err != nil {
		log.WithFields(log.Fields{
			"userId":     userID,
			"permission": permissionName,
		}).WithError(err).Error("Failed to add permission to user")
		return err
	}
	return nil
}

func (s *AuthService) GetUserPermissions(ctx context.Context, userID string) ([]models.Permission, error) {
	perms, err := s.userPermRepo.GetForUser(ctx, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"userId": userID,
		}).WithError(err).Error("Failed to get user permissions")
		return nil, err
	}
	return perms, nil
}

func (s *AuthService) GetAllPermissions(ctx context.Context) ([]models.Permission, error) {
	perms, err := s.permissionRepo.GetAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to get all permissions")
		return nil, err
	}
	return perms, nil
}

func (s *AuthService) AddUserPermission(ctx context.Context, userID int64, permissionID int64, createdBy string) error {
	err := s.userPermRepo.Add(ctx, strconv.FormatInt(userID, 10), uint(permissionID))
	if err != nil {
		log.WithFields(log.Fields{
			"userId":       userID,
			"permissionId": permissionID,
			"createdBy":    createdBy,
		}).WithError(err).Error("Failed to add user permission")
		return err
	}
	return nil
}

func (s *AuthService) RemoveUserPermission(ctx context.Context, userID int64, permissionID int64, createdBy string) error {
	err := s.userPermRepo.Remove(ctx, strconv.FormatInt(userID, 10), uint(permissionID))
	if err != nil {
		log.WithFields(log.Fields{
			"userId":       userID,
			"permissionId": permissionID,
			"createdBy":    createdBy,
		}).WithError(err).Error("Failed to remove user permission")
		return err
	}
	return nil
}

func (s *AuthService) ChangePassword(ctx context.Context, username string, newPassword string) error {
	encryptedPass := common.Encryptbase64(newPassword)
	err := s.userRepo.UpdatePassword(ctx, username, encryptedPass)
	if err != nil {
		log.WithFields(log.Fields{
			"username": username,
		}).WithError(err).Error("Failed to update password")
		return err
	}
	return nil
}

func (s *AuthService) GetAllUsersWithPermissions(ctx context.Context) ([]UserWithPermissions, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to get all users")
		return nil, err
	}

	var result []UserWithPermissions
	for _, user := range users {
		// Decrypt password before returning
		decryptedPass, err := common.DecryptBase64(user.Password)
		if err != nil {
			log.WithFields(log.Fields{
				"userId": user.ID,
			}).WithError(err).Error("Failed to decrypt user password")
			return nil, err
		}
		user.Password = decryptedPass

		perms, err := s.GetUserPermissions(ctx, strconv.FormatInt(user.ID, 10))
		if err != nil {
			log.WithFields(log.Fields{
				"userId": user.ID,
			}).WithError(err).Error("Failed to get user permissions")
			return nil, err
		}
		result = append(result, UserWithPermissions{
			User:        user,
			Permissions: perms,
		})
	}

	return result, nil
}
