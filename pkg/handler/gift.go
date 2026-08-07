package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create a new gift
// @Security ApiKeyAuth
// @Description Creates a new gift in the system.
// @Tags gifts
// @Accept json
// @Produce json
// @Param input body models.CreateGift true "Gift information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/gifts [post]
func (h *Handler) createGift(c *gin.Context) {
	var input models.CreateGift

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Gift.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Gifts
// @Security     ApiKeyAuth
// @Tags         gifts
// @Description  Get all gifts with pagination.
// @ID          get-all-gifts
// @Accept       json
// @Produce      json
// @Param 		 language header string false "Language (optional)"
// @Param        pageSize  query    int     false  "Number of gifts per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "Gift name to filter by (optional)"
// @Success      200        {object} models.GiftPage "Paginated list of gifts"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/gifts [get]
func (h *Handler) getAllGifts(c *gin.Context) {
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

	page, err := h.services.Gift.GetPage(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Gift status
// @Security     ApiKeyAuth
// @Tags         gifts
// @Description  Change Gift status.
// @ID          change-gift-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Gift id"
// @Success      200    {object} int "Success"  "Gift status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "Gift not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/gifts/{id}/lock [get]
func (h *Handler) changeGiftStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Gift.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get gift by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific gift by its ID.
// @Tags gifts
// @Produce json
// @Param id path int true "Gift ID"
// @Success 200 {object} models.Gift "Gift details"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/gifts/{id} [get]
func (h *Handler) getGiftById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	gift, err := h.services.Gift.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gift)
}

// @Summary      Update Gift information
// @Security     ApiKeyAuth
// @Tags         gifts
// @Description  Update gift.
// @ID          update-gift
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Gift id"
// @Param        input  body   models.CreateGift  true  "Gift information"
// @Success      200    {object}  int "Success" "Gift updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Gift not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/gifts/{id} [put]
func (h *Handler) updateGift(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateGift
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Gift.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a gift
// @Security ApiKeyAuth
// @Description Deletes a gift from the system.
// @Tags gifts
// @Produce json
// @Param id path int true "Gift ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/gifts/{id} [delete]
func (h *Handler) deleteGift(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Gift.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Gift Types
// @Security     ApiKeyAuth
// @Tags         gifts
// @Description  Get gift types.
// @ID          get-gift-types
// @Accept       json
// @Produce      json
// @Param 		 language header string false "Language (optional)"
// @Success      200        {object} []models.SelectObject "Paginated list of gifts"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/gifts/types [get]
func (h *Handler) getGiftTypes(c *gin.Context) {
	language := c.GetHeader("language")

	page, err := h.services.Gift.GetGiftTypes(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}
