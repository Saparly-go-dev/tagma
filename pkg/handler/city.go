package handler

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create City
// @Security ApiKeyAuth
// @Tags			city
// @Description	create city
// @ID				create-city
// @Accept			json
// @Produce		json
// @Param			input	body	tagma.CreateCity	true	"city info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/cities [post]
func (h *Handler) createCity(c *gin.Context) {
	var input tagma.CreateCity

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.City.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Cities
// @Security     ApiKeyAuth
// @Tags         city
// @Description  Get all cities with pagination.
// @ID          get-all-cities
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of cities per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "City name to filter by (optional)"
// @Success      200        {object} tagma.CityPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/cities [get]
func (h *Handler) getAllCities(c *gin.Context) {
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

	page, err := h.services.City.GetAll(pageSize, pageNumber, name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change City status
// @Security     ApiKeyAuth
// @Tags         city
// @Description  Get a city by ID.
// @ID          get-city-by-id
// @Accept       json
// @Produce      json
// @Param        cityId  path  int  true  "City ID"
// @Success      200    {object} int "Success"  "City status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "City not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/cities/{cityId}/lock [get]
func (h *Handler) statusCity(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.City.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get city information by id
// @Security     ApiKeyAuth
// @Tags         city
// @Description  Get city by id
// @ID          get-city-id
// @Accept       json
// @Produce      json
// @Param cityId path int true "City ID"
// @Success      200    {object}  int "Success"  "Get city information by Id"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "City not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/cities/{cityId} [get]
func (h *Handler) getCityById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	city, err := h.services.City.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, city)
}

// @Summary      Update City
// @Security     ApiKeyAuth
// @Tags         city
// @Description  Update a city.
// @ID          update-city
// @Accept       json
// @Produce      json
// @Param        cityId  path  int  true  "City ID"
// @Param        input  body   tagma.CreateCity  true  "City information"
// @Success      200    {object}  int "Success" "City updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "City not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/cities/{cityId} [put]
func (h *Handler) updateCity(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input tagma.CreateCity
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.City.Update(cityId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete City
// @Security     ApiKeyAuth
// @Tags         city
// @Description  Delete a city.
// @ID          delete-city
// @Accept       json
// @Produce      json
// @Param        cityId  path  int  true  "City ID"
// @Success      200    {object}  int "Success"  "City deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "City not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/cities/{cityId} [delete]
func (h *Handler) deleteCity(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.City.Delete(cityId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

func (h *Handler) test(c *gin.Context) {

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
