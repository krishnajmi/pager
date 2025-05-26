package models

import (
	"context"
	"time"

	"github.com/kp/pager/databases/sql"
)

const UserTableName = "pager_users"

func (User) TableName() string {
	return UserTableName
}

type User struct {
	ID        int64     `json:"id" gorm:"primary_key;autoIncrement:true;column:id"`
	Username  string    `json:"username" gorm:"column:username;size:255;not null;unique;index:idx_users_username"`
	Password  string    `json:"password" gorm:"column:password;size:255;not null"`
	Name      string    `json:"name" gorm:"column:name;size:255;not null"`
	EmailID   string    `json:"email_id" gorm:"column:email_id;size:255;not null;unique;index:idx_users_email"`
	UserType  string    `json:"user_type" gorm:"column:user_type;size:50;not null"` // Admin, User
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func CreateUser(ctx context.Context, tx interface{}, username, password, userType string) (*User, error) {
	db := sql.GetOrmQuearyable(ctx, tx)
	user := User{
		Username: username,
		Password: password,
		UserType: userType,
	}
	err := db.Create(&user).Error
	return &user, err
}

func GetUserByUsername(ctx context.Context, tx interface{}, username string) (*User, error) {
	var user User
	db := sql.GetOrmQuearyable(ctx, tx)
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func GetUserByID(ctx context.Context, tx interface{}, id string) (*User, error) {
	var user User
	db := sql.GetOrmQuearyable(ctx, tx)
	err := db.First(&user, id).Error
	return &user, err
}

func (user *User) Save(ctx context.Context, tx interface{}) error {
	db := sql.GetOrmQuearyable(ctx, tx)
	return db.Save(user).Error
}
