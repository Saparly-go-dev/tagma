package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create a new driver
// @Security ApiKeyAuth
// @Description Creates a new driver in the system.
// @Tags drivers
// @Accept json
// @Produce json
// @Param input body models.CreateDriver true "Driver information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/drivers [post]
func (h *Handler) createDriver(c *gin.Context) {
	var input models.CreateDriver

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Driver.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Drivers
// @Security     ApiKeyAuth
// @Tags         drivers
// @Description  Get all drivers with pagination.
// @ID          get-all-drivers
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        pageSize  query    int     false  "Number of drivers per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "Driver name to filter by (optional)"
// @Success      200        {object} models.DriverPage "Paginated list of drivers"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/drivers [get]
func (h *Handler) getAllDrivers(c *gin.Context) {
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

	page, err := h.services.Driver.GetAll(pageNumber, pageSize, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Driver status
// @Security     ApiKeyAuth
// @Tags         drivers
// @Description  Change Driver status.
// @ID          change-driver-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Driver id"
// @Success      200    {object} int "Success"  "Driver status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "Driver not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/drivers/{id}/lock [get]
func (h *Handler) statusDriver(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Driver.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get driver by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific driver by its ID.
// @Tags drivers
// @Produce json
// @Param id path int true "Driver ID"
// @Success 200 {object} models.Driver "Driver details"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/drivers/{id} [get]
func (h *Handler) getDriverById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	driver, err := h.services.Driver.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, driver)
}

// @Summary      Update Driver information
// @Security     ApiKeyAuth
// @Tags         drivers
// @Description  Update driver.
// @ID          update-driver
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Driver id"
// @Param        input  body   models.CreateDriver  true  "Driver information"
// @Success      200    {object}  int "Success" "Driver updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Driver not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/drivers/{id} [put]
func (h *Handler) updateDriver(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateDriver
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Driver.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a driver
// @Security ApiKeyAuth
// @Description Deletes a driver from the system.
// @Tags drivers
// @Produce json
// @Param id path int true "Driver ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/drivers/{id} [delete]
func (h *Handler) deleteDriver(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Driver.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
