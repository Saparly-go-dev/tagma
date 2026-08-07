package handler

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Create a new user
// @Security ApiKeyAuth
// @Description Creates a new user in the system.
// @Tags users
// @Accept json
// @Produce json
// @Param input body tagma.CreateUser true "User information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/users [post]
func (h *Handler) createUser(c *gin.Context) {
	var input tagma.CreateUser

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.User.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Users
// @Security     ApiKeyAuth
// @Tags         users
// @Description  Get all users with pagination.
// @ID          get-all-users
// @Accept       json
// @Produce      json
// @Param 		 language header string false "Language (optional)"
// @Param        pageSize  query    int     false  "Number of users per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Param        name       query    string  false  "User name to filter by (optional)"
// @Param        tip       query    string  false  "User type to filter by (optional)"
// @Success      200        {object} tagma.UserPage "Paginated list of users"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/users [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	// Get the 'language' header
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
	tip := c.Query("tip")

	page, err := h.services.User.GetPage(pageSize, pageNumber, name, tip, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary Change user password
// @Security ApiKeyAuth
// @Description Change user password .
// @Tags users
// @Accept json
// @Produce json
// @Param input body tagma.Change_Password true "User information"
// @Success 200 {integer} int "Success"
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/users/password [put]
func (h *Handler) changeUserPassword(c *gin.Context) {
	var input tagma.Change_Password

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, err := getUserId(c)
	if err != nil {
		return
	}
	role, err := getUserRole(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	if role != "admin" || input.UserId <= 0 {
		input.UserId = userID
	}

	err = h.services.User.ChangePassword(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Change User status
// @Security     ApiKeyAuth
// @Tags         users
// @Description  Change User status.
// @ID          change-user-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "User id"
// @Success      200    {object} int "Success"  "User status updated"
// @Failure      400    {object}  errorResponse "Invalid city ID"
// @Failure      404    {object}  errorResponse "User not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/users/{id}/lock [get]
func (h *Handler) changeUserStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.User.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update User information
// @Security     ApiKeyAuth
// @Tags         users
// @Description  Update user.
// @ID          update-user
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "User id"
// @Param        input  body   tagma.CreateUser  true  "User information"
// @Success      200    {object}  int "Success" "User updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "User not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/users/{id} [put]
func (h *Handler) updateUser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input tagma.CreateUser
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.User.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Delete a user
// @Security ApiKeyAuth
// @Description Deletes a user from the system.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} statusResponse
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /api/users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.User.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
