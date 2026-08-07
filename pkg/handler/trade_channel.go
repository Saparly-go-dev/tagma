package handler

import (
	"net/http"
	"strconv"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
)

// @Summary		Create Trade_Channel_Type
// @Security ApiKeyAuth
// @Tags			channel
// @Description	create channel type
// @ID				create-channel-type
// @Accept			json
// @Produce		json
// @Param			input	body	models.ChannelObject	true	"type info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/channels/types [post]
func (h *Handler) createChannelType(c *gin.Context) {
	var input models.ChannelObject

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Trade_Channel.CreateType(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Trade Channel Types
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get all channel types.
// @ID          get-all-channel-type
// @Produce      json
// @Success      200    {array}   []models.ChannelType "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/types [get]
func (h *Handler) GetChannelTypes(c *gin.Context) {

	page, err := h.services.Trade_Channel.GetAllTypes()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Delete Trade Channel Type
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Delete a trade channel type.
// @ID          delete-channel-type
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel Type ID"
// @Success      200    {object}  int "Success"  "Trade Channel deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/types/{id} [delete]
func (h *Handler) deleteChannelType(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.DeleteTypes(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Trade Channel type
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Update a channel type.
// @ID          update-channel-type
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel ID"
// @Param        input  body   models.ChannelObject  true  "Trade Channel Type information"
// @Success      200    {object}  int "Success" "Channel updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/types/{id} [put]
func (h *Handler) updateChannelType(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.ChannelObject
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.UpdateTypes(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary		Create Trade_Channel_Structure
// @Security ApiKeyAuth
// @Tags			channel
// @Description	create channel structure
// @ID				create-channel-structure
// @Accept			json
// @Produce		json
// @Param			input	body	models.ChannelObject	true	"structure info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/channels/structures [post]
func (h *Handler) createChannelStructure(c *gin.Context) {
	var input models.ChannelObject

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Trade_Channel.CreateStructure(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Trade Channel Structures
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get all channel structures.
// @ID          get-all-channel-structure
// @Produce      json
// @Success      200    {array}   []models.ChannelStructure "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/structures [get]
func (h *Handler) GetChannelStructures(c *gin.Context) {

	page, err := h.services.Trade_Channel.GetAllStructure()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Delete Trade Channel Structure
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Delete a trade channel structure.
// @ID          delete-channel-structure
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel Structure ID"
// @Success      200    {object}  int "Success"  "Trade Channel deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/structures/{id} [delete]
func (h *Handler) deleteChannelStructure(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.DeleteStructure(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Trade Channel structure
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Update a channel structure.
// @ID          update-channel-structure
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel ID"
// @Param        input  body   models.ChannelObject  true  "Trade Channel Structure information"
// @Success      200    {object}  int "Success" "Channel updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/structures/{id} [put]
func (h *Handler) updateChannelStructure(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.ChannelObject
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.UpdateStructure(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary		Create Trade_Channel_Size
// @Security ApiKeyAuth
// @Tags			channel
// @Description	create channel size
// @ID				create-channel-size
// @Accept			json
// @Produce		json
// @Param			input	body	models.ChannelObject	true	"size info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/channels/sizes [post]
func (h *Handler) createChannelSize(c *gin.Context) {
	var input models.ChannelObject

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Trade_Channel.CreateSize(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Trade Channel Sizes
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get all channel sizes.
// @ID          get-all-channel-size
// @Produce      json
// @Success      200    {array}   []models.ChannelSize "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/sizes [get]
func (h *Handler) getChannelSizes(c *gin.Context) {

	page, err := h.services.Trade_Channel.GetAllSizes()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Delete Trade Channel Size
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Delete a trade channel size.
// @ID          delete-channel-size
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel Size ID"
// @Success      200    {object}  int "Success"  "Trade Channel deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/sizes/{id} [delete]
func (h *Handler) deleteChannelSize(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.DeleteSize(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Trade Channel size
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Update a channel size.
// @ID          update-channel-size
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel ID"
// @Param        input  body   models.ChannelObject  true  "Trade Channel Size information"
// @Success      200    {object}  int "Success" "Channel updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/sizes/{id} [put]
func (h *Handler) updateChannelSize(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.ChannelObject
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.UpdateSize(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary		Create Trade_Channel_Management
// @Security ApiKeyAuth
// @Tags			channel
// @Description	create channel management
// @ID				create-channel-management
// @Accept			json
// @Produce		json
// @Param			input	body	models.ChannelObject	true	"management info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/channels/managements [post]
func (h *Handler) createChannelManagement(c *gin.Context) {
	var input models.ChannelObject

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Trade_Channel.CreateManagement(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Trade Channel Managements
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get all channel managements.
// @ID          get-all-channel-management
// @Produce      json
// @Success      200    {array}   []models.ChannelManagement "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/managements [get]
func (h *Handler) getChannelManagements(c *gin.Context) {

	page, err := h.services.Trade_Channel.GetAllManagements()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Delete Trade Channel Management
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Delete a trade channel management.
// @ID          delete-channel-management
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel Management ID"
// @Success      200    {object}  int "Success"  "Trade Channel deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/managements/{id} [delete]
func (h *Handler) deleteChannelManagement(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.DeleteManagement(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Trade Channel management
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Update a channel management.
// @ID          update-channel-management
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel ID"
// @Param        input  body   models.ChannelObject  true  "Trade Channel Management information"
// @Success      200    {object}  int "Success" "Channel updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/managements/{id} [put]
func (h *Handler) updateChannelManagement(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.ChannelObject
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.UpdateManagement(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary		Create Trade_Channel
// @Security ApiKeyAuth
// @Tags			channel
// @Description	create channel
// @ID				create-channel
// @Accept			json
// @Produce		json
// @Param			input	body	models.CreateChannel	true	"channel info"
// @Success		200 	{integer}	int "Success"
// @Failure		400 	404			{object}	errorResponse
// @Failure		500		{object}	errorResponse
// @Failure		default	{object}	errorResponse
// @Router			/api/channels [post]
func (h *Handler) createChannel(c *gin.Context) {
	var input models.CreateChannel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Trade_Channel.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Trade Channels
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get all channels.
// @ID          get-all-channels
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Param		 typeId query int false "channel Type Id "
// @Param        pageSize  query    int     false  "Number of trade channels per page"  Minimum: 1
// @Param        pageNumber query    int     false  "Page number (starts from 1)" Minimum: 1
// @Success      200    {array}   models.TradeChannelPage "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels [get]
func (h *Handler) getAllChannels(c *gin.Context) {
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

	typeId, err := strconv.Atoi(c.Query("typeId"))
	if err != nil {
		typeId = 0
	}

	page, err := h.services.Trade_Channel.GetChannelPage(pageSize, pageNumber, typeId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Change Trade Channel status
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Change Trade Channel status.
// @ID          change-trade-channel-status
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "TradeChannel ID"
// @Success      200    {object}  int "Success"  "Trade Channel updated"
// @Failure      400    {object}  errorResponse "Invalid Trade Channel ID"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/{id}/lock [get]
func (h *Handler) statusTradeChannel(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Trade_Channel.ChangeStatus(id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Update Trade Channel
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Update a channel.
// @ID          update-channel
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel ID"
// @Param        input  body   models.CreateChannel  true  "Trade Channel information"
// @Success      200    {object}  int "Success" "Channel updated"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/{id} [put]
func (h *Handler) updateChannel(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var input models.CreateChannel
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.Update(Id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete Trade Channel
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Delete a trade channel.
// @ID          delete-channel
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Trade Channel ID"
// @Success      200    {object}  int "Success"  "Trade Channel deleted"
// @Failure      400    {object}  errorResponse "Invalid request"
// @Failure      401    {object}  errorResponse "Unauthorized"
// @Failure      403    {object}  errorResponse "Forbidden"
// @Failure      404    {object}  errorResponse "Trade Channel not found"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/{id} [delete]
func (h *Handler) deleteChannel(c *gin.Context) {
	Id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.services.Trade_Channel.Delete(Id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get All Trade types
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get channels types.
// @ID          get-channels-types
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   models.SelectObject "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/type [get]
func (h *Handler) getTypesForChannels(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	page, err := h.services.Trade_Channel.GetTypeForChannel(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get All Trade structure
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get channels structure.
// @ID          get-channels-structure
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   models.SelectObject "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/structure [get]
func (h *Handler) getStructuresForChannels(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	page, err := h.services.Trade_Channel.GetStructureForChannel(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get All Trade size
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get channels size.
// @ID          get-channels-size
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   models.SelectObject "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/size [get]
func (h *Handler) getSizesForChannels(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	page, err := h.services.Trade_Channel.GetSizeForChannel(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}

// @Summary      Get All Trade management
// @Security     ApiKeyAuth
// @Tags         channel
// @Description  Get channels management.
// @ID          get-channels-management
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {array}   models.SelectObject "channels"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/channels/management [get]
func (h *Handler) getManagementsForChannels(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	page, err := h.services.Trade_Channel.GetManagementForChannel(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, page)
}
