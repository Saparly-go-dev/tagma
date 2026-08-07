package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Get Log types
// @Security ApiKeyAuth
// @Tags	logs
// @Description	get log types
// @ID				get-log-types
// @Accept			json
// @Param        language  header    string  false  "Language (optional)"
// @Produce		json
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/logs/types [get]
func (h *Handler) getLogTypes(c *gin.Context) {
	language := c.GetHeader("language")

	result, err := h.services.Log.GetLogTypes(language)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

}

// @Summary      Get All Logs
// @Security     ApiKeyAuth
// @Tags         logs
// @Description  Get all logs with pagination.
// @ID          get-all-logs
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of Log per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        logType       query    string  false  "Logname to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.TradeAgentPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/logs [get]
func (h *Handler) getAllLogs(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	pageNumber, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	logtype := c.Query("logType")

	page, err := h.services.Log.GetAll(language, pageSize, pageNumber, logtype)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}
