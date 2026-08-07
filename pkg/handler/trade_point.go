package handler

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create a trade point
// @Description Creates a new trade point with provided data and uploads images (requires API key).
// @Security ApiKeyAuth
// @Tags        point
// @Accept      multipart/form-data
// @Produce     json
// @Param       input  formData  models.CreateTradePoint   true      "JSON data for trade point (in 'input' field)"
// @Param       image1 formData  file      false     "Image 1 (optional)"
// @Param       image2 formData  file      false     "Image 2 (optional)"
// @Success    200     {object}   map[string]interface{}  "Success message with ID"
// @Failure    400     {object}   errorResponse          "Bad request (invalid data or image upload error)"
// @Failure    500     {object}   errorResponse          "Internal server error"
// @Router      /api/points [post]
func (h *Handler) createTradePoint(c *gin.Context) {
	// Retrieve JSON data from form input (assuming 'input' field contains JSON string)
	var input models.CreateTradePoint
	cityId, err := strconv.Atoi(c.PostForm("city_id"))
	nameRu := c.PostForm("name_ru")
	nameTm := c.PostForm("name_tm")
	locationRu := c.PostForm("location_ru")
	locationTm := c.PostForm("location_tm")
	orientirRu := c.PostForm("orientir_ru")
	orientirTm := c.PostForm("orientir_tm")
	status := c.PostForm("status")
	districtId, err := strconv.Atoi(c.PostForm("district_id"))
	tradeAgentId, err := strconv.Atoi(c.PostForm("trade_agent_id"))
	tradeChannelId, err := strconv.Atoi(c.PostForm("trade_channel_id"))
	dayId, err := strconv.Atoi(c.PostForm("day_id"))
	ekspeditorId, err := strconv.Atoi(c.PostForm("ekspeditor_id"))

	input.CityId = cityId
	input.Code = "10001"
	input.NameRu = nameRu
	input.NameTm = nameTm
	input.LocationRu = locationRu
	input.LocationTm = locationTm
	input.OrientirRu = orientirRu
	input.OrientirTm = orientirTm
	input.TradeAgentId = tradeAgentId
	input.TradeChannelId = tradeChannelId
	input.DistrictId = districtId
	input.DayId = dayId
	input.EkspeditorId = ekspeditorId
	if status == "false" {
		input.Status = false
	} else {
		input.Status = true
	}

	// Handle image1 upload
	image1, err := c.FormFile("image1")
	if err != nil {
		image1 = nil
	}

	// Handle image2 upload
	image2, err := c.FormFile("image2")
	if err != nil {
		image2 = nil
	}

	// Call service to create trade point with input and images
	id, err := h.services.TradePoint.Create(input, image1, image2)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with success and return the ID
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Trade Points on Excel File
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get all trade points on excel file.
// @ID          get-points-excel
// @Accept       json
// @Produce      json
// @Param       name       query    string  false  "TradePoint name to filter by (optional)"
// @Param       agentId  query    int     false      "Agent Id"
// @Param       cityId  query    int     false      "City Id"
// @Param       districtId  query    int     false      "District Id"
// @Param       dayId  query    int     false      "Day Id"
// @Param       ekspeditorId  query    int     false      "Ekspeditor Id"
// @Param       typeId  query    int     false      "Trade Channel Type Id"
// @Param       structureId  query    int     false      "Trade Channel Structure Id"
// @Param       sizeId  query    int     false      "Trade Channel Size Id"
// @Param       managementId  query    int     false      "Trade Channel Management Id"
// @Param       language  header    string  false  "Language (optional)"
// @Success      200        {object} models.TradePointPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/excel [get]
func (h *Handler) getTradePointExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	dayId, err := strconv.Atoi(c.Query("dayId"))
	if err != nil {
		dayId = 0
	}

	ekspeditorId, err := strconv.Atoi(c.Query("ekspeditorId"))
	if err != nil {
		ekspeditorId = 0
	}

	typeId, err := strconv.Atoi(c.Query("typeId"))
	if err != nil {
		typeId = 0
	}

	structureId, err := strconv.Atoi(c.Query("structureId"))
	if err != nil {
		structureId = 0
	}

	sizeId, err := strconv.Atoi(c.Query("sizeId"))
	if err != nil {
		sizeId = 0
	}

	managementId, err := strconv.Atoi(c.Query("managementId"))
	if err != nil {
		managementId = 0
	}

	result, err := h.services.TradePoint.GetExcel(cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=points.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// @Summary      Get All Trade Points
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get all trade points with pagination.
// @ID          get-all-points
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of Trade Point per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "TradePoint name to filter by (optional)"
// @Param        contact       query    string  false  "TradePoint contact number to filter by (optional)"
// @Param        code       query    string  false  "TradePoint code to filter"
// @Param        agentId  query    int     false      "Agent Id"
// @Param        cityId  query    int     false      "City Id"
// @Param        districtId  query    int     false      "District Id"
// @Param        dayId  query    int     false      "Day Id"
// @Param        ekspeditorId  query    int     false      "Ekspeditor Id"
// @Param        typeId  query    int     false      "Trade Channel Type Id"
// @Param        structureId  query    int     false      "Trade Channel Structure Id"
// @Param        sizeId  query    int     false      "Trade Channel Size Id"
// @Param        status  query    int     false      "Trade point Status"
// @Param        managementId  query    int     false      "Trade Channel Management Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.TradePointPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points [get]
func (h *Handler) getAllTradePoints(c *gin.Context) {
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

	name := c.Query("name")
	code := c.Query("code")
	contact := c.Query("contact")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	dayId, err := strconv.Atoi(c.Query("dayId"))
	if err != nil {
		dayId = 0
	}

	ekspeditorId, err := strconv.Atoi(c.Query("ekspeditorId"))
	if err != nil {
		ekspeditorId = 0
	}

	typeId, err := strconv.Atoi(c.Query("typeId"))
	if err != nil {
		typeId = 0
	}

	structureId, err := strconv.Atoi(c.Query("structureId"))
	if err != nil {
		structureId = 0
	}

	sizeId, err := strconv.Atoi(c.Query("sizeId"))
	if err != nil {
		sizeId = 0
	}

	managementId, err := strconv.Atoi(c.Query("managementId"))
	if err != nil {
		managementId = 0
	}

	status, err := strconv.Atoi(c.Query("status"))
	if err != nil {
		status = -1
	}

	page, err := h.services.TradePoint.GetAll(pageSize, pageNumber, cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, status, name, code, contact, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change TradePoint status
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Change a tradePoint status.
// @ID          change_point_status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradePoint ID"
// @Success      200    {object} int "Success"  "TradePoint status updated"
// @Failure      400    {object}  errorResponse "Invalid tradePoint ID"
// @Failure      404    {object}  errorResponse "TradePoint not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/points/{id}/lock [get]
func (h *Handler) statusTradePoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TradePoint.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get tradePoint information by id
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get tradePoint by id
// @ID          get-tradePoint-id
// @Accept       json
// @Produce      json
// @Param id path int true "TradePoint ID"
// @Success      200    {object}  int "Success"  "Get tradePoint information by id"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "TradePoint not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/points/{id} [get]
func (h *Handler) getTradePointById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tradePoint, err := h.services.TradePoint.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tradePoint)
}

// @Summary      Update TradePoint
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Update a tradePoint.
// @ID          update-tradePoint
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradePoint ID"
// @Param        input  body   models.CreateTradePoint  true  "TradePoint information"
// @Success      200    {object}  int "Success" "TradePoint updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "TradePoint not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/points/{id} [put]
func (h *Handler) updateTradePoint(c *gin.Context) {
	tradePointId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateTradePoint
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.TradePoint.Update(tradePointId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Update trade points agent
// @Security     ApiKeyAuth
// @Description Assigns multiple trade points to a specific agent
// @Tags point
// @Accept json
// @Produce json
// @Param id query int true "Agent ID to assign trade points to"
// @Param ekspeditor_id query int true "Ekspeditor ID to assign trade points to"
// @Param day_id query int true "day ID to assign trade points to"
// @Param ids body []int true "Array of trade point IDs to assign"
// @Success 200 {object} statusResponse "Successful operation"
// @Failure 400 {object} errorResponse "Invalid input data"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/points/change_agent [put]
func (h *Handler) updateTradePointsAgent(c *gin.Context) {
	agent_id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	ekspeditor_id, err := strconv.Atoi(c.Query("ekspeditor_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	day_id, err := strconv.Atoi(c.Query("day_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var ids []int
	if err = c.BindJSON(&ids); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result := h.services.TradePoint.ChangeTradePointAgent(agent_id, ekspeditor_id, day_id, ids)
	if result != nil {
		newErrorResponse(c, http.StatusInternalServerError, "<UNK>")
		return
	}
	c.JSON(http.StatusOK, statusResponse{Status: "OK"})
}

// @Summary      Delete TradePoint
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Delete a tradePoint.
// @ID          delete-tradePoint
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradePoint ID"
// @Success      200    {object}  int "Success"  "TradePoint deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "TradePoint not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/points/{id} [delete]
func (h *Handler) deleteTradePoint(c *gin.Context) {
	tradePointId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.TradePoint.Delete(tradePointId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get  Cities for TradePoints
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get cities for trade points.
// @ID           get-cities-point
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "City name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/cities [get]
func (h *Handler) getCitiesForPoints(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.TradePoint.GetCities(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get  districts for TradePoints
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get districts for trade points.
// @ID           get-districts-point
// @Accept       json
// @Produce      json
// @Param        cityId       query    int  false  "City Id to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/districts [get]
func (h *Handler) getDistrictsForCity(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.TradePoint.GetDistrictsByCity(cityId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get  Agents for TradePoints
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get agents for trade points.
// @ID           get-agents-point
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Trade Agents name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/agents [get]
func (h *Handler) getAgentsForPoints(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.TradePoint.GetAgents(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get  Channels for TradePoints
// @Security     ApiKeyAuth
// @Tags         point
// @Description  Get channels for trade points.
// @ID           get-channels-point
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/channels [get]
func (h *Handler) getChannelsForTradePoints(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.TradePoint.GetChannels(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get Trade Point information
// @Security     ApiKeyAuth
// @Tags         point
// @Description   Get trade points.
// @ID           get-general-information
// @Accept       json
// @Produce      json
// @Param id path int true "TradePoint ID"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.GeneralTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/{id}/general [get]
func (h *Handler) getGeneralInformationForTradePoint(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.services.TradePoint.GetGeneralInformation(Id, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Images for Point
// @Security     ApiKeyAuth
// @Tags         point
// @Description   Get trade points images.
// @ID           get-point-image
// @Accept       json
// @Produce      json
// @Param id path int true "TradePoint ID"
// @Success      200        {object} models.GeneralTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/images/{id} [get]
func (h *Handler) GetPointImages(c *gin.Context) {
	// Get the 'language' header

	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.services.TradePoint.GetPointImages(Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Delete Image
// @Security     ApiKeyAuth
// @Tags         point
// @Description   Get trade points.
// @ID           delete-image
// @Accept       json
// @Produce      json
// @Param id path int true "TradePoint ID"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/images/{id} [delete]
func (h *Handler) deleteImage(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.TradePoint.DeleteImage(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Save Point Image
// @Description Resave trade point image.
// @ID           resave-image
// @Security ApiKeyAuth
// @Tags        point
// @Accept      multipart/form-data
// @Produce     json
// @Param       point_id  formData  int   true      "Point Id"
// @Param       image formData  file      false     "Image 1 (optional)"
// @Success    200     {object}   map[string]interface{}  "Success message with ID"
// @Failure    400     {object}   errorResponse          "Bad request (invalid data or image upload error)"
// @Failure    500     {object}   errorResponse          "Internal server error"
// @Router      /api/points/resave [post]
func (h *Handler) ResavePointImage(c *gin.Context) {
	// Retrieve JSON data from form input (assuming 'input' field contains JSON string)
	pointId, err := strconv.Atoi(c.PostForm("point_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Handle image1 upload
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Image upload error: %s", err.Error()))
		return
	}

	if err := h.services.TradePoint.ReSaveImage(pointId, file); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Trade Point information about Orders
// @Security     ApiKeyAuth
// @Tags         point
// @Description   Get trade points.
// @ID           get-order-information
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade point Id"
// @Param        startDay       query    int  false "start day"
// @Param        startMonth       query    int  false  "start month"
// @Param        startYear       query    int  false  "start year"
// @Param        endDay       query    int  false  "end day"
// @Param        endMonth       query    int  false  "end month"
// @Param        endYear       query    int  false  "end year"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.GeneralTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/order [get]
func (h *Handler) getInformationAboutOrder(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	startDay, err := strconv.Atoi(c.Query("startDay"))
	if err != nil {
		startDay = 0
	}

	startMonth, err := strconv.Atoi(c.Query("startMonth"))
	if err != nil {
		startMonth = 0
	}

	startYear, err := strconv.Atoi(c.Query("startYear"))
	if err != nil {
		startYear = 0
	}

	endDay, err := strconv.Atoi(c.Query("endDay"))
	if err != nil {
		endDay = 0
	}

	endMonth, err := strconv.Atoi(c.Query("endMonth"))
	if err != nil {
		endMonth = 0
	}

	endYear, err := strconv.Atoi(c.Query("endYear"))
	if err != nil {
		endYear = 0
	}

	result, err := h.services.TradePoint.GetInformationAboutOrder(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Point information About Product
// @Security     ApiKeyAuth
// @Tags         point
// @Description   Get trade points.
// @ID           get-product-information
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade point Id"
// @Param        startDay       query    int  false "start day"
// @Param        startMonth       query    int  false  "start month"
// @Param        startYear       query    int  false  "start year"
// @Param        endDay       query    int  false  "end day"
// @Param        endMonth       query    int  false  "end month"
// @Param        endYear       query    int  false  "end year"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.GeneralTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/product [get]
func (h *Handler) getInformationAboutProducts(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	startDay, err := strconv.Atoi(c.Query("startDay"))
	if err != nil {
		startDay = 0
	}

	startMonth, err := strconv.Atoi(c.Query("startMonth"))
	if err != nil {
		startMonth = 0
	}

	startYear, err := strconv.Atoi(c.Query("startYear"))
	if err != nil {
		startYear = 0
	}

	endDay, err := strconv.Atoi(c.Query("endDay"))
	if err != nil {
		endDay = 0
	}

	endMonth, err := strconv.Atoi(c.Query("endMonth"))
	if err != nil {
		endMonth = 0
	}

	endYear, err := strconv.Atoi(c.Query("endYear"))
	if err != nil {
		endYear = 0
	}

	result, err := h.services.TradePoint.GetInformantionAboutProducts(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Point Sale History
// @Security     ApiKeyAuth
// @Tags         point
// @Description   Get trade points sale history.
// @ID           get-trade-point-sale-history
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.GeneralTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/points/sale/history [get]
func (h *Handler) getTradePointSaleHistoryForMonth(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.TradePoint.GetTradePointSaleHistoryForMonth(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
