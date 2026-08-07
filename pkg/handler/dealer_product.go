package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create DealerProduct
// @Security ApiKeyAuth
// @Tags			dealerProduct
// @Description	create dealerProduct
// @ID				create-dealerProduct
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateDealerProduct	true	"DealerProduct info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/dealer/product [post]
func (h *Handler) createDealerProduct(c *gin.Context) {
	var input models.CreateDealerProduct

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.DealerProduct.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All DealerProducts (Paginated)
// @Security     ApiKeyAuth
// @Tags         dealerProduct
// @Description  Get all dealerProducts with pagination, optional name filtering and language selection.
// @ID          get-all-dealerProducts
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of dealerProducts per page" Minimum: 1  Default: 10
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1  Default: 1
// @Param        name       query    string  false  "DealerProduct name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.DealerProductPage "Paginated list of dealerProducts"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/dealer/product [get]
func (h *Handler) getAllDealerProduct(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	// Parse 'pageSize' from query and handle errors
	pageSize, err := strconv.Atoi(c.Query("pageSize"))
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default page size if invalid or not provided
	}

	// Parse 'pageNumber' from query and handle errors
	pageNumber, err := strconv.Atoi(c.Query("pageNumber"))
	if err != nil || pageNumber <= 0 {
		pageNumber = 1 // Default page number if invalid or not provided
	}

	// Get 'name' query parameter
	name := c.Query("name")

	// Call the service layer to get the paginated list of dealerProducts
	page, err := h.services.DealerProduct.GetAll(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	c.JSON(http.StatusOK, page)
}

// @Summary      Change DealerProduct status
// @Security     ApiKeyAuth
// @Tags         dealerProduct
// @Description  Change DealerProduct status.
// @ID           change-dealerProduct-status
// @Accept       json
// @Produce      json
// @Param        DealerProductId  path  int  true  "DealerProduct ID"
// @Success      200    {object}  int "Success"  "DealerProduct status updated"
// @Failure      400    {object}  errorResponse "Invalid DealerProduct ID"
// @Failure      404    {object}  errorResponse "DealerProduct not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dealer/product/{DealerProductId}/lock [get]
func (h *Handler) statusDealerProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.DealerProduct.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update DealerProduct
// @Security     ApiKeyAuth
// @Tags         dealerProduct
// @Description  Update a dealerProduct.
// @ID          update-dealerProduct
// @Accept       json
// @Produce      json
// @Param        DealerProductId  path  int  true  "DealerProduct ID"
// @Param        input  body   models.CreateDealerProduct  true  "DealerProduct information"
// @Success      200    {object}  int "Success" "DealerProduct updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "DealerProduct not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dealer/product/{DealerProductId} [put]
func (h *Handler) updateDealerProduct(c *gin.Context) {
	PeoductId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateDealerProduct
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.DealerProduct.Update(PeoductId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete DealerProduct
// @Security     ApiKeyAuth
// @Tags         dealerProduct
// @Description  Delete a DealerProduct.
// @ID           delete-dealerProduct
// @Accept       json
// @Produce      json
// @Param        DealerProductId  path  int  true  "DealerProduct ID"
// @Success      200    {object}  int "Success"  "DealerProduct deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "DealerProduct not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/dealer/product/{DealerProductId} [delete]
func (h *Handler) deleteDealerProduct(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.DealerProduct.Delete(cityId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
