package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary      Save the Dealer Orders
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description  save the dealer_order.
// @ID           save-dealer-order
// @Produce      json
// @Param        dealerId       query    int  false  "Dealer Id"
// @Param        input body []models.CreateOrder true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dorders/save [post]
func (h *Handler) saveDealerOrder(c *gin.Context) {
	dealerId, err := strconv.Atoi(c.Query("dealerId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input []models.CreateOrder

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.DealerOrder.Save(dealerId, input)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update the Dealer Orders
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description  update the dealer_order.
// @ID           update-dealer-order
// @Produce      json
// @Param        dealerOrderId       query    int  false  "Dealer Order Id"
// @Param        input body []models.CreateOrder true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dorders/update [post]
func (h *Handler) updateDealerOrder(c *gin.Context) {
	dealerOrderId, err := strconv.Atoi(c.Query("dealerOrderId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input []models.CreateOrder

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.DealerOrder.Update(dealerOrderId, input)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Complete Dealer Order
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description  Complete Delaer Order By id.
// @ID          complete-dorder
// @Produce      json
// @Param        id       query    int  false  "Dealer Order id"
// @Success      200    {object}  object "Success"  "ReadOrderList list"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dorders/complete [post]
func (h *Handler) completeDealerOrder(c *gin.Context) {

	Id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.DealerOrder.CompleTheDealerOrder(Id)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Change Dealer Order status
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description  Change a dealer order status.
// @ID          change_dorder_status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Dealer Order ID"
// @Success      200    {object} int "Success"  "Dealer Order status updated"
// @Failure      400    {object}  errorResponse "Invalid Dealer Order ID"
// @Failure      404    {object}  errorResponse "Dealer Order not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dorders/{id}/lock [get]
func (h *Handler) changeDealerOrderStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.DealerOrder.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Dealer Order information
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description   Get Dealer Order page.
// @ID           get-dealer-order-page
// @Accept       json
// @Produce      json
// @Param        pageSize       query    int  false  "Page size"
// @Param        pageNumber       query    int  false  "Page number"
// @Param        dealerId       query    int  false  "Dealer Order Id"
// @Param        startDay       query    int  false "start day"
// @Param        startMonth       query    int  false  "start month"
// @Param        startYear       query    int  false  "start year"
// @Param        endDay       query    int  false  "end day"
// @Param        endMonth       query    int  false  "end month"
// @Param        endYear       query    int  false  "end year"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.DealerOrderPage "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/dorders/ [get]
func (h *Handler) getDealerOrderPage(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	dealerId, err := strconv.Atoi(c.Query("dealerId"))
	if err != nil {
		dealerId = 0
	}

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

	result, err := h.services.DealerOrder.GetOrders(pageSize, pageNumber, dealerId, startDay, startMonth, startYear, endDay, endMonth, endYear, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Dealer Orders List by Order Id
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description  Get OrderList By Id.
// @ID          get-dorders-list
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        id       query    int  false  "Dealer Order Id"
// @Success      200    {object}  []models.ReadOrderList "Success"  "ReadOrderList list"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dorders/list [get]
func (h *Handler) getDealerOrderListById(c *gin.Context) {
	language := c.GetHeader("language")

	Id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.DealerOrder.GetOrderList(Id, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Dealer Order information By Product
// @Security     ApiKeyAuth
// @Tags         dealer_order
// @Description   Get Dealer Order by products.
// @ID           get-dealer-order-products
// @Accept       json
// @Produce      json
// @Param        dealerId       query    int  false  "Dealer Id"
// @Param        startDay       query    int  false "start day"
// @Param        startMonth       query    int  false  "start month"
// @Param        startYear       query    int  false  "start year"
// @Param        endDay       query    int  false  "end day"
// @Param        endMonth       query    int  false  "end month"
// @Param        endYear       query    int  false  "end year"
// @Param		 status		   query	bool false "status of Product"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.ReadProductForTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/dorders/product [get]
func (h *Handler) getDealerOrderByProducts(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	dealerId, err := strconv.Atoi(c.Query("dealerId"))
	if err != nil {
		dealerId = 0
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

	status, err := strconv.ParseBool(c.Query("status"))
	if err != nil {
		status = false
	}

	result, err := h.services.DealerOrder.GetProduct(dealerId, startDay, startMonth, startYear, endDay, endMonth, endYear, status, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Delete dealer order
// @Description  Deletes a dealer order by its ID
// @Tags         dealer_order
// @Accept       json
// @Produce      json
// @Param        id   query     int     true   "Dealer Order ID"
// @Success      200  {object}  statusResponse   "Successfully deleted"
// @Failure      400  {object}  errorResponse    "Invalid dealer order ID"
// @Failure      500  {object}  errorResponse    "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/dorders [delete]
func (h *Handler) deleteDealerOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result_error := h.services.DealerOrder.DeleteDealerOrder(id)
	if result_error != nil {
		newErrorResponse(c, http.StatusInternalServerError, result_error.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "Ok"})
}
