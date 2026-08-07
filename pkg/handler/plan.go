package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary      Create plans
// @Security     ApiKeyAuth
// @Description  Create one or more plans for a specific month and year
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        month  query  int                        true  "Month"
// @Param        year   query  int                        true  "Year"
// @Param        input  body   []models.CreatePlan        true  "List of plans to create"
// @Success      200    {object}  statusResponse          "Successfully created"
// @Failure      400    {object}  errorResponse           ""
// @Failure      500    {object}  errorResponse           ""
// @Router       /api/plans [post]
func (h *Handler) createPlan(c *gin.Context) {
	var input []models.CreatePlan

	if err := c.BindJSON(&input); err != nil {
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

	err = h.services.Plan.Create(input, month, year)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})
}

// @Summary      Get plans
// @Security     ApiKeyAuth
// @Description  Retrieve plans for a specific month and year
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        language  header  string  true   "Language"
// @Param        month     query   int     true   "Month"
// @Param        year      query   int     true   "Year"
// @Success      200  {array}  []models.Plan             "List of plans"
// @Failure      400  {object}  errorResponse            ""
// @Failure      500  {object}  errorResponse            ""
// @Router       /api/plans [get]
func (h *Handler) getPlans(c *gin.Context) {
	language := c.GetHeader("language")

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

	//agent_id, err := strconv.Atoi(c.Query("agent_id"))
	//if err != nil {
	//	agent_id = 0
	//}

	result, err := h.services.Plan.GetPlans(month, year, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get main plan statistics
// @Security     ApiKeyAuth
// @Description  Retrieve main dashboard statistics for plans
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        language   header  string  true  "Language"
// @Param        month      query   int     true  "Month"
// @Param        year       query   int     true  "Year"
// @Param        agent_id   query   int     false "Agent ID (optional)"
// @Success      200  {object}  models.PlanStatisticsMain  "Main plan statistics"
// @Failure      400  {object}  errorResponse              ""
// @Failure      500  {object}  errorResponse              ""
// @Router       /api/plans/main [get]
func (h *Handler) getPlanStatisticsMain(c *gin.Context) {
	language := c.GetHeader("language")

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

	agent_id, err := strconv.Atoi(c.Query("agent_id"))
	if err != nil {
		agent_id = 0
	}

	result, err := h.services.Plan.GetPlanStatisticsForMain(month, year, agent_id, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get plan statistics for agent
// @Security     ApiKeyAuth
// @Description  Retrieve plan statistics for a specific agent
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        language   header  string  true  "Language"
// @Param        month      query   int     true  "Month"
// @Param        year       query   int     true  "Year"
// @Success      200  {array}  []models.PlanStatisticsForAgent  "Statistics for agent"
// @Failure      400  {object}  errorResponse               ""
// @Failure      500  {object}  errorResponse               ""
// @Router       /api/plans/agent [get]
func (h *Handler) getPlanStatisticsForAgent(c *gin.Context) {
	language := c.GetHeader("language")

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

	result, err := h.services.Plan.GetPlanStatisticsForMain(month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	for indx, result_item := range result {
		for _, list_item := range result_item.Data {
			result[indx].Count += list_item.Count
		}
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Update plans
// @Security     ApiKeyAuth
// @Description  Update one or more plans
// @Tags         plans
// @Accept       json
// @Produce      json
// @Param        month  query  int                        true  "Month"
// @Param        year   query  int                        true  "Year"
// @Param        input  body  []models.UpdatedPlan  true  "List of plans to update"
// @Success      200  {object}  statusResponse      "Successfully updated"
// @Failure      400  {object}  errorResponse       ""
// @Failure      500  {object}  errorResponse       ""
// @Router       /api/plans [put]
func (h *Handler) updatePlan(c *gin.Context) {
	var input []models.UpdatedPlan

	if err := c.BindJSON(&input); err != nil {
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

	err = h.services.Plan.Update(month, year, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})
}
