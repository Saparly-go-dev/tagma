package handler

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary		Create Contact
// @Security ApiKeyAuth
// @Tags			contacts
// @Description	create contact
// @ID				create-contact
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateContact	true	"contact info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/contacts [post]
func (h *Handler) createContact(c *gin.Context) {
	var input models.CreateContact

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Contact.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get all contacts
// @Security ApiKeyAuth
// @Description Retrieves all contacts for a specific point and language.
// @Tags        contacts
// @Accept      json
// @Produce     json
// @Param       id  query    int     true      "ID of the Trade point"
// @Param       language  header   string  false     "Language (optional)"
// @Success    200       {object}  interface{}  "List of contacts"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/contacts [get]
func (h *Handler) getAllContacts(c *gin.Context) {
	language := c.GetHeader("language")

	pointId, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	page, err := h.services.Contact.GetAll(pointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary Get contact by ID
// @Security ApiKeyAuth
// @Description Retrieves a specific contact by its ID.
// @Tags        contacts
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Contact ID"
// @Success    200       {object}  interface{}  "Contact details"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/contacts/{id} [get]
func (h *Handler) getContactById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	city, err := h.services.Contact.GetById(id)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, city)
}

// @Summary Update a contact
// @Security ApiKeyAuth
// @Description Updates an existing contact with provided information.
// @Tags        contacts
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Contact ID"
// @Param       input    body     models.CreateContact  true  "Contact update data"
// @Success    200       {object}  statusResponse   "Update status"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/contacts/{id} [put]
func (h *Handler) updateContact(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateContact
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Contact.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a contact
// @Security ApiKeyAuth
// @Description Deletes a contact by its ID.
// @Tags        contacts
// @Accept      json
// @Produce     json
// @Param       id       path     int     true      "Contact ID"
// @Success    200       {object}  statusResponse   "Deletion status"
// @Failure    400       {object}   errorResponse "Bad request"
// @Failure    500       {object}   errorResponse "Internal server error"
// @Router      /api/contacts/{id} [delete]
func (h *Handler) deleteContact(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Contact.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get  Kinds for Contacts
// @Security     ApiKeyAuth
// @Tags         contacts
// @Description  Get kind for contacts.
// @ID           get-kind-for-contact
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Kind name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/contacts/kinds [get]
func (h *Handler) GetKindForContacts(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.Contact.GetKinds(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get  Post for Contacts
// @Security     ApiKeyAuth
// @Tags         contacts
// @Description  Get post for contacts.
// @ID           get-post-for-contact
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Post name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/contacts/posts [get]
func (h *Handler) GetPostForContacts(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.Contact.GetPosts(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get  Uprs for Contacts
// @Security     ApiKeyAuth
// @Tags         contacts
// @Description  Get upr for contacts.
// @ID           get-upr-for-contact
// @Accept       json
// @Produce      json
// @Param        name       query    string  false  "Uprs name to filter by (optional)"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.SelectObject "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/contacts/uprs [get]
func (h *Handler) GetUprsForContacts(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	name := c.Query("name")

	page, err := h.services.Contact.GetUprs(name, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}
