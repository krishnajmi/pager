package login

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/kp/pager/login/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, username, password, userType, name, email string) (*models.User, error) {
	user := models.User{
		Username: username,
		Password: password,
		UserType: userType,
		Name:     name,
		EmailID:  email,
	}
	err := r.db.Create(&user).Error
	return &user, err
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	return &user, err
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) Create(ctx context.Context, name, description string) (*models.Permission, error) {
	permission := models.Permission{
		Name:        name,
		Description: description,
	}
	err := r.db.Create(&permission).Error
	return &permission, err
}

func (r *PermissionRepository) GetByName(ctx context.Context, name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

func (r *PermissionRepository) GetAll(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

type UserPermissionRepository struct {
	db *gorm.DB
}

func NewUserPermissionRepository(db *gorm.DB) *UserPermissionRepository {
	return &UserPermissionRepository{db: db}
}

func (r *UserPermissionRepository) Add(ctx context.Context, userID string, permissionID uint) error {
	userPermission := models.UserPermission{
		UserID:       userID,
		PermissionID: permissionID,
	}
	return r.db.Create(&userPermission).Error
}

func (r *UserPermissionRepository) Remove(ctx context.Context, userID string, permissionID uint) error {
	return r.db.Where("user_id = ? AND permission_id = ?", userID, permissionID).Delete(&models.UserPermission{}).Error
}

func (r *UserRepository) UpdateCreatedBy(ctx context.Context, userID string, createdBy string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("created_by", createdBy).Error
}

func (r *UserRepository) UpdatePassword(ctx context.Context, username string, newPassword string) error {
	return r.db.Model(&models.User{}).Where("username = ?", username).Update("password", newPassword).Error
}

func (r *UserPermissionRepository) GetForUser(ctx context.Context, userID string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Joins("JOIN pager_users_permissions ON pager_users_permissions.permission_id = pager_permissions.id").
		Where("pager_users_permissions.user_id = ?", userID).
		Find(&permissions).Error
	return permissions, err
}
