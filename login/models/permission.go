package models

import (
	"context"
	"time"

	"github.com/kp/pager/databases/sql"
)

const PermissionTableName = "pager_permissions"

func (Permission) TableName() string {
	return PermissionTableName
}

type Permission struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"size:255;not null;unique"`
	Description string    `json:"description" gorm:"size:500"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func CreatePermission(ctx context.Context, tx interface{}, name, description string) (*Permission, error) {
	db := sql.GetOrmQuearyable(ctx, tx)
	permission := Permission{
		Name:        name,
		Description: description,
	}
	err := db.Create(&permission).Error
	return &permission, err
}

func GetPermissionByName(ctx context.Context, tx interface{}, name string) (*Permission, error) {
	var permission Permission
	db := sql.GetOrmQuearyable(ctx, tx)
	err := db.Where("name = ?", name).First(&permission).Error
	return &permission, err
}

func GetPermissionByID(ctx context.Context, tx interface{}, id uint) (*Permission, error) {
	var permission Permission
	db := sql.GetOrmQuearyable(ctx, tx)
	err := db.First(&permission, id).Error
	return &permission, err
}

func (permission *Permission) Save(ctx context.Context, tx interface{}) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	return db.Save(permission).Error
}
