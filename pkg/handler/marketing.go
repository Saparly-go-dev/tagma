package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary      Get Marketing Main Page
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get marketing main page.
// @ID           get-marketing-main-page
// @Accept       json
// @Produce      json
// @Param        cityId       query    int  false  "City Id"
// @Param        districtId       query    int  false "District Id"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        agentId       query    int  false  "Trade agent Id"
// @Param        volume       query    float64  false  "Product volume"
// @Param        brand       query    int  false  "Product Brand Id"
// @Param        taste       query    string  false  "Product Taste name"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.MarketingMainPage "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/main [get]
func (h *Handler) getMarketingMainPage(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	brand, err := strconv.Atoi(c.Query("brand"))
	if err != nil {
		brand = 0
	}

	volumeIdStr := c.Query("volume")
	volume, err := strconv.ParseFloat(volumeIdStr, 64)
	if err != nil {
		volume = float64(0)
	}

	taste := c.Query("taste")

	result, err := h.services.Marketing.GetMarketingMainPage(cityId, districtId, tradePointId, agentId, brand, volume, taste, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Marketing Main Page
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get marketing main page Excel.
// @ID           get-marketing-main-page-excel
// @Accept       json
// @Produce      json
// @Param        cityId       query    int  false  "City Id"
// @Param        districtId       query    int  false "District Id"
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        agentId       query    int  false  "Trade agent Id"
// @Param        volume       query    float64  false  "Product volume"
// @Param        brand       query    int  false  "Product Brand Id"
// @Param        taste       query    string  false  "Product Taste name"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.MarketingMainPage "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/main/excel [get]
func (h *Handler) getMarketingMainPageExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	cityId, err := strconv.Atoi(c.Query("cityId"))
	if err != nil {
		cityId = 0
	}

	districtId, err := strconv.Atoi(c.Query("districtId"))
	if err != nil {
		districtId = 0
	}

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	agentId, err := strconv.Atoi(c.Query("agentId"))
	if err != nil {
		agentId = 0
	}

	brand, err := strconv.Atoi(c.Query("brand"))
	if err != nil {
		brand = 0
	}

	volumeIdStr := c.Query("volume")
	volume, err := strconv.ParseFloat(volumeIdStr, 64)
	if err != nil {
		volume = float64(0)
	}

	taste := c.Query("taste")

	result, err := h.services.Marketing.MarketingMainPageExport(cityId, districtId, tradePointId, agentId, brand, volume, taste, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=filtrs.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// @Summary      Get Distributin
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get distribution.
// @ID           get-distribution
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.ReadMarketingById "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/dist [get]
func (h *Handler) getDistribution(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	result, err := h.services.Marketing.GetDistribution(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Distributin
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get distribution excel.
// @ID           get-distribution-excel
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.ReadMarketingById "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/dist/excel [get]
func (h *Handler) getDistributionExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	result, err := h.services.Marketing.DistributionExport(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=filtrs.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// @Summary      Get Rasprodano
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get Rasprodano.
// @ID           get-rasprodano
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.ReadMarketingById "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/rasp [get]
func (h *Handler) getRasprodano(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	result, err := h.services.Marketing.GetRasprodano(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Rasprodano
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get Rasprodano.
// @ID           get-rasprodano-excel
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.ReadMarketingById "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/rasp/excel [get]
func (h *Handler) getRasprodanoExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	result, err := h.services.Marketing.RasprodanoExport(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=filtrs.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)
}

// @Summary      Get Wystawlaymost
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get Wystawlaymost.
// @ID           get-wystawlaymost
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.ReadMarketingById "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/wyst [get]
func (h *Handler) getWystawlaymost(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	result, err := h.services.Marketing.GetWystawlaymost(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Wystawlaymost
// @Security     ApiKeyAuth
// @Tags         marketing
// @Description   Get Wystawlaymost.
// @ID           get-wystawlaymost-excel
// @Accept       json
// @Produce      json
// @Param        tradePointId       query    int  false  "Trade Point Id"
// @Param        language  header    string  false  "Language (optional)"
// @Success      200        {object} []models.ReadMarketingById "result"
// @Failure      400        {object} errorResponse "Invalid request parameters"
// @Failure      500        {object} errorResponse "Internal server error"
// @Router       /api/marketing/wyst/excel [get]
func (h *Handler) getWystawlaymostExcel(c *gin.Context) {
	// Get the 'language' header
	language := c.GetHeader("language")

	tradePointId, err := strconv.Atoi(c.Query("tradePointId"))
	if err != nil {
		tradePointId = 0
	}

	result, err := h.services.Marketing.WystawlaymostExport(tradePointId, language)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Set headers to indicate a downloadable Excel file
	c.Header("Content-Disposition", "attachment; filename=filtrs.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", strconv.Itoa(len(result)))

	// Write the Excel file to the response
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", result)

}
