package handler

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create Trade Category
// @Security ApiKeyAuth
// @Tags			category
// @Description	create category
// @ID				create-category
// @Accept			json
// @Produce		json
// @Param			input	body	tagma.CreateTrade_category	true	"trade category info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/categories [post]
func (h *Handler) createCategory(c *gin.Context) {
	var input tagma.CreateTrade_category

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Trade_Category.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary      Get All Trade Categories
// @Security     ApiKeyAuth
// @Tags         category
// @Description  Get all trade categories.
// @ID          get-all-categories
// @Produce      json
// @Param        pageSize  query    int     false  "Number of trade categories per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "trade categories name to filter by (optional)"
// @Success      200    {array}   tagma.TradeCategoryPage "Trade Categories"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/categories [get]
func (h *Handler) getAllCategories(c *gin.Context) {
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

	page, err := h.services.Trade_Category.GetAll(pageNumber, pageSize, name)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)

}

// @Summary      Change TradeCategory status
// @Security     ApiKeyAuth
// @Tags         category
// @Description  Change Trade Category status.
// @ID          change-trade-category-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradeCategory ID"
// @Success      200    {object}  int "Success"  "Trade Category updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "City not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/categories/{id}/lock [get]
func (h *Handler) statusTradeCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Trade_Category.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Trade Categories by ID
// @Security     ApiKeyAuth
// @Tags         category
// @Description  Get a category by ID.
// @ID          get-trade-category-by-id
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Category ID"
// @Success      200    {object}  tagma.Trade_category "Trade Category"
// @Failure      400    {object}  errorResponse "Invalid Trade category ID"
// @Failure      404    {object}  errorResponse "Trade category not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/categories/{id} [get]
func (h *Handler) getCategoryById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	category, err := h.services.Trade_Category.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, category)
}

// @Summary      Update Trade Category
// @Security     ApiKeyAuth
// @Tags         category
// @Description  Update a trade category.
// @ID          update-category
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade category ID"
// @Param        input  body   tagma.CreateTrade_category  true  "Trade Category information"
// @Success      200    {object}  int "Success" "Trade Category updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Category not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/categories/{id} [put]
func (h *Handler) updateCategory(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input tagma.CreateTrade_category
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Category.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete Trade Category
// @Security     ApiKeyAuth
// @Tags         category
// @Description  Delete a trade category.
// @ID          delete-category
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Category ID"
// @Success      200    {object}  int "Success"  "Trade Category deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Category not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/categories/{id} [delete]
func (h *Handler) deleteCategory(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Category.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
