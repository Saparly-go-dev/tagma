package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary		Create TradeAgent
// @Security ApiKeyAuth
// @Tags	agent
// @Description	create tradeAgent
// @ID				create-tradeAgent
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateTradeAgent	true	"tradeAgent info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/agents [post]
func (h *Handler) createTradeAgent(c *gin.Context) {
	var input models.CreateTradeAgent

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TradeAgent.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

// @Summary      Get All Trade Agents
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get all trade agents with pagination.
// @ID          get-all-data
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of Trade Agent per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "TradeAgent name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.TradeAgentPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/agents [get]
func (h *Handler) getAllTradeAgents(c *gin.Context) {
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

	page, err := h.services.TradeAgent.GetAll(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change TradeAgent status
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Change a tradeAgent status.
// @ID          change_agent_status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradeAgent ID"
// @Success      200    {object} int "Success"  "TradeAgent status updated"
// @Failure      400    {object}  errorResponse "Invalid tradeAgent ID"
// @Failure      404    {object}  errorResponse "TradeAgent not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/agents/{id}/lock [get]
func (h *Handler) statusTradeAgent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TradeAgent.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get tradeAgent information by id
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get tradeAgent by id
// @ID          get-tradeAgent-id
// @Accept       json
// @Produce      json
// @Param id path int true "TradeAgent ID"
// @Success      200    {object}  int "Success"  "Get tradeAgent information by id"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "TradeAgent not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/agents/{id} [get]
func (h *Handler) getTradeAgentById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	tradeAgent, err := h.services.TradeAgent.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tradeAgent)
}

// @Summary      Update TradeAgent
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Update a tradeAgent.
// @ID          update-tradeAgent
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradeAgent ID"
// @Param        input  body   models.CreateTradeAgent  true  "TradeAgent information"
// @Success      200    {object}  int "Success" "TradeAgent updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "TradeAgent not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/agents/{id} [put]
func (h *Handler) updateTradeAgent(c *gin.Context) {
	tradeAgentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateTradeAgent
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.TradeAgent.Update(tradeAgentId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete TradeAgent
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Delete a tradeAgent.
// @ID          delete-tradeAgent
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradeAgent ID"
// @Success      200    {object}  int "Success"  "TradeAgent deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "TradeAgent not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/agents/{id} [delete]
func (h *Handler) deleteTradeAgent(c *gin.Context) {
	tradeAgentId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.TradeAgent.Delete(tradeAgentId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Driver for Agents
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get drivers.
// @ID          get-drives
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Driver name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/agents/drivers [get]
func (h *Handler) getDriverforAgents(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.TradeAgent.GetDrivers(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get All Ekspeditors for Agents
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get ekspeditors.
// @ID          get-ekspeditors
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Driver name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/agents/ekspeditors [get]
func (h *Handler) getEkspeditorforAgents(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.TradeAgent.GetEkspeditors(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get All Trade Agents Excel file
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get all trade agents with pagination.
// @ID          get-all-data-excel
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "TradeAgent name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.TradeAgentPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/agents/excel [get]
func (h *Handler) GetAgentExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	result, err := h.services.TradeAgent.AgentExcel(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=agents.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// @Summary      Get  Trade Agent debt
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get trade agent debt.
// @ID          get-agent_debt
// @Accept       json
// @Produce      json
// @Param        agentId  query    int     false  "Number of Trade Agent per page"  Minimum: 1
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.AgentDebt "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/agents/debt [get]
func (h *Handler) getAgentDebt(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	page, err := h.services.TradeAgent.GetAgentDebt(agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get trade agent revenue for month
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Возвращает сумму подтверждённых оплат агента за месяц.
// @ID           get-agent-revenue
// @Accept       json
// @Produce      json
// @Param        month    query   int  true  "Месяц (1-12)"
// @Param        year     query   int  true  "Год"
// @Param        language header string false "Язык"
// @Success      200 {object} map[string]int "revenue"
// @Failure      400 {object} errorResponse
// @Failure      500 {object} errorResponse
// @Router       /api/agents/revenue [get]
func (h *Handler) getTradeAgentRevenue(c *gin.Context) {
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		newErrorResponse(c, http.StatusBadRequest, "invalid month")
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 1 {
		newErrorResponse(c, http.StatusBadRequest, "invalid year")
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	language := c.GetHeader("language")

	sum, err := h.services.TradeAgent.GetTradeAgentRevenue(month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"revenue": sum})
}

// @Summary      Get factory CSV for agent/day
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Возвращает CSV со сводным списком продуктов для фабрики за указанную дату.
// @ID           get-factory-csv
// @Accept       json
// @Produce      text/csv
// @Param        day     query    int  true  "День"
// @Param        month   query    int  true  "Месяц"
// @Param        year    query    int  true  "Год"
// @Param        language header  string false "Язык"
// @Success      200 {file} file "CSV файл"
// @Failure      400 {object} errorResponse "Invalid request parameters"
// @Failure      500 {object} errorResponse "Internal server error"
// @Router       /api/agents/factory/csv [get]
func (h *Handler) getFactoryCSV(c *gin.Context) {
	day, err := strconv.Atoi(c.Query("day"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid day")
		return
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid month")
		return
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid year")
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	language := c.GetHeader("language")

	result, err := h.services.TradeAgent.CreateCSVFileForFactory(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	filename := fmt.Sprintf("%02d-%02d-%04d.csv", day, month, year)

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	c.Data(http.StatusOK, "text/csv", result)
}

// @Summary      Get Trade Points Merchandising
// @Security     ApiKeyAuth
// @Tags         agent
// @Description  Get trade points merchendising.
// @ID          get-trade-points-merchendising
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   models.PointsMerchendising_List{}  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/agents/merchendising [get]
func (h *Handler) GetMerchandiserInformationForAgent(c *gin.Context) {
	language := c.GetHeader("language")

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

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.TradeAgent.GetMerchendisingInformation(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
