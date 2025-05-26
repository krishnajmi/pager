package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kp/pager/templates"
)

func TemplateRouterGroup(servicePrefix string, db *gorm.DB, middlewares ...gin.HandlerFunc) RouterGroup {
	return RouterGroup{
		Prefix:      servicePrefix,
		Routes:      templateRoutes(db, servicePrefix),
		Middlewares: middlewares}
}

func templateRoutes(db *gorm.DB, prefix string) []Route {
	// Initialize services
	templateService := templates.NewTemplateService(db)

	// Initialize controllers
	templateCtrl := templates.NewTemplateController(templateService)

	return []Route{
		newRoute(http.MethodPost, "/", templateCtrl.CreateTemplate, prefix),
		newRoute(http.MethodPut, "/:id/", templateCtrl.UpdateTemplate, prefix),
		newRoute(http.MethodGet, "/:id/", templateCtrl.GetTemplate, prefix),
		newRoute(http.MethodGet, "/", templateCtrl.GetAllTemplates, prefix),
	}
}
