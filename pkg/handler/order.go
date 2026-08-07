package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary      Get Trade Points
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get trade points.
// @ID          get-trade-points
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.TradePointMobile  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/points [get]
func (h *Handler) getTradePointsForOrder(c *gin.Context) {
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

	result, err := h.services.Order.GetTradePoints(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Points For Ekspeditor
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get trade points For Ekspeditor.
// @ID          get-trade-points-ekspeditor
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.TradePointMobile  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/ekspeditor/points [get]
func (h *Handler) getTradePointsForOrderFromEkspeditorId(c *gin.Context) {
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
	if err != nil || agentId == 0 {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.Order.GetTradePoints(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Trade Points For marketing
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get trade points.
// @ID          get-trade-points-marketing
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []models.TradePointMobile  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/points/marketing [get]
func (h *Handler) getTradePointsForMarketing(c *gin.Context) {
	language := c.GetHeader("language")

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

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

	agentId, err := h.services.GetTradeAgentId(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	result, err := h.services.Order.GetTradePointsForMarketing(day, month, year, agentId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Products
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get all products.
// @ID          get-products-to-order
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {object} []models.ProductDistinct "products"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/products [get]
func (h *Handler) getProductsForOrders(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	result, err := h.services.Order.GetProducts(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Brands
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get all brands.
// @ID          get-all-brands
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []string "brands"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/brands [get]
func (h *Handler) getBrandsForOrders(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	result, err := h.services.Order.GetBrand(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Models
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get models.
// @ID          get-models
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        brandId       query    int  false  "product brand Id"
// @Param        productTypeId       query    int  false  "product type id"
// @Success      200    {array}   []string "models"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/models [get]
func (h *Handler) getModelsForOrder(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	brandId, err := strconv.Atoi(c.Query("brandId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	productTypeId, err := strconv.Atoi(c.Query("productTypeId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Order.GetModel(brandId, productTypeId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Product Types
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get all types.
// @ID          get-all-types
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   []string "brands"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/types [get]
func (h *Handler) getTypesForOrders(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	result, err := h.services.Order.GetTypes(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Volumes
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get volumes.
// @ID          get-volumes
// @Produce      json
// @Param        brand       query    string  false  "name of product brand"
// @Param        model       query    string  false  "name of product model"
// @Param        type       query    string  false  "name of product model"
// @Success      200    {array}   []models.SelectObject  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/volumes [get]
func (h *Handler) getVolumesForOrder(c *gin.Context) {

	brand := c.Query("brand")
	model := c.Query("model")
	tip := c.Query("type")

	result, err := h.services.Order.GetVolume(brand, model, tip)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Models
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get sum.
// @ID          get-sum-of-order
// @Produce      json
// @Param        id       query    int  false  "product Id"
// @Param        count       query    int  false  "product Count"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/sum [get]
func (h *Handler) getSumForOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Order.GetSum(id, count)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Save the Orders
// @Security     ApiKeyAuth
// @Tags         order
// @Description  save the order.
// @ID           save-order
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        first       query    bool  false  "Trade Point Id"
// @Param		 is_credit query bool false "Is the order created for Credit"
// @Param		 paymentTypeId	query 	int	false "Payment type Id for Order"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        input body []models.CreateOrder true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/save [post]
func (h *Handler) saveTheOrder(c *gin.Context) {
	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

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

	first, err := strconv.ParseBool(c.Query("first"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	is_credit, err := strconv.ParseBool(c.Query("is_credit"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	paymentTypeId, err := strconv.Atoi(c.Query("paymentTypeId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input []models.CreateOrder

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//userId, error := getUserId(c)
	//if error != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, error.Error())
	//	return
	//}

	err1 := h.services.Order.SaveOrder(day, month, year, tradePointId, paymentTypeId, first, is_credit, input)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Marketing
// @Security     ApiKeyAuth
// @Tags         order
// @Description  marketing post.
// @ID           marketing
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        input body []models.MarketingPost true "list of Marketing Post"
// @Success      200    {object}  int "Success"  "Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/marketing [post]
func (h *Handler) marketingForAgent(c *gin.Context) {
	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

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

	var input []models.MarketingPost

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.Order.Marketing(day, month, year, tradePointId, input)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Orders page
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get Orders Page.
// @ID           get-orders-page
// @Produce      json
// @Param        language      header    string  false  "Language (optional)"
// @Param        day           query     int     false  "Day of the order date (optional, defaults to 0)"
// @Param        month         query     int     false  "Month of the order date (optional, defaults to 0)"
// @Param        year          query     int     false  "Year of the order date (optional, defaults to 0)"
// @Param        pageSize      query     int     true   "Page size (REQUIRED)"
// @Param        pageNumber    query     int     true   "Page number (REQUIRED)"
// @Param        agentId       query     int     false  "Trade Agent Id (optional, defaults to 0)"
// @Param        cityId        query     int     false  "City Id (optional, defaults to 0)"
// @Param        districtId    query     int     false  "District Id (optional, defaults to 0)"
// @Param        location      query     string  false  "Trade Point Location (optional)"
// @Param        code          query     string  false  "Trade Point code (optional)"
// @Param        name          query     string  false  "Trade Point Name (optional)"
// @Param        order         query     string  false  "Field name for sorting (e.g., 'date')"
// @Param        status        query     bool    false  "Order status (e.g., 'true' or 'false'. Defaults to false on error/missing)"
// @Success      200           {object}  models.OrderPage "Success"  "Orders List"
// @Failure      400           {object}  errorResponse "Bad request (e.g., pageSize or pageNumber missing/invalid)"
// @Failure      500           {object}  errorResponse "Internal server error"
// @Router       /api/orders/page [get]
func (h *Handler) getOrdersPage(c *gin.Context) {
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

	day, err := strconv.Atoi(c.Query("day"))
	if err != nil {
		day = 0
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		month = 0
	}

	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		year = 0
	}

	name := c.Query("name")
	location := c.Query("location")
	code := c.Query("code")
	order := c.Query("order")

	tradeAgentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		tradeAgentId = 0
	}

	is_status_active := true

	status, err := strconv.ParseBool(c.Query("status"))
	if err != nil {
		is_status_active = false
	}

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	result, err := h.services.Order.GetOrderPage(pageSize, pageNumber, tradeAgentId, cityId, districtId, day, month, year, name, language, location, code, order, status, is_status_active)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Orders List by Order Id
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get OrderList By Id.
// @ID          get-orders-list
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        orderId       query    int  false  "Order Id"
// @Success      200    {object}  []models.ReadOrderList "Success"  "ReadOrderList list"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/list [get]
func (h *Handler) getOrderListById(c *gin.Context) {
	language := c.GetHeader("language")

	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Order.GetOrderList(orderId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Orders excel
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get Orders Page.
// @ID          get-orders-excel
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        agentId       query    int  false  "Trade Agent Id"
// @Param        cityId       query    int  false  "City Id"
// @Param        districtId       query    int  false  "District Id"
// @Param        name       query    string  false  "Trade Point Name"
// @Success      200    {object}  models.OrderPage "Success"  "Orders List"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/excel [get]
func (h *Handler) orderExcel(c *gin.Context) {
	language := c.GetHeader("language")

	name := c.Query("name")

	tradeAgentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		tradeAgentId = 0
	}

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	result, err := h.services.Order.OrderExcel(tradeAgentId, cityId, districtId, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=orders.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// @Summary      Get Orders
// @Security     ApiKeyAuth
// @Tags         order
// @Description  Get Orders.
// @ID          get-orders
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        day       query    int  false  "day"
// @Param        month       query    int  false  "month"
// @Param        year       query    int  false  "year"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Success      200    {object}  models.GetOrderObject "Success"  "Orders List"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders [get]
func (h *Handler) getOrders(c *gin.Context) {
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

	result, err := h.services.Order.GetOrders(day, month, year, tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Close The Order
// @Security     ApiKeyAuth
// @Tags         order
// @Description  close the order.
// @ID           close-order
// @Produce      json
// @Param        orderId       query    int  false  "Order Id number"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/orders/close [post]
func (h *Handler) closeTheOrder(c *gin.Context) {
	orderId, err := strconv.Atoi(c.Query("orderId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.Order.CloseTheOrder(orderId)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
