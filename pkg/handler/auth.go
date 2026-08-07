package handler

import (
	"net/http"

	"github.com/Saparly-go-dev/tagma"
	"github.com/gin-gonic/gin"
)

// @Summary		SignUp
// @Tags			auth
// @Description	create account
// @ID				create-account
// @Accept			json
// @Produce		json
// @Param			input	body	tagma.User	true	"account info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/auth/sign-up [post]
func (h *Handler) singUp(c *gin.Context) {
	var input tagma.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary		SignIp
// @Tags			auth
// @Description	login
// @ID				login
// @Accept			json
// @Produce		json
// @Param			input	body		signInInput	true	"credentails"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/auth/sign-in [post]
func (h *Handler) singIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	tip, err := h.services.Authorization.GetUserType(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, "invalid credentials")
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
		"type":  tip,
	})

}
