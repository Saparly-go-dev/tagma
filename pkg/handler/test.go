package handler

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary		Create test
// @Security ApiKeyAuth
// @Tags			test
// @Description	create test
// @ID				create-test
// @Accept			json
// @Produce		json
// @Param			input	body	tagma.Test	true	"test info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/test [post]
func (h *Handler) createTest(c *gin.Context) {
	var input tagma.Test

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Test.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
