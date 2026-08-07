package handler

import (
	_ "github.com/Saparly-go-dev/tagma/docs"
	"github.com/Saparly-go-dev/tagma/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(requestResponseLogger(), gin.Recovery())

	router.Static("/images", "./images")

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", h.userIdentity, h.requireRoles("admin"), h.singUp)
		}

		auth.POST("/sign-in", h.singIn)

		agents := api.Group("/agents", h.userIdentity, h.requireRoles("admin", "agent"))
		{
			agents.POST("/", h.createTradeAgent)
			agents.GET("/", h.getAllTradeAgents)
			agents.GET("/:id/lock", h.statusTradeAgent)
			agents.GET("/:id", h.getTradeAgentById)
			agents.PUT("/:id", h.updateTradeAgent)
			agents.DELETE("/:id", h.deleteTradeAgent)
			agents.GET("/drivers", h.getDriverforAgents)
			agents.GET("/ekspeditors", h.getEkspeditorforAgents)
			agents.GET("/excel", h.GetAgentExcel)
			agents.GET("/debt", h.getAgentDebt)
			agents.GET("/merchendising", h.GetMerchandiserInformationForAgent)
			agents.GET("/factory/csv", h.getFactoryCSV)
			agents.GET("/revenue", h.getTradeAgentRevenue)
		}

		logs := api.Group("/logs", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			//logs.POST("/", h.saveLog)
			logs.GET("/", h.getAllLogs)
			logs.GET("/types", h.getLogTypes)
		}

		points := api.Group("/points", h.userIdentity, h.requireRoles("admin", "agent", "ekspeditor", "merchandiser"))
		{
			points.POST("/", h.createTradePoint)
			points.GET("/", h.getAllTradePoints)
			points.GET("/:id/lock", h.statusTradePoint)
			points.GET("/:id", h.getTradePointById)
			points.PUT("/:id", h.updateTradePoint)
			points.PUT("/change_agent", h.updateTradePointsAgent)
			points.DELETE("/:id", h.deleteTradePoint)
			points.GET("/cities", h.getCitiesForPoints)
			points.GET("/districts", h.getDistrictsForCity)
			points.GET("/agents", h.getAgentsForPoints)
			points.GET("/channels", h.getChannelsForTradePoints)
			points.GET("/:id/general", h.getGeneralInformationForTradePoint)
			points.DELETE("/images/:id", h.deleteImage)
			points.POST("/resave", h.ResavePointImage)
			points.GET("/images/:id", h.GetPointImages)
			points.GET("/excel", h.getTradePointExcel)
			points.GET("/order", h.getInformationAboutOrder)
			points.GET("/product", h.getInformationAboutProducts)
			points.GET("/sale/history", h.getTradePointSaleHistoryForMonth)

			debts := points.Group("debts", h.userIdentity)
			{
				debts.POST("/create", h.CreateTradePointDebt)
				debts.GET("/mobile", h.GetTradePointDebtsForMobile)
			}
		}

		contacts := api.Group("/contacts", h.userIdentity, h.requireRoles("admin"))
		{
			contacts.POST("/", h.createContact)
			contacts.GET("/", h.getAllContacts)
			contacts.GET("/:id", h.getContactById)
			contacts.PUT("/:id", h.updateContact)
			contacts.DELETE("/:id", h.deleteContact)
			contacts.GET("/kinds", h.GetKindForContacts)
			contacts.GET("/posts", h.GetPostForContacts)
			contacts.GET("/uprs", h.GetUprsForContacts)
		}

		cities := api.Group("/cities", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			cities.POST("/", h.createCity)
			cities.GET("/", h.getAllCities)
			cities.GET("/:id/lock", h.statusCity)
			cities.GET("/:id", h.getCityById)
			cities.PUT("/:id", h.updateCity)
			cities.DELETE("/:id", h.deleteCity)

		}

		channels := api.Group("/channels", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			channels.POST("/types", h.createChannelType)
			channels.GET("/types", h.GetChannelTypes)
			channels.DELETE("/types/:id", h.deleteChannelType)
			channels.PUT("/types/:id", h.updateChannelType)
			channels.POST("/structures", h.createChannelStructure)
			channels.GET("/structures", h.GetChannelStructures)
			channels.DELETE("/structures/:id", h.deleteChannelStructure)
			channels.PUT("/structures/:id", h.updateChannelStructure)
			channels.POST("/sizes", h.createChannelSize)
			channels.GET("/sizes", h.getChannelSizes)
			channels.DELETE("/sizes/:id", h.deleteChannelSize)
			channels.PUT("/sizes/:id", h.updateChannelSize)
			channels.POST("/managements", h.createChannelManagement)
			channels.GET("/managements", h.getChannelManagements)
			channels.DELETE("/managements/:id", h.deleteChannelManagement)
			channels.PUT("/managements/:id", h.updateChannelManagement)
			channels.POST("/", h.createChannel)
			channels.GET("/", h.getAllChannels)
			channels.GET("/:id/lock", h.statusTradeChannel)
			//hannels.GET("/:id", h.getChannelById)
			channels.PUT("/:id", h.updateChannel)
			channels.DELETE("/:id", h.deleteChannel)
			channels.GET("/type", h.getTypesForChannels)
			channels.GET("/structure", h.getStructuresForChannels)
			channels.GET("/size", h.getSizesForChannels)
			channels.GET("/management", h.getManagementsForChannels)
		}

		category := api.Group("/categories", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			category.POST("/", h.createCategory)
			category.GET("/", h.getAllCategories)
			category.GET("/:id/lock", h.statusTradeCategory)
			category.GET("/:id", h.getCategoryById)
			category.PUT("/:id", h.updateCategory)
			category.DELETE("/:id", h.deleteCategory)
		}

		product := api.Group("/products", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			types := product.Group("/types")
			{
				types.POST("/", h.createProductType)
				types.GET("/all", h.getProductTypesForSelect)
				types.GET("/", h.getAllProductType)
				types.GET("/:id/lock", h.statusProducttype)
				types.GET("/:id", h.GetProductTypeById)
				types.PUT("/:id", h.updateProductType)
				types.DELETE("/:id", h.deleteProductType)
				types.GET("/volume", h.getProductTypeVolumes)
			}

			product.POST("/", h.createProduct)
			product.GET("/", h.getAllProduct)
			product.GET("/:id/lock", h.statusProduct)
			product.GET("/:id", h.getProductById)
			product.GET("/brands", h.getBrandsForProduct)
			product.PUT("/:id", h.updateProduct)
			product.DELETE("/:id", h.deleteProduct)
			product.GET("/sales", h.getInformationAboutSales)
			product.GET("/daily/excel", h.getDailySalesProductInExcel)
		}

		driver := api.Group("/drivers", h.userIdentity, h.requireRoles("admin"))
		{
			driver.POST("/", h.createDriver)
			driver.GET("/", h.getAllDrivers)
			driver.GET("/:id/lock", h.statusDriver)
			driver.GET("/:id", h.getDriverById)
			driver.PUT("/:id", h.updateDriver)
			driver.DELETE("/:id", h.deleteDriver)
		}

		ekspeditor := api.Group("/ekspeditors", h.userIdentity, h.requireRoles("admin", "ekspeditor", "agent"))
		{
			ekspeditor.POST("/", h.createEkspeditor)
			ekspeditor.GET("/", h.getAllEkspeditors)
			ekspeditor.GET("/:id/lock", h.statusEkspeditor)
			ekspeditor.GET("/:id", h.getEkspeditorById)
			ekspeditor.PUT("/:id", h.updateEkspeditor)
			ekspeditor.DELETE("/:id", h.deleteEkspeditor)
			ekspeditor.GET("/product/all", h.GetProductsCountForEkspeditor)
			ekspeditor.GET("/product/point", h.GetPointProductListForEkspeditor)
			ekspeditor.POST("/save", h.SaveEkspeditorProductList)
			ekspeditor.GET("/points", h.getTradePointsForEkspeditors)
			ekspeditor.POST("/product/delivered", h.ChangeOrderStatus)
			ekspeditor.GET("/list", h.getEkspeditorsProductList)
			ekspeditor.GET("/galyndy", h.GetEkspeditorGalyndy)
			ekspeditor.GET("/galyndy/web", h.getEkspeditorGalyndyForAdmin)
			ekspeditor.DELETE("/delete", h.deleteEkspdata)
			ekspeditor.GET("/product/minus", h.GetProductsCountForEkspeditorWithMinus)
			ekspeditor.GET("/car/galyndy", h.getProductInTheCarWithDifference)
			ekspeditor.GET("/car/galyndy/price", h.getProductPriceInTheCarWithDifference)
			ekspeditor.GET("/car/galyndy/array", h.getProductInTheCarWithDifferenceInArray)
		}

		route := api.Group("/routes", h.userIdentity, h.requireRoles("admin", "agent", "ekspeditor"))
		{
			route.POST("/", h.createRoute)
			route.GET("/", h.getAllRoutes)
			route.GET("/:id/lock", h.statusRoute)
			route.GET("/:id", h.getRouteById)
			route.PUT("/:id", h.updateRoute)
			route.DELETE("/:id", h.deleteRoute)
			route.GET("/days", h.GetDaysForRoute)
			route.GET("/points", h.GetTradePointsForRoute)
			route.GET("/ekspeditors/:agentId", h.GETEkspeditorsForRoute)
			route.GET("/orders/get", h.get_route_orders)
			route.POST("orders/post", h.post_route_order)
			route.GET("/excel", h.excelRouteOrders)

		}

		order := api.Group("/orders", h.userIdentity, h.requireRoles("admin", "agent", "ekspeditor"))
		{
			order.GET("/brands", h.getBrandsForOrders)
			order.GET("/models", h.getModelsForOrder)
			order.GET("/types", h.getTypesForOrders)
			order.GET("/volumes", h.getVolumesForOrder)
			order.GET("/points", h.getTradePointsForOrder)
			order.GET("/ekspeditor/points", h.getTradePointsForOrderFromEkspeditorId)
			order.GET("/points/marketing", h.getTradePointsForMarketing)
			order.GET("/products", h.getProductsForOrders)
			order.GET("/sum", h.getSumForOrder)
			order.POST("/save", h.saveTheOrder)
			order.POST("/marketing", h.marketingForAgent)
			order.GET("/page", h.getOrdersPage)
			order.GET("/excel", h.orderExcel)
			order.GET("/list", h.getOrderListById)
			order.GET("/", h.getOrders)
			order.POST("/close", h.closeTheOrder)
		}

		discount := api.Group("/discounts", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			discount.POST("/", h.createDiscount)
			discount.GET("/", h.getAllDiscounts)
			discount.GET("/:id/lock", h.changeDiscountStatus)
			discount.PUT("/:id", h.updateDiscount)
			discount.DELETE("/:id", h.deleteDiscount)
			discount.PUT("/info", h.updateDiscountNameAndTime)
			discount.GET("get-points-by-agent", h.getTradePointByAgent)
		}

		gift := api.Group("/gifts", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			gift.POST("/", h.createGift)
			gift.GET("/", h.getAllGifts)
			gift.GET("/:id/lock", h.changeGiftStatus)
			gift.GET("/:id", h.getGiftById)
			gift.PUT("/:id", h.updateGift)
			gift.DELETE("/:id", h.deleteGift)
			gift.GET("/types", h.getGiftTypes)
		}

		district := api.Group("districts", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			district.POST("/", h.createDistrict)
			district.GET("/", h.getAllDistricts)
			district.GET("/:id/lock", h.statusDistrict)
			district.GET("/:id", h.getDistrictById)
			district.PUT("/:id", h.updateDistrict)
			district.DELETE("/:id", h.deleteDistrict)
		}

		note := api.Group("/notes", h.userIdentity, h.requireRoles("admin", "agent", "ekspeditor"))
		{
			note.POST("/", h.createNote)
			note.GET("/", h.getAllNotes)
			note.PUT("/:id", h.updateNote)
			note.DELETE("/:id", h.deleteNote)
			note.POST("/points/", h.createPointNote12)
			note.GET("/point/", h.getAllPointNotes)
			note.PUT("/point/:id", h.updatePointNote)
			note.DELETE("/point/:id", h.deletePointNote)
		}

		user := api.Group("/users", h.userIdentity)
		{
			user.PUT("/password", h.changeUserPassword)
			user.POST("/", h.requireRoles("admin"), h.createUser)
			user.GET("/", h.requireRoles("admin"), h.getAllUsers)
			user.GET("/:id/lock", h.requireRoles("admin"), h.changeUserStatus)
			user.PUT("/:id", h.requireRoles("admin"), h.updateUser)
			user.DELETE("/:id", h.requireRoles("admin"), h.deleteUser)
		}

		payment := api.Group("/payments", h.userIdentity, h.requireRoles("admin", "agent", "ekspeditor"))
		{
			payment.GET("/points", h.getTradePointsForPayment)
			payment.GET("ekspeditor/points", h.getTradePointsForPaymentFromEkspeditor)
			payment.POST("/save", h.saveThePayment)
			payment.POST("/ekspeditor/save", h.saveThePaymentFromEkspeditor)
			payment.GET("/types", h.getPaymentTypes)
			payment.GET("/order", h.getPaymentsForOrder)
			payment.GET("/list", h.GetAgentPayments)
			payment.POST("/confirm", h.confirmPayment)
			payment.GET("update", h.updatePayment)
			payment.GET("/debt", h.updateDebt)
			payment.GET("/add_debt", h.addSelfDebtForAgent)
			payment.POST("/add_debt", h.addSelfDebtForAgent)
			payment.GET("/trade-point-debt", h.getTradePointAllDebtSum)
			payment.POST("/create", h.PostPayment)
			payment.GET("/history", h.GetTradePointPaymentsHistory)
			payment.GET("/daily", h.getDailyInformationAboutPayment)
			payment.GET("/trade-points-debt-list", h.getTradePointsDebtList)
		}

		manual := api.Group("/manual", h.userIdentity, h.requireRoles("admin", "agent", "ekspeditor", "merchandiser", "viewer"))
		{
			manual.POST("/create", h.createManual)
			manual.GET("/get", h.getManual)
			manual.POST("/update", h.updateManual)
			manual.DELETE("/delete", h.deleteManualImage)
			manual.POST("/save", h.saveManualImage)
			manual.GET("/mobile", h.getManualForMobile)
		}

		marketing := api.Group("/marketing", h.userIdentity, h.requireRoles("admin", "viewer", "merchandiser"))
		{
			marketing.GET("/main", h.getMarketingMainPage)
			marketing.GET("/dist", h.getDistribution)
			marketing.GET("/rasp", h.getRasprodano)
			marketing.GET("/wyst", h.getWystawlaymost)
			marketing.GET("/main/excel", h.getMarketingMainPageExcel)
			marketing.GET("/dist/excel", h.getDistributionExcel)
			marketing.GET("/rasp/excel", h.getRasprodanoExcel)
			marketing.GET("/wyst/excel", h.getWystawlaymostExcel)
		}

		dealer := api.Group("/dealer", h.userIdentity, h.requireRoles("admin"))
		{
			dealer.POST("/", h.createDealer)
			dealer.GET("/", h.getDealersByPage)
			dealer.PUT("/:id", h.updateDealer)
			dealer.GET("/:id/lock", h.changeDealerStatus)
			dealer.GET("/all", h.getDealersByName)
			dealer.POST("/topup", h.TopUpDealer)

			dealer_product := api.Group("/product", h.userIdentity, h.requireRoles("admin"))
			{
				dealer_product.POST("/", h.createDealerProduct)
				dealer_product.GET("/", h.getAllDealerProduct)
				dealer_product.GET("/:id/lock", h.statusDealerProduct)
				dealer_product.PUT("/:id", h.updateDealerProduct)
				dealer_product.DELETE("/:id", h.deleteDealerProduct)
			}
		}

		dorder := api.Group("/dorders", h.userIdentity, h.requireRoles("admin"))
		{
			dorder.POST("/save", h.saveDealerOrder)
			dorder.POST("/update", h.updateDealerOrder)
			dorder.GET("/:id/lock", h.changeDealerOrderStatus)
			dorder.GET("/", h.getDealerOrderPage)
			dorder.GET("/product", h.getDealerOrderByProducts)
			dorder.GET("/list", h.getDealerOrderListById)
			dorder.POST("/complete", h.completeDealerOrder)
			dorder.DELETE("/", h.deleteDealerOrder)
		}

		promotion := api.Group("/promotions", h.userIdentity, h.requireRoles("admin", "viewer"))
		{
			promotion.POST("/", h.createPromotion)
			promotion.GET("/", h.getAllPromotions)
			promotion.GET("/:id/lock", h.changePromotionStatus)
			promotion.PUT("/:id", h.updatePromotion)
			promotion.DELETE("/:id", h.deletePromotion)
			promotion.PUT("/info", h.updatePromotionNameAndTime)
		}

		extra := api.Group("/extra", h.userIdentity, h.requireRoles("admin", "agent"))
		{
			extra.PUT("/:id", h.updateExtraTradePointDay)
			extra.GET("/:id", h.getExtraTradePointDay)
		}
		extra_order := api.Group("/extra/order", h.userIdentity, h.requireRoles("admin", "agent"))
		{
			extra_order.GET("/agent", h.getTradePointsForExtraOrder)
			extra_order.GET("/items", h.getOrderItemsForExtraOrder)
		}

		plan := api.Group("/plans", h.userIdentity, h.requireRoles("admin", "agent", "viewer"))
		{
			plan.POST("/", h.createPlan)
			plan.GET("/", h.getPlans)
			plan.GET("/main", h.getPlanStatisticsMain)
			plan.GET("/agent", h.getPlanStatisticsForAgent)
			plan.PUT("/", h.updatePlan)
		}

		furniture := api.Group("/furniture", h.userIdentity, h.requireRoles("admin", "merchandiser"))
		{
			furniture.GET("/mobile", h.GetTradePointFurnituresForMobile)
			furniture.GET("/", h.GetFurnituresPage)
			furniture.POST("/", h.CreateFurniture)
			furniture.PUT("/", h.UpdateFurniture)
			furniture.DELETE("/", h.DeleteFurniture)
			furniture.GET("/select", h.GetFurnitureListForSelect)
			furniture.GET("/point", h.GetFurnituresForTradePoint)
			furniture.GET("/point_create", h.CreateFurnitureForTradePoint)
			furniture.POST("/point_create", h.CreateFurnitureForTradePoint)
			furniture.PUT("/point_update", h.UpdateFurnitureForTradePoint)
			furniture.DELETE("/point", h.DeleteFurnitureForTradePoint)
		}

		merchandiser := api.Group("/merchandisers", h.userIdentity, h.requireRoles("admin", "merchandiser"))
		{
			merchandiser.POST("/", h.createMerchandiser)
			merchandiser.GET("/", h.getAllMerchandisers)
			merchandiser.GET("/:id/lock", h.statusMerchandiser)
			merchandiser.GET("/:id", h.getMerchandiserById)
			merchandiser.PUT("/:id", h.updateMerchandiser)
			merchandiser.DELETE("/:id", h.deleteMerchandiser)
			merchandiser.GET("/select", h.getMerchandiserForSelect)
			merchandiser.GET("/marketing", h.getTradePointsListForMerchandiser)
		}

		promo := api.Group("/promo", h.userIdentity, h.requireRoles("admin"))
		{
			promo.POST("/", h.createPromo)
			promo.GET("/page", h.getPromoPage)
			promo.PUT("/:id", h.updatePromoInformation)
			promo.GET("/:id/lock", h.changePromoStatus)
			promo.DELETE("/:id", h.deletePromo)
		}

	}
	return router
}
