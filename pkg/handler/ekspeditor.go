package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new ekspeditor
// @Security ApiKeyAuth
// @Description Creates a new ekspeditor in the system.
// @Tags ekspeditors
// @Accept json
// @Produce json
// @Param input body models.CreateEkspeditor true "Ekspeditor information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/ekspeditors [post]
func (h *Handler) createEkspeditor(c *gin.Context) {
	var input models.CreateEkspeditor

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Ekspeditor.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Ekspeditors
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get all ekspeditors with pagination.
// @ID          get-all-ekspeditors
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        pageSize  query    int     false  "Number of ekspeditors per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "Ekspeditor name to filter by (optional)"
// @Success      200        {object} models.DriverPage "Paginated list of drivers"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/ekspeditors [get]
func (h *Handler) getAllEkspeditors(c *gin.Context) {
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

	page, err := h.services.Ekspeditor.GetAll(pageNumber, pageSize, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Ekspeditor status
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Change Ekspeditor status.
// @ID          change-ekspeditor-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Ekspeditor id"
// @Success      200    {object} int "Success"  "Ekspeditor status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "Driver not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/{id}/lock [get]
func (h *Handler) statusEkspeditor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Ekspeditor.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get ekspeditor by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific ekspeditor by its ID.
// @Tags ekspeditors
// @Produce json
// @Param id path int true "Ekspeditor ID"
// @Success 200 {object} models.Ekspeditor "Ekspeditor details"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/ekspeditors/{id} [get]
func (h *Handler) getEkspeditorById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	city, err := h.services.Ekspeditor.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, city)
}

// @Summary      Update Ekspeditor information
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Update ekspeditor.
// @ID          update-ekspeditor
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Ekspeditor id"
// @Param        input  body   models.CreateEkspeditor  true  "Ekspeditor information"
// @Success      200    {object}  int "Success" "Ekspeditor updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Ekspeditor not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/{id} [put]
func (h *Handler) updateEkspeditor(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateEkspeditor
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Ekspeditor.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete  ekspeditor
// @Security ApiKeyAuth
// @Description Deletes  ekspeditor from the system.
// @Tags ekspeditors
// @Produce json
// @Param id path int true "Ekspeditor ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/ekspeditors/{id} [delete]
func (h *Handler) deleteEkspeditor(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Ekspeditor.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Product  List
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  get prdouct list.
// @Produce      json
// @Param        day       query    int  false  "Day"
// @Param        month       query    int  false  "Month"
// @Param        year       query    int  false  "Year"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/product/all [get]
func (h *Handler) GetProductsCountForEkspeditor(c *gin.Context) {

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

	result, err := h.services.Ekspeditor.GetProductCount(day, month, year, ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Product List By Point
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  get product list by point.
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        day       query    int  false  "Day"
// @Param        month       query    int  false  "Month"
// @Param        year       query    int  false  "Year"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/product/point [get]
func (h *Handler) GetPointProductListForEkspeditor(c *gin.Context) {
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

	//userId, err := getUserId(c)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}
	//
	//ekspeditorId, err := h.services.GetEkspeditorId(userId)
	//if err != nil {
	//	newErrorResponse(c, http.StatusInternalServerError, err.Error())
	//	return
	//}

	result, err := h.services.Ekspeditor.GetPointProductList(day, month, year, tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Save Ekspeditor productlist
// @Security ApiKeyAuth
// @Description Save Ekspeditor productlist.
// @Tags ekspeditors
// @Accept json
// @Produce json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Param input body []models.CreateOrder true "ProductList"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/ekspeditors/save [post]
func (h *Handler) SaveEkspeditorProductList(c *gin.Context) {
	var input []models.CreateOrder

	if err := c.BindJSON(&input); err != nil {
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

	if err1 := h.services.Ekspeditor.SaveProductList(day, month, year, ekspeditorId, input); err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err1.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Trade Points
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get trade points.
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Success      200    {array}   []models.TradePointMobile  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/points [get]
func (h *Handler) getTradePointsForEkspeditors(c *gin.Context) {
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

	result, err := h.services.Ekspeditor.GetTradePoints(day, month, year, ekspeditorId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Order Status
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  confirm order status.
// @Produce      json
// @Param        day       query    int  false  "Day"
// @Param        month       query    int  false  "Month"
// @Param        year       query    int  false  "Year"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Success      200    {object}  int "Success"  "Order Status Changed"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/product/delivered [post]
func (h *Handler) ChangeOrderStatus(c *gin.Context) {

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

	if err1 := h.services.Ekspeditor.OrderStatus(day, month, year, tradePointId, ekspeditorId); err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err1.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Ekspeditor Galyndy For Mobile
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get Ekspeditor Gelyndy For Mobile.
// @Produce      json
// @Success      200    {object}  int "Success"  "Product List"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/galyndy [get]
func (h *Handler) GetEkspeditorGalyndy(c *gin.Context) {

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

	result, err := h.services.Ekspeditor.GetGalyndy(ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get ekspeditor List
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get ekspeditors list.
// @Produce      json
// @Param        day       query    int  false  "Day number"
// @Param        month       query    int  false  "Month number"
// @Param        year       query    int  false  "Year number"
// @Success      200    {array}   []models.EkspData  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/list [get]
func (h *Handler) getEkspeditorsProductList(c *gin.Context) {

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

	result, err := h.services.Ekspeditor.EkspeditorProducts(day, month, year, ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Product Galyndy For Admin
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get Product Galyndy For Admin.
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        ekspeditorId       query    int  false  "EkspeditorId"
// @Success      200    {array}   []models.ReadProductForTradePoint  "volumes"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/galyndy/web [get]
func (h *Handler) getEkspeditorGalyndyForAdmin(c *gin.Context) {
	language := c.GetHeader("language")

	ekspeditorId, err := strconv.Atoi(c.Query("ekspeditorId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Ekspeditor.GetEkspeditorGalyndy(ekspeditorId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Delete EkspData
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Ekspeditoryn ahli galyndylaryny pozyar.
// @Produce      json
// @Success      200    {object}  int "Success"  "Product List"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/delete [delete]
func (h *Handler) deleteEkspdata(c *gin.Context) {

	if err := h.services.Ekspeditor.DeleEkspData(); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Product  List For Mobile with minus eksp_data
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get Product  List For Mobile with minus eksp_data.
// @Produce      json
// @Param        day       query    int  false  "Day"
// @Param        month       query    int  false  "Month"
// @Param        year       query    int  false  "Year"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/product/minus [get]
func (h *Handler) GetProductsCountForEkspeditorWithMinus(c *gin.Context) {

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

	result, err := h.services.Ekspeditor.GetProductCountWithMinus(day, month, year, ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Product In The Car With Difference
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get products remaining in the car based on ekspeditor ID with calculated difference.
// @Produce      json
// @Param        language  header   string  false  "Language"
// @Success      200    {object}  interface{}    "Success"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/car/galyndy [get]
func (h *Handler) getProductInTheCarWithDifference(c *gin.Context) {
	language := c.GetHeader("language")

	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ekspeditorId, err := h.services.GetEkspeditorId(userId)
	if err != nil {
		agentId, err := h.services.GetTradeAgentId(userId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		ekspeditorId, err = h.services.GetEkspeditoryIdFromAgentId(agentId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

	}

	result, err := h.services.Ekspeditor.GetProductInTheCarWithDifference(ekspeditorId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Product Price In The Car With Difference
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get total price of products remaining in the car based on ekspeditor ID with calculated difference.
// @Produce      json
// @Success      200    {object}  map[string]int "Success"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/car/galyndy/price [get]
func (h *Handler) getProductPriceInTheCarWithDifference(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	ekspeditorId, err := h.services.GetEkspeditorId(userId)
	if err != nil {
		agentId, err := h.services.GetTradeAgentId(userId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}

		ekspeditorId, err = h.services.GetEkspeditoryIdFromAgentId(agentId)
		if err != nil {
			newErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	price, err := h.services.Ekspeditor.GetProductPriceInTheCarWithDifference(ekspeditorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"price": price})
}

// @Summary      Get Product In The Car With Difference (Array)
// @Security     ApiKeyAuth
// @Tags         ekspeditors
// @Description  Get products remaining in the car with calculated difference in array format.
// @Produce      json
// @Param        language  header   string  false  "Language"
// @Success      200    {array}   interface{}    "Success"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/ekspeditors/car/galyndy/array [get]
func (h *Handler) getProductInTheCarWithDifferenceInArray(c *gin.Context) {
	language := c.GetHeader("language")

	result, err := h.services.Ekspeditor.GetProductInTheCarWithDifferenceInArray(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
