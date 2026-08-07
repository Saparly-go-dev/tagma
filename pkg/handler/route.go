package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create Route
// @Security ApiKeyAuth
// @Tags			routes
// @Description	create route
// @ID				create-route
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateRoute	true	"route info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/routes [post]
func (h *Handler) createRoute(c *gin.Context) {
	var input models.CreateRoute

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Route.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get all routes
// @Security ApiKeyAuth
// @Description Retrieves all routes for a specific point and language.
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       pageSize  query    int     true      "page element count"
// @Param       pageNumber  query    int     true      "number of page"
// @Param       agentId  query    int     false      "Agent Id"
// @Param       cityId  query    int     false      "City Id"
// @Param       dayId  query    int     false      "Day Id"
// @Param       language  header   string  false     "Language (optional)"
// @Success    200       {object}  interface{}  "List of routes"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes [get]
func (h *Handler) getAllRoutes(c *gin.Context) {
	language := c.GetHeader("language")

	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	pageNumber, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	dayId, err := strconv.Atoi(c.Query("dayId"))
	if err != nil {
		dayId = 0
	}

	page, err := h.services.Route.GetAll(pageSize, pageNumber, agentId, cityId, dayId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary Change route status
// @Security ApiKeyAuth
// @Description Change route status.
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Route ID"
// @Param       language  header   string  false     "Language (optional)"
// @Success    200       {object}  interface{}  "Route details"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/{id}/lock [get]
func (h *Handler) statusRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Route.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get route by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific route by its ID.
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Route ID"
// @Param       language  header   string  false     "Language (optional)"
// @Success    200       {object}  interface{}  "Route details"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/{id} [get]
func (h *Handler) getRouteById(c *gin.Context) {
	language := c.GetHeader("language")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	city, err := h.services.Route.GetById(id, language)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, city)
}

// @Summary Update a route
// @Security ApiKeyAuth
// @Description Updates an existing route with provided information.
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Route ID"
// @Param       input    body     models.CreateRoute  true  "Route update data"
// @Success    200       {object}  statusResponse   "Update status"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/{id} [put]
func (h *Handler) updateRoute(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateRoute
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Route.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a route
// @Security ApiKeyAuth
// @Description Deletes a route by its ID.
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Route ID"
// @Success    200       {object}  statusResponse   "Deletion status"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/{id} [delete]
func (h *Handler) deleteRoute(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Route.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Days for Route
// @Security     ApiKeyAuth
// @Tags         routes
// @Description  Get days.
// @ID          get-days
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Day name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/routes/days [get]
func (h *Handler) GetDaysForRoute(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.Route.GetDays(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary Delete a route
// @Security ApiKeyAuth
// @Description Deletes a route by its ID.
// @ID          ekspeditors-for-route
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       agentId       path     int     true      "Trade Agent ID"
// @Param       language  header    string  false  "Language (optional)"
// @Success    200        {object} []models.SelectObject "result"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/ekspeditors/{agentId} [get]
func (h *Handler) GETEkspeditorsForRoute(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	agentId, err := strconv.Atoi(c.Param("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.services.Route.GetEkspeditors(agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade points for Route
// @Security     ApiKeyAuth
// @Tags         routes
// @Description  Get points.
// @ID          get-points
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Trade point name to filter by (optional)"
// @Param        cityId		query	 int 	false 	"City Id"
// @Param        districtId		query	 int	 false	 "District Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/routes/points [get]
func (h *Handler) GetTradePointsForRoute(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	result, err := h.services.Route.GetPoints(cityId, districtId, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Get all route Orders
// @Security ApiKeyAuth
// @Description Retrieves all route orders for a specific point and language.
// @Tags        routes
// @Accept      json
// @Produce     json
// @Param       agentId  query    int     false      "Agent Id"
// @Param       dayId  query    int     false      "Day Id"
// @Param       language  header   string  false     "Language (optional)"
// @Success    200       {object}  []models.SelectObject  "List of Route Orders"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/orders/get [get]
func (h *Handler) get_route_orders(c *gin.Context) {
	language := c.GetHeader("language")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	dayId, err := strconv.Atoi(c.Query("dayId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	page, err := h.services.Route.GetRouteOrder(agentId, dayId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary		Post route Orders
// @Security ApiKeyAuth
// @Tags			routes
// @Description	post route orders
// @Accept			json
// @Produce		json
// @Param       agentId  query    int     false      "Agent Id"
// @Param       dayId  query    int     false      "Day Id"
// @Param			input	body	[]models.SelectObject	true	"route info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/routes/orders/post [post]
func (h *Handler) post_route_order(c *gin.Context) {
	var input []models.SelectObject

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	dayId, err := strconv.Atoi(c.Query("dayId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err_result := h.services.Route.PostRouteOrder(agentId, dayId, input); err_result != nil {
		newErrorResponse(c, http.StatusInternalServerError, err_result.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Trade points for Route in excel file
// @Security     ApiKeyAuth
// @Tags         routes
// @Description  Get points Excel file.
// @ID          route-points excel
// @Accept       json
// @Produce      json
// @Param       agentId  query    int     false      "Agent Id"
// @Param       dayId  query    int     false      "Day Id"
// @Param       language  header   string  false     "Language (optional)"
// @Success    200       {object}  []models.SelectObject  "List of Route Orders"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/routes/excel [get]
func (h *Handler) excelRouteOrders(c *gin.Context) {
	language := c.GetHeader("language")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	dayId, err := strconv.Atoi(c.Query("dayId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.Route.ExcelRouteOrder(agentId, dayId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=route_orders.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)

}
