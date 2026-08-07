package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new promotion
// @Security ApiKeyAuth
// @Description Creates a new promotion in the system.
// @Tags promotions
// @Accept json
// @Produce json
// @Param input body models.CreatePromotion true "Promotion information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promotions [post]
func (h *Handler) createPromotion(c *gin.Context) {
	var input models.CreatePromotion

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Promotion.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Promotions
// @Security     ApiKeyAuth
// @Tags         promotions
// @Description  Get all promotions with pagination.
// @ID          get-all-promotions
// @Accept       json
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param        pageSize  query    int     false  "Number of promotions per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "Promotion name to filter by (optional)"
// @Success      200        {object} models.PromotionPage "Paginated list of promotions"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/promotions [get]
func (h *Handler) getAllPromotions(c *gin.Context) {
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

	page, err := h.services.Promotion.GetAll(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Promotion status
// @Security     ApiKeyAuth
// @Tags         promotions
// @Description  Change Promotion status.
// @ID          change-promotion-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Promotion id"
// @Success      200    {object} int "Success"  "Promotion status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "Promotion not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/promotions/{id}/lock [get]
func (h *Handler) changePromotionStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Promotion.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Promotion information
// @Security     ApiKeyAuth
// @Tags         promotions
// @Description  Update promotion.
// @ID          update-promotion
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Promotion id"
// @Param        input  body   models.CreatePromotion  true  "Promotion information"
// @Success      200    {object}  int "Success" "Promotion updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Promotion not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/promotions/{id} [put]
func (h *Handler) updatePromotion(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreatePromotion
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Promotion.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Promotion Name and Time
// @Security     ApiKeyAuth
// @Tags         promotions
// @Description  Обновить название и срок действия акции.
// @ID           update-promotion-name-and-time
// @Accept       json
// @Produce      json
// @Param        id     path   int                        true  "Promotion ID"
// @Param        input  body   models.UpdateDiscountInfo  true  "New promotion name and time"
// @Success      200    {object}  statusResponse          "Success"
// @Failure      400    {object}  errorResponse           "Invalid request"
// @Failure      401    {object}  errorResponse           "Unauthorized"
// @Failure      403    {object}  errorResponse           "Forbidden"
// @Failure      404    {object}  errorResponse           "Promotion not found"
// @Failure      500    {object}  errorResponse           "Internal server error"
// @Router       /api/promotions/info [put]
func (h *Handler) updatePromotionNameAndTime(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.UpdateDiscountInfo
	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err = h.services.Promotion.UpdateNameAndTime(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a promotion
// @Security ApiKeyAuth
// @Description Deletes a promotion from the system.
// @Tags promotions
// @Produce json
// @Param id path int true "Promotion ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promotions/{id} [delete]
func (h *Handler) deletePromotion(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Promotion.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
