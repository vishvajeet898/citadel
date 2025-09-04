package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	commonUtils "github.com/Orange-Health/citadel/common/utils"
)

// @Summary		Get Templates By Type
// @Description	Get Templates By Type
// @Tags			templates
// @Accept			json
// @Produce		json
// @Param			template_type	query		string					false	"Template Type"
// @Success		200				{object}	[]structures.Template	"Templates"
// @Router			/api/v1/templates [get]
func (templateController *Template) GetTemplatesByType(c *gin.Context) {
	templateType := c.Query("template_type")
	templates, err := templateController.TemplateService.GetTemplatesByType(templateType)
	if err != nil {
		commonUtils.HandleErrorResponse(c, err.StatusCode, err.Message)
		return
	}

	c.JSON(http.StatusOK, templates)
}
