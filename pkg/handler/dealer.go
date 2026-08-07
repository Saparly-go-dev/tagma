package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create Dealer
// @Security ApiKeyAuth
// @Tags	dealer
// @Description	create dealer
// @ID				create-dealer
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateDealer	true	"dealer info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/dealer [post]
func (h *Handler) createDealer(c *gin.Context) {
	var input models.CreateDealer

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Dealer.Create(input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Dealers Table
// @Security     ApiKeyAuth
// @Tags         dealer
// @Description  Get all  dealers with pagination.
// @ID          get-dealers-all-data
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of Dealers per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param		 cityId 	query	int 	false	"City Id"
// @Param        name       query    string  false  "Dealer name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.DealerPage "Paginated list of cities"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/dealer [get]
func (h *Handler) getDealersByPage(c *gin.Context) {
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

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	name := c.Query("name")

	page, err := h.services.Dealer.GetPage(pageSize, pageNumber, cityId, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Update Dealer
// @Security     ApiKeyAuth
// @Tags         dealer
// @Description  Update a dealer.
// @ID          update-dealer
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Dealer ID"
// @Param        input  body   models.CreateDealer  true  "Dealer information"
// @Success      200    {object}  int "Success" "Dealer updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Dealer not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dealer/{id} [put]
func (h *Handler) updateDealer(c *gin.Context) {
	dealerId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateDealer
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Dealer.Update(dealerId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Change Dealer  status
// @Security     ApiKeyAuth
// @Tags         dealer
// @Description  Change a dealer  status.
// @ID          change_dealer_status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Dealer ID"
// @Success      200    {object} int "Success"  "Dealer Order status updated"
// @Failure      400    {object}  errorResponse "Invalid Dealer Order ID"
// @Failure      404    {object}  errorResponse "Dealer Order not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dealer/{id}/lock [get]
func (h *Handler) changeDealerStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Dealer.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get  Dealers by name
// @Security     ApiKeyAuth
// @Tags         dealer
// @Description  Get all  dealers by name.
// @ID          get-dealers-by-name
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Dealer name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.SelectObject "list of dealers"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/dealer/all [get]
func (h *Handler) getDealersByName(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.Dealer.GetDealerForSelect(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Top up Dealers account
// @Security     ApiKeyAuth
// @Tags         dealer
// @Description  top up dealers account.
// @ID          top-up-dealer
// @Accept      multipart/form-data
// @Produce      json
// @Param        dealerId       formData    string  false  "Dealer Id"
// @Param        price1       formData    string  false  "Currency in Manat"
// @Param        price2       formData    string  false  "Currency in Tenne"
// @Success      200    {object} int "Success"  "Dealer account top up"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/dealer/topup [post]
func (h *Handler) TopUpDealer(c *gin.Context) {
	// Get the 'language' header
	dealerId, err := strconv.Atoi(c.PostForm("dealerId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	price1, err := strconv.Atoi(c.PostForm("price1"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	price2, err := strconv.Atoi(c.PostForm("price2"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if result_error := h.services.Dealer.TopUpDealerAccount(dealerId, price1, price2); result_error != nil {
		newErrorResponse(c, http.StatusInternalServerError, result_error.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
