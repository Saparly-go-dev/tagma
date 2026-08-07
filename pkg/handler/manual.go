package handler

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary      Create Manual
// @Security     ApiKeyAuth
// @Tags         manual
// @Description  create manual information
// @ID           create_manual
// @Produce      json
// @Param        input body models.CreateManual true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/manual/create [post]
func (h *Handler) createManual(c *gin.Context) {

	var input models.CreateManual

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.services.Manual.Create(input)

	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Get Manual
// @Security     ApiKeyAuth
// @Tags         manual
// @Description  get manual information
// @ID           get_maual
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {object}  models.ReadManual "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/manual/get [get]
func (h *Handler) getManual(c *gin.Context) {
	language := c.GetHeader("language")

	result, err := h.services.Manual.Get(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Manual For Mobile
// @Security     ApiKeyAuth
// @Tags         manual
// @Description  get manual information For Mobile
// @ID           get_maual_mobile
// @Produce      json
// @Param        language  header    string  false  "Language (optional)"
// @Success      200    {object}  models.ReadManualForMobile "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/manual/mobile [get]
func (h *Handler) getManualForMobile(c *gin.Context) {
	language := c.GetHeader("language")

	result, err := h.services.Manual.GetForMobile(language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Update Manual
// @Security     ApiKeyAuth
// @Tags         manual
// @Description  update manual information
// @ID           update
// @Produce      json
// @Param        manualId       query    int  false  "Manual Id"
// @Param        input body models.CreateManual true "list of order"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/manual/update [post]
func (h *Handler) updateManual(c *gin.Context) {

	manualId, err := strconv.Atoi(c.Query("manualId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var input models.CreateManual

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.Manual.UpdateManual(manualId, input)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary      Delete Manual Image
// @Security     ApiKeyAuth
// @Tags         manual
// @Description  delete manual image
// @ID           delete_manual_image
// @Produce      json
// @Param        ImageId       query    int  false  "Image Id"
// @Success      200    {object}  int "Success"  "Order Saved"
// @Failure      500    {object}  errorResponse "Internal server error"
// @Router       /api/manual/delete [delete]
func (h *Handler) deleteManualImage(c *gin.Context) {
	imageId, err := strconv.Atoi(c.Query("ImageId"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err1 := h.services.Manual.DeleteImage(imageId)

	if err1 != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}

// @Summary Save  Image
// @Description Save Manual image.
// @ID           save_manual_image
// @Security ApiKeyAuth
// @Tags        manual
// @Accept      multipart/form-data
// @Produce     json
// @Param       manualId  formData  int   true      "Manual Id"
// @Param       image formData  file      false     "Image 1 (optional)"
// @Success    200     {object}    int "Success"   "Success message with ID"
// @Failure    400     {object}   errorResponse          "Bad request (invalid data or image upload error)"
// @Failure    500     {object}   errorResponse          "Internal server error"
// @Router      /api/manual/save [post]
func (h *Handler) saveManualImage(c *gin.Context) {
	// Retrieve JSON data from form input (assuming 'input' field contains JSON string)
	manualId, err := strconv.Atoi(c.PostForm("manualId"))

	// Handle image1 upload
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Image upload error: %s", err.Error()))
		return
	}

	if err := h.services.Manual.SaveImage(manualId, file); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "Ok",
	})
}
