package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new discount
// @Security ApiKeyAuth
// @Description Creates a new discount in the system.
// @Tags discounts
// @Accept json
// @Produce json
// @Param input body models.CreateDiscount true "Discount information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/discounts [post]
func (h *Handler) createDiscount(c *gin.Context) {
	var input models.CreateDiscount

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Discount.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Discounts
// @Security     ApiKeyAuth
// @Tags         discounts
// @Description  Get all discounts with pagination.
// @ID          get-all-discounts
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        pageSize  query    int     false  "Number of discounts per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "Discount name to filter by (optional)"
// @Success      200        {object} models.DiscountPage "Paginated list of discounts"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/discounts [get]
func (h *Handler) getAllDiscounts(c *gin.Context) {
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

	page, err := h.services.Discount.GetAll(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Discount status
// @Security     ApiKeyAuth
// @Tags         discounts
// @Description  Change Discount status.
// @ID          change-discount-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Discount id"
// @Success      200    {object} int "Success"  "Discount status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "Discount not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/discounts/{id}/lock [get]
func (h *Handler) changeDiscountStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Discount.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Discount information
// @Security     ApiKeyAuth
// @Tags         discounts
// @Description  Update discount.
// @ID          update-discount
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Discount id"
// @Param        input  body   models.CreateDiscount  true  "Discount information"
// @Success      200    {object}  int "Success" "Discount updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Discount not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/discounts/{id} [put]
func (h *Handler) updateDiscount(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateDiscount
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Discount.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a discount
// @Security ApiKeyAuth
// @Description Deletes a discount from the system.
// @Tags discounts
// @Produce json
// @Param id path int true "Discount ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/discounts/{id} [delete]
func (h *Handler) deleteDiscount(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Discount.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Discount Name and Time
// @Security     ApiKeyAuth
// @Tags         discounts
// @Description  Обновить название и срок действия скидки.
// @ID           update-discount-name-and-time
// @Accept       json
// @Produce      json
// @Param        id     query   int                        true  "Discount ID"
// @Param        input  body   models.UpdateDiscountInfo  true  "New discount name and time"
// @Success      200    {object}  statusResponse          "Success"
// @Failure      400    {object}  errorResponse           "Invalid request"
// @Failure      401    {object}  errorResponse           "Unauthorized"
// @Failure      403    {object}  errorResponse           "Forbidden"
// @Failure      404    {object}  errorResponse           "Discount not found"
// @Failure      500    {object}  errorResponse           "Internal server error"
// @Router       /api/discounts/info [put]
func (h *Handler) updateDiscountNameAndTime(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.UpdateDiscountInfo
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	fmt.Print(id)
	if err = h.services.Discount.UpdateNameAndDate(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get TradePoints by Trade agent
// @Security     ApiKeyAuth
// @Tags         discounts
// @Description  Get trade points by agent.
// @ID          get-trade-point-by-agent
// @Accept       json
// @Produce      json
// @Param        agentId query    int     false  "Trade Agent Id" Minimum: 1
// @Success      200        {object} []models.SelectObject "List of Trade Points"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/discounts/get-points-by-agent [get]
func (h *Handler) getTradePointByAgent(c *gin.Context) {
	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.services.Discount.GetTradePointByTradeAgent(agentId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}
