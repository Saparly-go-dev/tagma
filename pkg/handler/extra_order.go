package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary      Get All Trade Points For Extra Order
// @Security     ApiKeyAuth
// @Tags         extraorder
// @Description  Get all trade points for extra order.
// @ID          get-points-extra_order
// @Accept       json
// @Produce      json
// @Param        day query    int     false  "Day Id" Minimum: 1
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.TradePointMobile "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/extra/order/agent [get]
func (h *Handler) getTradePointsForExtraOrder(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

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

	day_id, err := strconv.Atoi(c.Query("day"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.ExtraOrder.GetTradePoint(day_id, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Order
// @Security     ApiKeyAuth
// @Tags         extraorder
// @Description  Get Orders
// @ID          get-order
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        day       query    int  false  "day"
// @Param        month       query    int  false  "month"
// @Param        year       query    int  false  "year"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Success      200    {object}  models.GetOrderObject "Success"  "Orders List"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/extra/order/items [get]
func (h *Handler) getOrderItemsForExtraOrder(c *gin.Context) {
	// Get the 'language' header
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

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.ExtraOrder.GetOrder(day, month, year, tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
