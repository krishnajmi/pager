package templates

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateController struct {
	templateService TemplateService
}

func NewTemplateController(templateService TemplateService) *TemplateController {
	return &TemplateController{
		templateService: templateService,
	}
}

func (c *TemplateController) CreateTemplate(ctx *gin.Context) {
	var request TemplateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("createTemplateView:unableToBindJSON", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := c.templateService.CreateTemplate(ctx.Request.Context(), request.Name, request.Subject, request.Content)
	if err != nil {
		slog.Error("createTemplateView:unableToCreateTemplate", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": false,
			"msg":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"error":  nil,
		"status": true,
		"msg":    "Template created successfully",
		"data":   template,
	})
}

func (c *TemplateController) UpdateTemplate(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var request TemplateRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		slog.Error("updateTemplateView:unableToBindJSON", slog.Any("error", err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template, err := c.templateService.UpdateTemplate(ctx.Request.Context(), id, request.Name, request.Subject, request.Content)
	if err != nil {
		slog.Error("updateTemplateView:unableToUpdateTemplate",
			slog.Int64("template_id", id),
			slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": false,
			"msg":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":  nil,
		"status": true,
		"msg":    "Template updated successfully",
		"data":   template,
	})
}

func (c *TemplateController) GetTemplate(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	template, err := c.templateService.GetTemplate(ctx.Request.Context(), id)
	if err != nil {
		slog.Error("getTemplateView:unableToGetTemplate",
			slog.Int64("template_id", id),
			slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": false,
			"msg":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":  nil,
		"status": true,
		"msg":    "Template retrieved successfully",
		"data":   template,
	})
}

func (c *TemplateController) GetAllTemplates(ctx *gin.Context) {
	templates, err := c.templateService.GetAllTemplates(ctx.Request.Context())
	if err != nil {
		slog.Error("getAllTemplatesView:unableToGetTemplates", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  err.Error(),
			"status": false,
			"msg":    err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"error":  nil,
		"status": true,
		"msg":    "Templates retrieved successfully",
		"data":   templates,
	})
}
