package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create a new merchandiser
// @Security ApiKeyAuth
// @Description Creates a new merchandiser in the system.
// @Tags merchandisers
// @Accept json
// @Produce json
// @Param input body models.CreateMerchandiser true "Merchandiser information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/merchandisers [post]
func (h *Handler) createMerchandiser(c *gin.Context) {
	var input models.CreateMerchandiser

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Merchandiser.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Merchandisers
// @Security     ApiKeyAuth
// @Tags         merchandisers
// @Description  Get all merchandisers with pagination.
// @ID          get-all-merchandisers
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        pageSize  query    int     false  "Number of merchandisers per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "Merchandiser name to filter by (optional)"
// @Success      200        {object} models.MerchandiserPage "Paginated list of merchandisers"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/merchandisers [get]
func (h *Handler) getAllMerchandisers(c *gin.Context) {
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

	name := c.Query("name")

	page, err := h.services.Merchandiser.GetAll(pageNumber, pageSize, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Merchandiser status
// @Security     ApiKeyAuth
// @Tags         merchandisers
// @Description  Change Merchandiser status.
// @ID          change-merchandiser-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Merchandiser id"
// @Success      200    {object} int "Success"  "Merchandiser status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "Merchandiser not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/merchandisers/{id}/lock [get]
func (h *Handler) statusMerchandiser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Merchandiser.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get merchandiser by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific merchandiser by its ID.
// @Tags merchandisers
// @Produce json
// @Param id path int true "Merchandiser ID"
// @Success 200 {object} models.Merchandiser "Merchandiser details"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/merchandisers/{id} [get]
func (h *Handler) getMerchandiserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	merchandiser, err := h.services.Merchandiser.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, merchandiser)
}

// @Summary      Update Merchandiser information
// @Security     ApiKeyAuth
// @Tags         merchandisers
// @Description  Update merchandiser.
// @ID          update-merchandiser
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Merchandiser id"
// @Param        input  body   models.CreateMerchandiser  true  "Merchandiser information"
// @Success      200    {object}  int "Success" "Merchandiser updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Merchandiser not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/merchandisers/{id} [put]
func (h *Handler) updateMerchandiser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateMerchandiser
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Merchandiser.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a merchandiser
// @Security ApiKeyAuth
// @Description Deletes a merchandiser from the system.
// @Tags merchandisers
// @Produce json
// @Param id path int true "Merchandiser ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/merchandisers/{id} [delete]
func (h *Handler) deleteMerchandiser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Merchandiser.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get merchandisers for select
// @Security ApiKeyAuth
// @Description  Returns a list of merchandisers formatted for select inputs, with language support and optional filtering by name.
// @Tags         merchandisers
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language preference (e.g., 'ru', 'tm')"
// @Param        name      query     string  false  "Filter by merchandiser name"
// @Success      200       {array}   []models.SelectObject  "List of merchandisers for select"
// @Failure      500       {object}  errorResponse          "Internal server error"
// @Router       /api/merchandisers/select [get]
func (h *Handler) getMerchandiserForSelect(c *gin.Context) {
	language := c.GetHeader("language")
	name := c.Query("name")

	result, err := h.services.Merchandiser.GetMerchandiserForSelect(language, name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Points For marketing
// @Security     ApiKeyAuth
// @Tags         merchandisers
// @Description  Get trade points.
// @ID           get-trade-points-for-merchandiser
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.TradePointMobile  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/merchandisers/marketing [get]
func (h *Handler) getTradePointsListForMerchandiser(c *gin.Context) {
	language := c.GetHeader("language")

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	day, err := strconv.Atoi(c.Query("day"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Merchandiser.GetMarketingPageForMobile(day, month, year, userId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
