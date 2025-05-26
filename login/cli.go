package login

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	"github.com/kp/pager/login/models"
)

type UserCLI struct {
	DB *gorm.DB
}

func (c *UserCLI) CreateAdmin(username, password string) error {
	// Check if admin exists
	var existing models.User
	if err := c.DB.Where("username = ? AND user_type = ?", username, "Admin").First(&existing).Error; err == nil {
		return fmt.Errorf("admin user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create admin
	admin := models.User{
		Username: username,
		Password: string(hashedPassword),
		UserType: "Admin",
	}

	if err := c.DB.Create(&admin).Error; err != nil {
		return fmt.Errorf("failed to create admin: %v", err)
	}

	fmt.Printf("Admin user '%s' created successfully\n", username)
	return nil
}

func (c *UserCLI) CreateUser(username, password, userType string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	user := models.User{
		Username: username,
		Password: string(hashedPassword),
		UserType: userType,
	}

	if err := c.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	fmt.Printf("User '%s' created successfully\n", username)
	return nil
}
