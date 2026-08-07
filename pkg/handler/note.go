package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary Create a new note
// @Security ApiKeyAuth
// @Description Creates a new note in the system.
// @Tags notes
// @Accept json
// @Produce json
// @Param input body models.CreateNote true "Note information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/notes [post]
func (h *Handler) createNote(c *gin.Context) {
	var input models.CreateNote

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Note.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Notes
// @Security     ApiKeyAuth
// @Tags         notes
// @Description  Get all notes with pagination.
// @ID          get-all-notes
// @Accept       json
// @Produce      json
// @Param        pageSize  query    int     false  "Number of notes per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Success      200        {object} models.NotePage "Paginated list of notes"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/notes [get]
func (h *Handler) getAllNotes(c *gin.Context) {
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

	page, err := h.services.Note.GetAll(pageSize, pageNumber)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Update Note information
// @Security     ApiKeyAuth
// @Tags         notes
// @Description  Update note.
// @ID          update-note
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Note id"
// @Param        input  body   models.CreateNote  true  "Note information"
// @Success      200    {object}  int "Success" "Note updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Note not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/notes/{id} [put]
func (h *Handler) updateNote(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateNote
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Note.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a note
// @Security ApiKeyAuth
// @Description Deletes a note from the system.
// @Tags notes
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/notes/{id} [delete]
func (h *Handler) deleteNote(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Note.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Create a new point note
// @Security ApiKeyAuth
// @Description Creates a new point note in the system.
// @Tags notes
// @Accept json
// @Produce json
// @Param        trade_point_id  query    int     false  "Trade Point id" Minimum : 1
// @Param        description       query    string  false  "Note description"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/notes/points [get]
func (h *Handler) createPointNote12(c *gin.Context) {
	var input models.CreatePointNote
	fmt.Println("bashlady")

	tradePointId, err := strconv.Atoi(c.Query("trade_point_id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	description := c.Query("description")
	input.TradePointId = tradePointId
	input.Description = description

	fmt.Println(input)

	err = h.services.Note.CreatePointNote(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Point Notes
// @Security     ApiKeyAuth
// @Tags         notes
// @Description  Get all point notes with pagination.
// @ID          get-all-point-notes
// @Accept       json
// @Produce      json
// @Param 		 language header string false "Language (optional)"
// @Param        tradePointId  query    int     false  "TradePoint Id"
// @Param        code       query    string  false  "Trade Point code"
// @Param        pageSize  query    int     false  "Number of pointnotes per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Success      200        {object} models.PointNotePage "Paginated list of notes"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/notes/point [get]
func (h *Handler) getAllPointNotes(c *gin.Context) {
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

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	code := c.Query("code")

	page, err := h.services.Note.GetAllPointNote(pageSize, pageNumber, tradePointId, language, code)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Update Point Note information
// @Security     ApiKeyAuth
// @Tags         notes
// @Description  Update point note.
// @ID          update-point-note
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "PointNote id"
// @Param        input  body   models.CreatePointNote  true  "Note information"
// @Success      200    {object}  int "Success" "Note updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Note not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/notes/point/{id} [put]
func (h *Handler) updatePointNote(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreatePointNote
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Note.UpdatePointNote(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a point note
// @Security ApiKeyAuth
// @Description Deletes a point note from the system.
// @Tags notes
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/notes/point/{id} [delete]
func (h *Handler) deletePointNote(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Note.DeletePointNote(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
