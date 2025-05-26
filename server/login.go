package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	login "github.com/kp/pager/login"
)

func AuthRouterGroup(servicePrefix string, db *gorm.DB, middlewares ...gin.HandlerFunc) RouterGroup {
	return RouterGroup{
		Prefix:      servicePrefix,
		Routes:      authRoutes(db, servicePrefix),
		Middlewares: middlewares}
}

func authRoutes(db *gorm.DB, prefix string) []Route {
	// Initialize repositories
	userRepo := login.NewUserRepository(db)
	permRepo := login.NewPermissionRepository(db)
	userPermRepo := login.NewUserPermissionRepository(db)

	// Initialize services
	authService := login.NewAuthService(userRepo, permRepo, userPermRepo)

	// Initialize controllers
	authCtrl := login.NewAuthController(authService)

	return []Route{
		newRoute(http.MethodPost, "/login/", authCtrl.Login, prefix),
		newRoute(http.MethodPost, "/register/", authCtrl.Register, prefix, login.PagerAdminAccess),
		newRoute(http.MethodPost, "/permissions/", authCtrl.AddPermission, prefix, login.PagerAdminAccess),
		newRoute(http.MethodGet, "/permissions/user/:user_id/", authCtrl.GetPermissions, prefix, login.PagerAdminAccess),
		newRoute(http.MethodGet, "/permissions/", authCtrl.GetAllPermissions, prefix, login.PagerAdminAccess),
		newRoute(http.MethodGet, "/users/:user_id/", authCtrl.GetUserDetails, prefix, login.PagerAdminAccess),
		newRoute(http.MethodPost, "/permissions/add/", authCtrl.AddUserPermission, prefix, login.PagerAdminAccess),
		newRoute(http.MethodPost, "/permissions/remove/", authCtrl.RemoveUserPermission, prefix, login.PagerAdminAccess),
		newRoute(http.MethodPost, "/reset/password", authCtrl.ChangePassword, prefix, login.PagerAdminAccess),
		newRoute(http.MethodGet, "/", authCtrl.GetAllUsers, prefix, login.PagerAdminAccess),
	}
}
