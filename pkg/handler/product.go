package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary		Create Product
// @Security ApiKeyAuth
// @Tags			product
// @Description	create product
// @ID				create-product
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateProduct	true	"Product info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/products [post]
func (h *Handler) createProduct(c *gin.Context) {
	var input models.CreateProduct

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Product.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Products (Paginated)
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Get all products with pagination, optional name filtering and language selection.
// @ID          get-all-products
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of products per page" Minimum: 1  Default: 10
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1  Default: 1
// @Param        name       query    string  false  "Product name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.ProductPage "Paginated list of products"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/products [get]
func (h *Handler) getAllProduct(c *gin.Context) {
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

	// Call the service layer to get the paginated list of products
	page, err := h.services.Product.GetAll(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	c.JSON(http.StatusOK, page)
}

// @Summary      Get ProductTypes for select (Paginated)
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Get all product types .
// @ID          get-all-products-types
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "list of products types"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/products/types/all [get]
func (h *Handler) getProductTypesForSelect(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	// Call the service layer to get the paginated list of products
	page, err := h.services.Product.GetTypes(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	c.JSON(http.StatusOK, page)
}

// @Summary      Change Product status
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Change Product status.
// @ID           change-product-status
// @Accept       json
// @Produce      json
// @Param        ProductId  path  int  true  "Product ID"
// @Success      200    {object}  int "Success"  "Product status updated"
// @Failure      400    {object}  errorResponse "Invalid Product ID"
// @Failure      404    {object}  errorResponse "Product not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/{ProductId}/lock [get]
func (h *Handler) statusProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Product.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Product by ID
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Get a product by ID.
// @ID          get-product-by-id
// @Accept       json
// @Produce      json
// @Param        ProductId  path  int  true  "Product ID"
// @Success      200    {object}  models.CreateProduct "Product"
// @Failure      400    {object}  errorResponse "Invalid product ID"
// @Failure      404    {object}  errorResponse "product not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/{ProductId} [get]
func (h *Handler) getProductById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	city, err := h.services.Product.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, city)
}

// @Summary      Get Product Brands for select (Paginated)
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Get all product brands .
// @ID          get-all-products-brands
// @Accept       json
// @Produce      json
// @Success      200        {object} []models.SelectObject "list of products brands"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/products/brands [get]
func (h *Handler) getBrandsForProduct(c *gin.Context) {

	// Call the service layer to get the paginated list of products
	page, err := h.services.Product.GetProductBrands()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Send the response
	c.JSON(http.StatusOK, page)
}

// @Summary      Update Product
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Update a product.
// @ID          update-product
// @Accept       json
// @Produce      json
// @Param        ProductId  path  int  true  "Product ID"
// @Param        input  body   models.CreateProduct  true  "Product information"
// @Success      200    {object}  int "Success" "Product updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Product not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/{ProductId} [put]
func (h *Handler) updateProduct(c *gin.Context) {
	PeoductId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateProduct
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Product.Update(PeoductId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete Product
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Delete a Product.
// @ID           delete-product
// @Accept       json
// @Produce      json
// @Param        ProductId  path  int  true  "Product ID"
// @Success      200    {object}  int "Success"  "Product deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Product not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/products/{ProductId} [delete]
func (h *Handler) deleteProduct(c *gin.Context) {
	cityId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Product.Delete(cityId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get sales Information
// @Security     ApiKeyAuth
// @Tags         product
// @Description  Get sales information.
// @ID           get-sales-information
// @Accept       json
// @Produce      json
// @Param        brandId       query    int  false  "Trade point Id"
// @Param        productTypeId       query    int  false  "Trade point Id"
// @Param        agentId       query    int  false  "Trade agent Id"
// @Param        cityId       query    int  false  "Trade point Id"
// @Param        districtId       query    int  false  "Trade point Id"
// @Param        tradePointId       query    int  false  "Trade point Id"
// @Param        channel_type_id       query    int  false  "Channel type Id"
// @Param        startDay       query    int  false "start day"
// @Param        startMonth       query    int  false  "start month"
// @Param        startYear       query    int  false  "start year"
// @Param        endDay       query    int  false  "end day"
// @Param        endMonth       query    int  false  "end month"
// @Param        endYear       query    int  false  "end year"
// @Param        taste       query    string  false  "Product name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} models.GeneralTradePoint "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/products/sales [get]
func (h *Handler) getInformationAboutSales(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	brandId, err := strconv.Atoi(c.Query("brandId"))
	if err != nil {
		brandId = 0
	}

	productTypeId, err := strconv.Atoi(c.Query("productTypeId"))
	if err != nil {
		productTypeId = 0
	}

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	channelTypeId, err := strconv.Atoi(c.Query("channelTypeId"))
	if err != nil {
		channelTypeId, err = strconv.Atoi(c.Query("channel_type_id"))
		if err != nil {
			channelTypeId = 0
		}
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

	volume, err := strconv.ParseFloat(c.Query("volume"), 64)
	if err != nil {
		volume = 0
	}

	product_taste := c.Query("taste")

	result, err := h.services.Product.GetInformantionAboutSales(brandId, productTypeId, agentId, cityId, districtId, tradePointId, channelTypeId, startDay, startMonth, startYear, endDay, endMonth, endYear, volume, product_taste, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get daily sales products in Excel
// @Description  Returns daily sales product information as a downloadable Excel file
// @Tags         product
// @Accept       json
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param        agentId    query     int     false  "Agent ID (optional, default=0)"
// @Param        day        query     int     true   "Day (1-31)"
// @Param        month      query     int     true   "Month (1-12)"
// @Param        year       query     int     true   "Year (e.g. 2025)"
// @Param        language   header    string  false  "Language preference (e.g. en, ru, tm)"
// @Success      200        {file}    file           "Excel file with daily sales products"
// @Failure      400        {object} errorResponse  "Invalid request parameters"
// @Failure      500        {object} errorResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /api/products/daily/excel [get]
func (h *Handler) getDailySalesProductInExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	Day, err := strconv.Atoi(c.Query("day"))
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

	result, err := h.services.Product.GetDailySalesProductsInExcel(agentId, Day, month, year, language)
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
