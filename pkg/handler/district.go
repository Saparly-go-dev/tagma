package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create a new district
// @Security ApiKeyAuth
// @Description Creates a new district in the system.
// @Tags districts
// @Accept json
// @Produce json
// @Param input body models.CreateDistrict true "District information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/districts [post]
func (h *Handler) createDistrict(c *gin.Context) {
	var input models.CreateDistrict

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.District.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Districts
// @Security     ApiKeyAuth
// @Tags         districts
// @Description  Get all districts with pagination.
// @ID          get-all-districts
// @Accept       json
// @Produce      json
// @Param 		 language header string false "Language (optional)"
// @Param        pageSize  query    int     false  "Number of districts per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "District name to filter by (optional)"
// @Success      200        {object} models.DistrictPage "Paginated list of districts"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/districts [get]
func (h *Handler) getAllDistricts(c *gin.Context) {
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

	page, err := h.services.District.GetPage(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change District status
// @Security     ApiKeyAuth
// @Tags         districts
// @Description  Change District status.
// @ID          change-district-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "District id"
// @Success      200    {object} int "Success"  "District status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "District not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/districts/{id}/lock [get]
func (h *Handler) statusDistrict(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.District.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get district by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific district by its ID.
// @Tags districts
// @Produce json
// @Param id path int true "District ID"
// @Success 200 {object} models.District "District details"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/districts/{id} [get]
func (h *Handler) getDistrictById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	district, err := h.services.District.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, district)
}

// @Summary      Update District information
// @Security     ApiKeyAuth
// @Tags         districts
// @Description  Update district.
// @ID          update-district
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "District id"
// @Param        input  body   models.CreateDistrict  true  "District information"
// @Success      200    {object}  int "Success" "District updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "District not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/districts/{id} [put]
func (h *Handler) updateDistrict(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateDistrict
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.District.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a district
// @Security ApiKeyAuth
// @Description Deletes a district from the system.
// @Tags districts
// @Produce json
// @Param id path int true "District ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/districts/{id} [delete]
func (h *Handler) deleteDistrict(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.District.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
