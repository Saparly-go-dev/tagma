package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary      Get Trade Points
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  Get trade points.
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.TradePointForPayment  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/points [get]
func (h *Handler) getTradePointsForPayment(c *gin.Context) {
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

	result, err := h.services.Payment.GetTradePointsForPayment(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Points for ekspeditors
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  Get trade points for ekspeditors.
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.TradePointForPayment  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/ekspeditor/points [get]
func (h *Handler) getTradePointsForPaymentFromEkspeditor(c *gin.Context) {
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

	ekspeditorId, err := h.services.GetEkspeditorId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentIdFromEkspeditorId(ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.Payment.GetTradePointsForPayment(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

type PostPayment struct {
	OrderId       int     `json:"order_id"`
	TradePointId  int     `json:"trade_point_id"`
	Currency      float64 `json:"currency"`
	PaymentTypeId int     `json:"payment_type_id"`
}

// @Summary      Save the Payment
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  save the payment.
// @ID           save-payment
// @Produce      json
// @Param        input body PostPayment true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/save [post]
func (h *Handler) saveThePayment(c *gin.Context) {

	var input PostPayment

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	currency := input.Currency

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

	var paymentErr error
	if input.OrderId > 0 {
		paymentErr = h.services.Payment.SavePayment(agentId, models.CreatePayment{
			OrderId:       input.OrderId,
			Currency:      currency,
			PaymentTypeId: input.PaymentTypeId,
			Is_Agent:      true,
		})
	} else if input.TradePointId > 0 {
		paymentErr = h.services.Payment.PostPayment(currency, agentId, input.TradePointId, input.PaymentTypeId, true)
	} else {
		newErrorResponse(c, http.StatusBadRequest, "order_id or trade_point_id is required")
		return
	}

	if paymentErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, paymentErr.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Save the Payment from ekspeditor
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  save the payment form ekspeditor.
// @ID           save-payment-ekspeditor
// @Produce      json
// @Param        input body PostPayment true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/ekspeditor/save [post]
func (h *Handler) saveThePaymentFromEkspeditor(c *gin.Context) {

	var input PostPayment

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	currency := input.Currency

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ekspeditorId, err := h.services.GetEkspeditorId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentIdFromEkspeditorId(ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var paymentErr error
	if input.OrderId > 0 {
		paymentErr = h.services.Payment.SavePayment(agentId, models.CreatePayment{
			OrderId:       input.OrderId,
			Currency:      currency,
			PaymentTypeId: input.PaymentTypeId,
			Is_Agent:      false,
		})
	} else if input.TradePointId > 0 {
		paymentErr = h.services.Payment.PostPayment(currency, agentId, input.TradePointId, input.PaymentTypeId, false)
	} else {
		newErrorResponse(c, http.StatusBadRequest, "order_id or trade_point_id is required")
		return
	}

	if paymentErr != nil {
		newErrorResponse(c, http.StatusInternalServerError, paymentErr.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Payment Types
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  Get payment types.
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.SelectObject  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/types [get]
func (h *Handler) getPaymentTypes(c *gin.Context) {
	language := c.GetHeader("language")

	result, err := h.services.Payment.GetPaymentTypes(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Points
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  Get trade points.
// @Produce      json
// @Param        orderId       query    int  false  "Day number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.ReadPayment  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/order [get]
func (h *Handler) getPaymentsForOrder(c *gin.Context) {
	language := c.GetHeader("language")

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Payment.GetPaymentForOrder(orderId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Agent Payments
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  Get trade agent payments.
// @Produce      json
// @Param        agentId       query    int  false  "Trade agent id"
// @Param        paymentTypeId       query    int  false  "Payment type id"
// @Param        status       query    bool  false  "Payment status"
// @Param		 code		query		string	false "Trade Point Code"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.ReadPayment  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/list [get]
func (h *Handler) GetAgentPayments(c *gin.Context) {
	language := c.GetHeader("language")

	code := c.Query("code")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	paymentTypeId, err := strconv.Atoi(c.Query("paymentTypeId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	status, err := strconv.ParseBool(c.Query("status"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Payment.GetAgentPayments(agentId, paymentTypeId, status, code, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Confirm the Payment
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  confirm the payment.
// @ID           confirm-payment
// @Produce      json
// @Param        agentId       query    int  false  "Trade agent id"
// @Param        price1       query    int  false  "Price 1"
// @Param        price2       query    int  false  "Price 2"
// @Param		 ids   body []int	false "Payment id list"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/confirm [post]
func (h *Handler) confirmPayment(c *gin.Context) {

	agentId, err := strconv.Atoi(c.Query("agentId"))
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

	var ids []int
	if err := c.BindJSON(&ids); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.Payment.ConfirmPayments(agentId, price1, price2, ids)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err1.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Payment currency
// @Security     ApiKeyAuth
// @Tags         update payment
// @Description  debt.
// @ID           update-payment
// @Produce      json
// @Param        paymentId       query    int  false  "payment id"
// @Param        price1       query    int  false  "Manat coefficient"
// @Param        price2       query    int  false  "tenne coefficient"
// @Success      200    {object}  int "Success"  "Success"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/update [get]
func (h *Handler) updatePayment(c *gin.Context) {

	payment_id, err := strconv.Atoi(c.Query("paymentId"))
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

	err1 := h.services.Payment.UpdatePayment(payment_id, price1, price2)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err1.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Debt
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  debt.
// @ID           debt
// @Produce      json
// @Param        agentId       query    int  false  "Day number"
// @Param        price1       query    int  false  "Day number"
// @Param        price2       query    int  false  "Month number"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/debt [get]
func (h *Handler) updateDebt(c *gin.Context) {

	agentId, err := strconv.Atoi(c.Query("agentId"))
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

	err1 := h.services.Payment.UpdateDebt(agentId, price1, price2)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err1.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Debt
// @Security     ApiKeyAuth
// @Tags         payment
// @Description  debt.
// @ID           addSelfDebt
// @Produce      json
// @Param        agentId       query    int  false  "Day number"
// @Param        price1       query    int  false  "Day number"
// @Param        price2       query    int  false  "Month number"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/payments/add_debt [get]
func (h *Handler) addSelfDebtForAgent(c *gin.Context) {
	agentId, err := strconv.Atoi(c.Query("agentId"))
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

	err = h.services.Payment.AddSelfDebt(agentId, currency)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

// @Summary      Get total debt for trade point
// @Security     ApiKeyAuth
// @Description  Returns the total debt amount for a specific trade point
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        tradePointId  query     int  true  "Trade Point ID"
// @Success      200           {object}  object      "Total debt amount"
// @Failure      400           {object}  errorResponse  "Invalid tradePointId parameter"
// @Router       /api/payments/trade-point-debt [get]
func (h *Handler) getTradePointAllDebtSum(c *gin.Context) {
	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result := h.services.Payment.GetTradePointAllDebtSum(tradePointId)

	c.JSON(http.StatusOK, result)
}

type Post_Payment struct {
	Currency      float64 `json:"currency"`
	TradePointId  int     `json:"trade_point_id"`
	PaymentTypeId int     `json:"payment_type_id"`
}

// @Summary      Create payment
// @Description  Creates a new payment for a trade point
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        input  body      Post_Payment  true  "Payment data"
// @Success      200    {object}  statusResponse  "Payment created successfully"
// @Failure      400    {object}  errorResponse   "Invalid input data"
// @Failure      500    {object}  errorResponse   "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/payments/create [post]
func (h *Handler) PostPayment(c *gin.Context) {

	var input Post_Payment
	is_agent := true

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	currency := input.Currency

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		ekspeditorId, err2 := h.services.GetEkspeditorId(userId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}

		agentId, err2 = h.services.GetTradeAgentIdFromEkspeditorId(ekspeditorId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}
		is_agent = false
	}

	result := h.services.PostPayment(currency, agentId, input.TradePointId, input.PaymentTypeId, is_agent)

	if result != nil {
		newErrorResponse(c, http.StatusInternalServerError, result.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}

// @Summary      Get trade point payment history
// @Description  Returns payment history for a specific trade point
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        tradePointId  query     int     true   "Trade Point ID"
// @Param        language      header    string  false  "Language preference"
// @Success      200           {object}  object         "Payment history"
// @Failure      400           {object}  errorResponse  "Invalid tradePointId parameter"
// @Failure      500           {object}  errorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/payments/history [get]
func (h *Handler) GetTradePointPaymentsHistory(c *gin.Context) {
	language := c.GetHeader("language")
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		ekspeditorId, err2 := h.services.GetEkspeditorId(userId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}

		agentId, err2 = h.services.GetTradeAgentIdFromEkspeditorId(ekspeditorId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}
	}

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Payment.GetTradePointPaymentsHistory(agentId, tradePointId, language)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get daily payment information
// @Description  Returns payment information for a specific day
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        language   header    string  false  "Language preference (e.g. en, ru, tm)"
// @Success      200        {object}  object         "Daily payment information"
// @Failure      400        {object}  errorResponse  "Invalid request parameters"
// @Failure      500        {object}  errorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/payments/daily [get]
func (h *Handler) getDailyInformationAboutPayment(c *gin.Context) {
	language := c.GetHeader("language")

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		ekspeditorId, err2 := h.services.GetEkspeditorId(userId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}

		agentId, err2 = h.services.GetTradeAgentIdFromEkspeditorId(ekspeditorId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}
	}

	now := time.Now()
	day := now.Day()
	month := int(now.Month())
	year := now.Year()

	result, err := h.services.Payment.GetDailyPaymentsInformations(agentId, day, month, year, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get All TradePoints list with debts
// @Description  Returns payment information for right now
// @Tags         payment
// @Accept       json
// @Produce      json
// @Param        language   header    string  false  "Language preference (e.g. en, ru, tm)"
// @Success      200        {object}  object         "Daily payment information"
// @Failure      400        {object}  errorResponse  "Invalid request parameters"
// @Failure      500        {object}  errorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/payments/trade-points-debt-list [get]
func (h *Handler) getTradePointsDebtList(c *gin.Context) {
	language := c.GetHeader("language")

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		ekspeditorId, err2 := h.services.GetEkspeditorId(userId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}

		agentId, err2 = h.services.GetTradeAgentIdFromEkspeditorId(ekspeditorId)
		if err2 != nil {
			newErrorResponse(c, http.StatusInternalServerError, err2.Error())
			return
		}
	}

	result, err := h.services.Payment.GetTradePointsDebtList(agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
