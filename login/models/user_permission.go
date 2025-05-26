package models

import (
	"context"
	"time"

	"github.com/kp/pager/databases/sql"
)

const UserPermissionTableName = "pager_users_permissions"

func (UserPermission) TableName() string {
	return UserPermissionTableName
}

type UserPermission struct {
	ID           int64     `json:"id" gorm:"primary_key;autoIncrement:true;column:id"`
	UserID       string    `json:"user_id" gorm:"index:idx_user_permissions_user_id;size:36"`
	PermissionID uint      `json:"permission_id" gorm:"index:idx_user_permissions_permission_id"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

/*
CREATE TABLE user_permissions (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id       BIGINT NOT NULL,
    permission_id BIGINT NOT NULL,
    created_at    TIMESTAMP WITHOUT TIME ZONE,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id)
);
CREATE INDEX idx_user_permissions_user_id ON user_permissions (user_id);
CREATE INDEX idx_user_permissions_permission_id ON user_permissions (permission_id);
*/

func AddUserPermission(ctx context.Context, tx interface{}, userID string, permissionID uint) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	userPermission := UserPermission{
		UserID:       userID,
		PermissionID: permissionID,
	}
	return db.Create(&userPermission).Error
}

func GetUserPermissions(ctx context.Context, tx interface{}, userID string) ([]Permission, error) {
	var permissions []Permission
	db := sql.GetOrmQuearyable(ctx, tx)
	err := db.Joins("JOIN user_permissions ON user_permissions.permission_id = permissions.id").
		Where("user_permissions.user_id = ?", userID).
		Find(&permissions).Error
	return permissions, err
}

func DeleteUserPermission(ctx context.Context, tx interface{}, userID string, permissionID uint) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	return db.Where("user_id = ? AND permission_id = ?", userID, permissionID).
		Delete(&UserPermission{}).Error
}
