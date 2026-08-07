package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new promo
// @Security ApiKeyAuth
// @Description Creates a new promo in the system.
// @Tags promo
// @Accept json
// @Produce json
// @Param input body models.UpdatePromo true "Promo info"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promo [post]
func (h *Handler) createPromo(c *gin.Context) {
	var input models.UpdatePromo

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Promo.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Get promo page
// @Security ApiKeyAuth
// @Description Returns a paginated list of promos.
// @Tags promo
// @Produce json
// @Param language header string false "Language"
// @Param pageSize query int true "Page size"
// @Param pageNumber query int true "Page number"
// @Param name query string false "Promo name filter"
// @Success 200 {object} []models.Promo
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promo/page [get]
func (h *Handler) getPromoPage(c *gin.Context) {
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

	page, err := h.services.Promo.GetPage(pageSize, pageNumber, name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary Update promo information
// @Security ApiKeyAuth
// @Description Updates promo information by ID.
// @Tags promo
// @Accept json
// @Produce json
// @Param id path int true "Promo ID"
// @Param input body models.UpdatePromo true "Promo update info"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promo/{id} [put]
func (h *Handler) updatePromoInformation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.UpdatePromo
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Promo.Update(id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Change promo status
// @Security ApiKeyAuth
// @Description Toggles the status of a promo (active/inactive).
// @Tags promo
// @Produce json
// @Param id path int true "Promo ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promo/{id}/lock [get]
func (h *Handler) changePromoStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Promo.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a promo
// @Security ApiKeyAuth
// @Description Deletes a promo from the system.
// @Tags promo
// @Produce json
// @Param id path int true "Promo ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/promo/{id} [delete]
func (h *Handler) deletePromo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Promo.Delete(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
