package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary      Get furniture information for a trade point (mobile)
// @Security     ApiKeyAuth
// @Tags         furniture
// @Description  Get all furniture items for a specific trade point, with language support for mobile applications
// @Accept       json
// @Produce      json
// @Param        language header string false "Language preference (e.g., 'ru', 'tm')"
// @Param        trade_point_id query int true "Trade Point ID to filter furniture"
// @Success      200 {object} object "List of furniture items with localized names"
// @Failure      400 {object} errorResponse "Invalid trade point ID"
// @Failure      500 {object} errorResponse "Internal server error"
// @Router       /api/furniture/mobile [get]
func (h *Handler) GetTradePointFurnituresForMobile(c *gin.Context) {
	language := c.GetHeader("language")

	trade_point_id, err := strconv.Atoi(c.Query("trade_point_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Furniture.GetMobile(trade_point_id, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get paginated furniture
// @Description  Returns a paginated list of furniture, optionally filtered by name.
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        pageSize    query     int    true    "Number of items per page"
// @Param        pageNumber  query     int    true    "Page number to retrieve"
// @Param        name        query     string  false   "Filter by furniture name"
// @Success      200       {object}  models.FurniturePage  "Paginated list of furniture"
// @Failure      400       {object}  errorResponse                     "Invalid page size or number"
// @Failure      500       {object}  errorResponse                     "Internal server error"
// @Router       /api/furniture [get]
func (h *Handler) GetFurnituresPage(c *gin.Context) {

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

	result, err := h.services.Furniture.Get(name, pageSize, pageNumber)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Create a new furniture item
// @Description  Creates a new furniture item with the provided details.
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        input  body      models.CreateFurniture  true  "Furniture object to be created"
// @Success      200    {object}  statusResponse          "Successfully created"
// @Failure      400    {object}  errorResponse         "Invalid request body"
// @Failure      500    {object}  errorResponse         "Internal server error"
// @Router       /api/furniture [post]
func (h *Handler) CreateFurniture(c *gin.Context) {
	var input models.CreateFurniture

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Furniture.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update a furniture item (ID in query)
// @Description  Updates an existing furniture item with the provided ID (passed in the query).
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        id      query     int                     true  "Furniture ID to be updated"
// @Param        input   body      models.CreateFurniture  true  "Furniture object with updated values"
// @Success      200     {object}  statusResponse          "Successfully updated"
// @Failure      400     {object}  errorResponse         "Invalid query parameter or request body"
// @Failure      500     {object}  errorResponse         "Internal server error"
// @Router       /api/furniture [put]
func (h *Handler) UpdateFurniture(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input models.CreateFurniture

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Furniture.Update(id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete a furniture item (ID in query)
// @Description  Deletes the furniture item with the specified ID (passed in the query).
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        id  query     int    true    "Furniture ID to be deleted"
// @Success      200 {object}  statusResponse  "Successfully deleted"
// @Failure      400 {object}  errorResponse   "Invalid query parameter"
// @Failure      500 {object}  errorResponse   "Internal server error"
// @Router       /api/furniture [delete]
func (h *Handler) DeleteFurniture(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Furniture.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get furniture list for select
// @Description  Returns a list of furniture items formatted for select inputs, with language support.
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        language header string false "Language preference (e.g., 'ru', 'tm')"
// @Param        name     query  string false "Filter by furniture name"
// @Success      200   {array}   []models.SelectObject  "List of furniture for select"
// @Failure      500   {object}  errorResponse          "Internal server error"
// @Router       /api/furniture/select [get]
func (h *Handler) GetFurnitureListForSelect(c *gin.Context) {
	language := c.GetHeader("language")
	name := c.Query("name")

	result, err := h.services.Furniture.GetFurnitureList(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get furniture for a trade point
// @Description  Get all furniture items associated with a specific trade point, with language support.
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        language        header string false "Language preference (e.g., 'ru', 'tm')"
// @Param        trade_point_id  query  int    true  "Trade Point ID to filter furniture"
// @Success      200   {array}   []object  "List of furniture items for the trade point"
// @Failure      400   {object}  errorResponse "Invalid trade point ID"
// @Failure      500   {object}  errorResponse "Internal server error"
// @Router       /api/furniture/point [get]
func (h *Handler) GetFurnituresForTradePoint(c *gin.Context) {
	language := c.GetHeader("language")
	trade_point_id, err := strconv.Atoi(c.Query("trade_point_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Furniture.GetTradePointFurnitures(trade_point_id, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Create a furniture association for a trade point
// @Description  Associates a furniture item with a specific trade point.
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        trade_point_id  query  int    true  "Trade Point ID to associate with"
// @Param        furniture_id    query  int    true  "Furniture ID to be associated"
// @Param        count           query  int    true  "Number of furniture items"
// @Success      200   {object}  statusResponse "Successfully created association"
// @Failure      400   {object}  errorResponse "Invalid parameters"
// @Failure      500   {object}  errorResponse "Internal server error"
// @Router       /api/furniture/point_create [post]
func (h *Handler) CreateFurnitureForTradePoint(c *gin.Context) {
	trade_point_id, err := strconv.Atoi(c.Query("trade_point_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	furniture_id, err := strconv.Atoi(c.Query("furniture_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Furniture.CreateTradePointFurniture(trade_point_id, furniture_id, count)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Create a furniture association for a trade point
// @Description  Associates a furniture item with a specific trade point.
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        id  query  int    true  "ID to associate with"
// @Param        furniture_id    query  int    true  "Furniture ID to be associated"
// @Param        count           query  int    true  "Number of furniture items"
// @Success      200   {object}  statusResponse "Successfully created association"
// @Failure      400   {object}  errorResponse "Invalid parameters"
// @Failure      500   {object}  errorResponse "Internal server error"
// @Router       /api/furniture/point_update [put]
func (h *Handler) UpdateFurnitureForTradePoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	furniture_id, err := strconv.Atoi(c.Query("furniture_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	count, err := strconv.Atoi(c.Query("count"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Furniture.UpdateTradePointFurniture(id, furniture_id, count)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete a furniture association for a trade point (ID in query)
// @Description  Deletes the association of a furniture item with a specific trade point (ID passed in the query).
// @Tags         furniture
// @Accept       json
// @Produce      json
// @Param        id  query  int    true  "Association ID to delete"
// @Success      200 {object}  statusResponse "Successfully deleted association"
// @Failure      400 {object}  errorResponse "Invalid ID"
// @Failure      500 {object}  errorResponse "Internal server error"
// @Router       /api/furniture/point [delete]
func (h *Handler) DeleteFurnitureForTradePoint(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Furniture.DeleteTradePointFurniture(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{Status: "OK"})
}
