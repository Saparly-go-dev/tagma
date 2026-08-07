package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary      Create a trade point debt
// @Security ApiKeyAuth
// @Description  Creates a new debt record for a specific trade point.
// @Tags         trade-point-debts
// @Accept       json
// @Produce      json
// @Param        id       query     int     true   "Trade Point ID"
// @Param        price1       query    int  false  "Price 1"
// @Param        price2       query    int  false  "Price 2"
// @Success      200    {object}  object  "Successfully created debt"
// @Failure      400    {object}  errorResponse         "Invalid Trade Point ID or Debt Amount format"
// @Failure      500    {object}  errorResponse         "Internal server error"
// @Router       /api/points/debts/create [post]
func (h *Handler) CreateTradePointDebt(c *gin.Context) {
	trade_point_id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	price1, err := strconv.Atoi(c.Query("price1"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	price2, err := strconv.Atoi(c.Query("price2"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	currency := float64(price1) + float64(price2)/100

	err = h.services.TradePointDebt.Create(trade_point_id, currency)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"Ok"})
}

// @Summary      Get trade point debts for mobile
// @Security ApiKeyAuth
// @Description  Retrieves the debt information for a specific trade point for mobile applications.
// @Tags         trade-point-debts
// @Accept       json
// @Produce      json
// @Param        id  query     int    true  "Trade Point ID"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200 {object}  object      "Debt information for the trade point"
// @Failure      400 {object}  errorResponse "Invalid Trade Point ID format"
// @Failure      500 {object}  errorResponse "Internal server error"
// @Router       /api/points/debts/mobile [get]
func (h *Handler) GetTradePointDebtsForMobile(c *gin.Context) {
	language := c.GetHeader("language")
	trade_point_id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.TradePointDebt.GetDebtForMobile(trade_point_id, language)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
