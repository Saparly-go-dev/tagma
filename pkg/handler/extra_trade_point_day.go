package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary      Update ExtraTradePointDay information
// @Security     ApiKeyAuth
// @Tags         extra
// @Description  Update.
// @ID          update-point-extra-day
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Point id"
// @Param        input  body   []int  true  "ExtraTradePointDay information"
// @Success      200    {object}  int "Success" "ExtraTradePointDay updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "ExtraTradePointDay not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/extra/{id} [put]
func (h *Handler) updateExtraTradePointDay(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input []int
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.ExtraTradePointDay.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All ExtraTradePointDays
// @Security     ApiKeyAuth
// @Tags         extra
// @Description  Get all trade point days.
// @ID          get-extra-days
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Point id"
// @Success      200        {object} []int "days list"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/extra/{id} [get]
func (h *Handler) getExtraTradePointDay(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	page := h.services.ExtraTradePointDay.Get(Id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}
