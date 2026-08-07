package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create ProductType
// @Security ApiKeyAuth
// @Tags			productType
// @Description	create productType
// @ID				create-productType
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateProductType	true	"ProductType info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/products/types [post]
func (h *Handler) createProductType(c *gin.Context) {
	var input models.CreateProductType

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.ProductType.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All ProductTypes
// @Security     ApiKeyAuth
// @Tags         productType
// @Description  Get all ProductTypes with pagination.
// @ID          get-all-productTypes
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of ProductTypes per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "ProductType name to filter by (optional)"
// @Success      200        {object} models.ProductTypePage "Paginated list of productTypes"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/products/types  [get]
func (h *Handler) getAllProductType(c *gin.Context) {
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

	page, err := h.services.ProductType.GetAll(pageSize, pageNumber, name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change ProductType status
// @Security     ApiKeyAuth
// @Tags         productType
// @Description  Change ProductType status.
// @ID          change-productType-status
// @Accept       json
// @Produce      json
// @Param        ProductTypeId  path  int  true  "ProductType ID"
// @Success      200    {object}  int "Success"  "ProductType status updated"
// @Failure      400    {object}  errorResponse "Invalid ProductType ID"
// @Failure      404    {object}  errorResponse "ProductType not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/types/{ProductTypeId}/lock [get]
func (h *Handler) statusProducttype(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.ProductType.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get ProductType by ID
// @Security     ApiKeyAuth
// @Tags         productType
// @Description  Get a productType by ID.
// @ID          get-productType-by-id
// @Accept       json
// @Produce      json
// @Param        ProductTypeId  path  int  true  "Product Type ID"
// @Success      200    {object}  models.ProductType "ProductType"
// @Failure      400    {object}  errorResponse "Invalid productType ID"
// @Failure      404    {object}  errorResponse "productType not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/types/{ProductTypeId} [get]
func (h *Handler) GetProductTypeById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	productType, err := h.services.ProductType.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, productType)
}

// @Summary      Update ProductType
// @Security     ApiKeyAuth
// @Tags         productType
// @Description  Update a productType.
// @ID          update-productType
// @Accept       json
// @Produce      json
// @Param        ProductTypeId  path  int  true  "Product Type ID"
// @Param        input  body   models.CreateProductType  true  "ProductType information"
// @Success      200    {object}  int "Success" "ProductType updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "ProductType not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/types/{ProductTypeId} [put]
func (h *Handler) updateProductType(c *gin.Context) {
	PeoductTypeId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateProductType
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.ProductType.Update(PeoductTypeId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete ProductType
// @Security     ApiKeyAuth
// @Tags         productType
// @Description  Delete a ProductType.
// @ID          delete-productType
// @Accept       json
// @Produce      json
// @Param        ProductTypeId  path  int  true  "Product Type ID"
// @Success      200    {object}  int "Success"  "ProductType deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "ProductType not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/types/{ProductTypeId} [delete]
func (h *Handler) deleteProductType(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.ProductType.Delete(cityId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary     Get Product type volumes
// @Security     ApiKeyAuth
// @Tags         productType
// @Description  Get volumes for  ProductType.
// @ID          get-volume
// @Accept       json
// @Produce      json
// @Success      200    {object}  []models.Volume  "volumes"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Volumes not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/types/volume [get]
func (h *Handler) getProductTypeVolumes(c *gin.Context) {

	result, err := h.services.ProductType.GetAllVolumes()

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)

	//c.JSON(http.StatusOK, statusResponse{
	//	Status: "Ok",
	//})
}
