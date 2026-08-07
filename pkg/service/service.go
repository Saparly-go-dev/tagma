package service

import (
	"mime/multipart"

	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type Authorization interface {
	CreateUser(user tagma.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, string, error)
	GetUserType(username, password string) (string, error)
	GetTradeAgentId(Id int) (int, error)
	GetEkspeditorId(Id int) (int, error)
	GetTradeAgentIdFromEkspeditorId(Id int) (int, error)
	GetEkspeditoryIdFromAgentId(Id int) (int, error)
}

type City interface {
	Create(city tagma.CreateCity) (int, error)
	GetAll(pageSize, pageNumber int, name string) (*tagma.CityPage, error)
	ChangeStatus(cityId int) error
	GetById(cityId int) (tagma.City, error)
	Delete(cityId int) error
	Update(cityId int, city tagma.CreateCity) error
}

type Trade_Channel interface {
	CreateType(tip models.ChannelObject) error
	GetAllTypes() ([]models.ChannelType, error)
	DeleteTypes(Id int) error
	UpdateTypes(Id int, tip models.ChannelObject) error
	CreateStructure(structure models.ChannelObject) error
	GetAllStructure() ([]models.ChannelStructure, error)
	DeleteStructure(Id int) error
	UpdateStructure(Id int, structure models.ChannelObject) error
	CreateSize(size models.ChannelObject) error
	GetAllSizes() ([]models.ChannelSize, error)
	DeleteSize(Id int) error
	UpdateSize(Id int, size models.ChannelObject) error
	CreateManagement(management models.ChannelObject) error
	GetAllManagements() ([]models.ChannelManagement, error)
	DeleteManagement(Id int) error
	UpdateManagement(Id int, management models.ChannelObject) error
	Create(channel models.CreateChannel) error
	GetChannelPage(pageSize, pageNumber, typeId int, language string) (*models.TradeChannelPage, error)
	ChangeStatus(Id int) error
	Delete(Id int) error
	Update(Id int, channel models.CreateChannel) error
	GetTypeForChannel(language string) (*[]models.SelectObject, error)
	GetStructureForChannel(language string) (*[]models.SelectObject, error)
	GetSizeForChannel(language string) (*[]models.SelectObject, error)
	GetManagementForChannel(language string) (*[]models.SelectObject, error)
}

type Trade_Category interface {
	Create(category tagma.CreateTrade_category) (int, error)
	GetAll(pageNumber, pageSize int, name string) (*tagma.TradeCategoryPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (tagma.Trade_category, error)
	Delete(Id int) error
	Update(Id int, category tagma.CreateTrade_category) error
}

type ProductType interface {
	Create(category models.CreateProductType) (int, error)
	GetAll(pageNumber, pageSize int, name string) (*models.ProductTypePage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.ProductType, error)
	Delete(Id int) error
	Update(Id int, productType models.CreateProductType) error
	GetAllVolumes() (*[]models.Volume, error)
}

type Product interface {
	Create(category models.CreateProduct) (int, error)
	GetAll(pageNumber, pageSize int, name, language string) (*models.ProductPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.CreateProduct, error)
	Delete(Id int) error
	Update(Id int, productType models.CreateProduct) error
	GetTypes(language string) ([]models.SelectObject, error)
	GetProductBrands() ([]models.SelectObject, error)
	GetInformantionAboutSales(brandId, productTypeId, agentId, cityId, districtId, tradePointId, channelTypeId, startDay, startMonth, startYear, endDay, endMonth, endYear int, volume float64, product_taste, language string) ([]models.ReadProductForTradePoint, error)
	GetDailySalesProductsInExcel(agentId, day, month, year int, language string) ([]byte, error)
}

type Driver interface {
	Create(Id models.CreateDriver) (int, error)
	GetAll(pageNumber, pageSize int, name, language string) (*models.DriverPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.Driver, error)
	Delete(Id int) error
	Update(Id int, driver models.CreateDriver) error
}

type Ekspeditor interface {
	Create(Id models.CreateEkspeditor) (int, error)
	GetAll(pageNumber, PageSize int, name, language string) (*models.EkspeditorPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.Ekspeditor, error)
	Delete(Id int) error
	Update(Id int, driver models.CreateEkspeditor) error
	GetProductCount(day, month, year, ekspeditorId int) (*map[int]int, error)
	GetPointProductList(day, month, year, tradePointId int, language string) (models.GetOrderObject, error)
	SaveProductList(day, month, year, ekspeditorId int, orders []models.CreateOrder) error
	GetTradePoints(day, month, year, ekspeditorId int, language string) (*[]models.TradePointMobile, error)
	EkspeditorProducts(day, month, year, ekspeditorId int) (*[]models.EkspData, error)
	OrderStatus(day, month, year, tradePointId, ekspeditorId int) error
	GetGalyndy(ekspeditorId int) (*[]models.CreateOrder, error)
	GetEkspeditorGalyndy(ekspeditorId int, language string) ([]models.ReadProductForTradePoint, error)
	DeleEkspData() error
	GetProductCountWithMinus(day, month, year, ekspeditorId int) (*map[int]int, error)
	GetProductInTheCarWithDifference(ekspeditorId int, language string) ([]models.ReadProductForTradePoint, error)
	GetProductInTheCarWithDifferenceInArray(language string) ([]models.GetGalyndyInArray, error)
	GetProductPriceInTheCarWithDifference(ekspeditorId int) (float64, error)
}

type TradeAgent interface {
	Create(agent models.CreateTradeAgent) (int, error)
	GetAll(pageSize, pageNumber int, name, language string) (*models.TradeAgentPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.TradeAgent, error)
	Delete(Id int) error
	Update(Id int, agent models.CreateTradeAgent) error
	GetDrivers(name, language string) (*[]models.SelectObject, error)
	GetEkspeditors(name, language string) (*[]models.SelectObject, error)
	AgentExcel(name, language string) ([]byte, error)
	GetAgentDebt(AgentId int, language string) (models.AgentDebt, error)
	GetMerchendisingInformation(day, month, year, agentid int, language string) (models.PointsMerchendising_List, error)
	CreateCSVFileForFactory(day, month, year, agentid int, language string) ([]byte, error)
	GetTradeAgentRevenue(month, year, agent_id int, language string) (float64, error)
}

type TradePoint interface {
	Create(point models.CreateTradePoint, image1, image2 *multipart.FileHeader) (int, error)
	GetExcel(cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId int, name, language string) ([]byte, error)
	GetAll(pageSize, pageNumber, cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, status int, name, code, contact, language string) (*models.TradePointPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.TradePoint, error)
	Delete(Id int) error
	Update(Id int, point models.CreateTradePoint) error
	ChangeTradePointAgent(agent_id, ekspeditor_id, day_id int, ids []int) error
	GetCities(name, language string) (*[]models.SelectObject, error)
	GetDistrictsByCity(cityId int, language string) (*[]models.SelectObject, error)
	GetAgents(name, language string) (*[]models.SelectObject, error)
	GetChannels(name, language string) (*[]models.SelectObject, error)
	GetGeneralInformation(Id int, language string) (*models.GeneralTradePoint, error)
	DeleteImage(Id int) error
	ReSaveImage(PointId int, file *multipart.FileHeader) error
	GetPointImages(Id int) (*[]models.ReadImage, error)
	GetInformationAboutOrder(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) ([]models.ReadOrder, error)
	GetInformantionAboutProducts(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) ([]models.ReadProductForTradePoint, error)
	GetTradePointSaleHistoryForMonth(tradePointId int, language string) ([]models.SaleHistoryTradePoint, error)
}

type Contact interface {
	Create(contact models.CreateContact) (int, error)
	GetAll(PointId int, language string) (*[]models.ReadContact, error)
	GetById(Id int) (*models.Contact, error)
	Delete(Id int) error
	Update(Id int, contact models.CreateContact) error
	GetKinds(name, language string) (*[]models.SelectObject, error)
	GetPosts(name, language string) (*[]models.SelectObject, error)
	GetUprs(name, language string) (*[]models.SelectObject, error)
}

type Route interface {
	Create(route models.CreateRoute) (int, error)
	GetById(Id int, language string) (*models.ReadRoute, error)
	GetAll(pageSize, pageNumber, agentId, cityId, dayId int, language string) (*models.RoutePage, error)
	ChangeStatus(Id int) error
	Delete(Id int) error
	Update(Id int, route models.CreateRoute) error
	GetDays(name, language string) (*[]models.SelectObject, error)
	GetPoints(cityId, districtId int, name, language string) (*[]models.PointSelect, error)
	GetEkspeditors(TradeAgentId int, language string) (*[]models.SelectObject, error)
	GetRouteOrder(agent_id, day_id int, language string) ([]models.SelectObject, error)
	PostRouteOrder(agent_id, day_id int, list []models.SelectObject) error
	ExcelRouteOrder(agent_id, day_id int, language string) ([]byte, error)
}

type Log interface {
	GetAll(language string, pageSize, pageNumber int, logType string) (*models.LogPage, error)
	GetLogTypes(language string) (*[]string, error)
}

type Order interface {
	GetTradePoints(day, month, year, agentid int, language string) (*[]models.TradePointMobile, error)
	GetTradePointsForMarketing(day, month, year, agentid int, language string) (*[]models.TradePointMobile, error)
	GetProducts(language string) (*[]models.ProductDistinct, error)
	GetBrand(language string) (*[]string, error)
	GetModel(brandId, ProductTypeId int, language string) (*[]string, error)
	GetTypes(language string) (*[]string, error)
	GetVolume(brand, model, tip string) (*[]models.SelectObject, error)
	GetSum(Id, Count int) (float64, error)
	SaveOrder(day, month, year, tradePointId, paymentTypeId int, first, is_credit bool, orders []models.CreateOrder) error
	Marketing(day, month, year, tradePointId int, list []models.MarketingPost) error
	GetOrderPage(pageSize, pageNumber, tradeAgentId, cityId, districtId, day, month, year int, name, language, location, code, order string, status, is_status_active bool) (*models.OrderPage, error)
	OrderExcel(tradeAgentId, cityId, districtId int, name, language string) ([]byte, error)
	GetOrderList(OrderId int, language string) (*[]models.ReadOrderList, error)
	GetOrders(day, month, year, tradePointId int, language string) (models.GetOrderObject, error)
	CloseTheOrder(orderId int) error
}

type Discount interface {
	Create(discount models.CreateDiscount) error
	GetAll(pageSize, pageNumber int, name, language string) (*models.DiscountPage, error)
	ChangeStatus(Id int) error
	Delete(Id int) error
	Update(Id int, discount models.CreateDiscount) error
	UpdateNameAndDate(id int, data models.UpdateDiscountInfo) error
	GetTradePointByTradeAgent(tradeAgentId int) ([]models.SelectObject, error)
}
type Gift interface {
	Create(gift models.CreateGift) error
	GetPage(pageSize, pageNumber int, name, language string) (*models.GiftPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (*models.Gift, error)
	Delete(Id int) error
	Update(Id int, gift models.CreateGift) error
	GetGiftTypes(language string) (*[]models.SelectObject, error)
}

type District interface {
	Create(district models.CreateDistrict) error
	GetPage(pageSize, pageNumber int, name, language string) (*models.DistrictPage, error)
	GetById(Id int) (*models.District, error)
	ChangeStatus(Id int) error
	Delete(Id int) error
	Update(Id int, district models.CreateDistrict) error
}

type Note interface {
	Create(note models.CreateNote) error
	GetAll(pageSize, pageNumber int) (*models.NotePage, error)
	Update(Id int, note models.CreateNote) error
	Delete(Id int) error
	CreatePointNote(note models.CreatePointNote) error
	GetAllPointNote(pageSize, pageNumber, tradePointId int, code, language string) (*models.PointNotePage, error)
	UpdatePointNote(Id int, pointNote models.CreatePointNote) error
	DeletePointNote(Id int) error
}

type User interface {
	Create(user tagma.CreateUser) error
	GetPage(pageSize, pageNumber int, name, tip, language string) (tagma.UserPage, error)
	ChangePassword(data tagma.Change_Password) error
	ChangeStatus(Id int) error
	Delete(Id int) error
	Update(Id int, data tagma.CreateUser) error
}

type Manual interface {
	Create(data models.CreateManual) error
	Get(language string) (*models.ReadManual, error)
	GetForMobile(language string) (*models.ReadManualForMobile, error)
	UpdateManual(Id int, data models.CreateManual) error
	DeleteImage(Id int) error
	SaveImage(Id int, file *multipart.FileHeader) error
}

type Payment interface {
	GetTradePointsForPayment(day, month, year, agentId int, language string) (*[]models.TradePointForPayment, error)
	SavePayment(agentId int, data models.CreatePayment) error
	GetPaymentTypes(language string) ([]models.SelectObject, error)
	GetPaymentForOrder(orderId int, language string) ([]models.ReadPayment, error)
	GetAgentPayments(agent_id, payment_type_id int, status bool, point_code, language string) ([]models.ReadPayment, error)
	ConfirmPayments(agent_id, price1, price2 int, ids []int) error
	UpdatePayment(id, price1, price2 int) error
	UpdateDebt(agent_id, price1, price2 int) error
	AddSelfDebt(agent_id int, currency float64) error
	GetTradePointAllDebtSum(tradePointId int) float64
	PostPayment(currency float64, agent_id, tradepointId, PaymentTypeId int, is_agent bool) error
	GetTradePointPaymentsHistory(agent_id, tradePointId int, language string) ([]models.ReadPayment, error)
	GetDailyPaymentsInformations(agent_id, day, month, year int, language string) (models.ObjectWithIdNameAndCurrency, error)
	GetTradePointsDebtList(agentid int, language string) ([]models.TradePointMobile, error)
}

type Marketing interface {
	GetMarketingMainPage(cityId, districtId, tradePointId, agentId, brand int, volume float64, taste, language string) ([]models.MarketingMainPage, error)
	MarketingMainPageExport(cityId, districtId, tradePointId, agentId, brand int, volume float64, taste, language string) ([]byte, error)
	GetDistribution(tradePointId int, language string) ([]models.ReadMarketingById, error)
	DistributionExport(tradePointId int, language string) ([]byte, error)
	GetRasprodano(tradePointId int, language string) ([]models.ReadMarketingById, error)
	RasprodanoExport(tradePointId int, language string) ([]byte, error)
	GetWystawlaymost(tradePointId int, language string) ([]models.ReadMarketingById, error)
	WystawlaymostExport(tradePointId int, language string) ([]byte, error)
}

type Dealer interface {
	Create(dealer models.CreateDealer) error
	GetPage(pageSize, pageNumber, cityId int, name, language string) (models.DealerPage, error)
	Update(Id int, dealer models.CreateDealer) error
	ChangeStatus(Id int) error
	GetDealerForSelect(name, language string) ([]models.SelectObject, error)
	TopUpDealerAccount(dealerId, price1, price2 int) error
}

type DealerProduct interface {
	Create(product models.CreateDealerProduct) (int, error)
	GetAll(pageSize, pageNumber int, name, language string) (*models.DealerProductPage, error)
	Delete(productId int) error
	Update(DealerProductId int, DealerProduct models.CreateDealerProduct) error
	ChangeStatus(Id int) error
}

type DealerOrder interface {
	Save(DealerId int, list []models.CreateOrder) error
	Update(DealerOrderId int, list []models.CreateOrder) error
	CompleTheDealerOrder(dealerOrderId int) error
	ChangeStatus(Id int) error
	GetOrders(pageSize, pageNumber, DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) (models.DealerOrderPage, error)
	GetProduct(DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear int, status bool, language string) ([]models.ReadProductForTradePoint, error)
	GetOrderList(Id int, language string) (*[]models.ReadOrderList, error)
	DeleteDealerOrder(orderId int) error
}

type Promotion interface {
	Create(promotion models.CreatePromotion) error
	ChangeStatus(Id int) error
	GetAll(pageSize, pageNumber int, name, language string) (*models.PromotionPage, error)
	Update(id int, data models.CreatePromotion) error
	UpdateNameAndTime(id int, data models.UpdateDiscountInfo) error
	Delete(id int) error
}

type ExtraTradePointDay interface {
	Update(trade_point_id int, list []int) error
	Get(trade_point_id int) []int
}

type Test interface {
	Create(test tagma.Test) (int, error)
	GetAll() ([]tagma.Test, error)
}

type ExtraOrder interface {
	GetTradePoint(day_id, trade_agent_id int, language string) (*[]models.TradePointMobile, error)
	GetOrder(day, month, year, tradePointId int, language string) (models.GetOrderObject, error)
}

type Plan interface {
	Create(list []models.CreatePlan, month, year int) error
	GetPlans(month, year int, language string) ([]models.ReadPlan, error)
	GetPlanStatisticsForMain(month, year, agent_id int, language string) ([]models.PlanStatisticsMain, error)
	GetPlansForAgent(month, year, agent_id int, language string) ([]models.PlanStatisticsForAgent, error)
	Update(month, year int, list []models.UpdatedPlan) error
}

type Furniture interface {
	GetMobile(trade_point_id int, language string) ([]models.ReadFurnitureForMobile, error)
	Get(name string, pageSize, pageNUmber int) (models.FurniturePage, error)
	Create(data models.CreateFurniture) error
	Update(id int, data models.CreateFurniture) error
	Delete(id int) error
	GetFurnitureList(name, language string) ([]models.SelectObject, error)
	GetTradePointFurnitures(trade_point_id int, language string) ([]models.TradePointFurnitureRead, error)
	CreateTradePointFurniture(trade_point_id, furniture_id, count int) error
	UpdateTradePointFurniture(id, furniture_id, count int) error
	DeleteTradePointFurniture(id int) error
}

type Merchandiser interface {
	Create(merchandiser models.CreateMerchandiser) (int, error)
	GetAll(pageNumber, pageSize int, name, language string) (*models.MerchandiserPage, error)
	ChangeStatus(Id int) error
	GetById(Id int) (models.Merchandiser, error)
	Delete(Id int) error
	Update(Id int, merchandiser models.CreateMerchandiser) error
	GetMerchandiserForSelect(language, name string) ([]models.SelectObject, error)
	GetMarketingPageForMobile(day, month, year, userId int, language string) (*[]models.TradePointMobile, error)
}

type TradePointDebt interface {
	Create(tradePointId int, currency float64) error
	GetDebtForMobile(trade_point_id int, language string) ([]models.TradePointDebtResponseForMobile, error)
}

type Promo interface {
	Create(promo models.UpdatePromo) error
	Update(promo_id int, promo models.UpdatePromo) error
	GetPage(pageSize, pageNumber int, name, language string) (models.PromoPage, error)
	ChangeStatus(id int) error
	Delete(id int) error
}

type Service struct {
	Authorization
	City
	Trade_Channel
	Trade_Category
	Test
	ProductType
	Product
	Driver
	Ekspeditor
	TradeAgent
	TradePoint
	Contact
	Route
	Log
	Order
	Discount
	Gift
	District
	Note
	User
	Payment
	Manual
	Marketing
	Dealer
	DealerProduct
	DealerOrder
	Promotion
	ExtraTradePointDay
	ExtraOrder
	Plan
	Furniture
	Merchandiser
	TradePointDebt
	Promo
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization:      NewAuthService(repos.Authorization),
		City:               NewCityService(repos.City),
		Trade_Channel:      NewTradeChannelService(repos.Trade_Channel),
		Trade_Category:     NewTradeCategoryService(repos.Trade_Category),
		Test:               NewTestService(repos.Test),
		ProductType:        NewProductTypeService(repos.ProductType),
		Product:            NewProductService(repos.Product),
		Driver:             NewDriverService(repos.Driver),
		Ekspeditor:         NewEkspeditorService(repos.Ekspeditor),
		TradeAgent:         NewTradeAgentService(repos.TradeAgent),
		TradePoint:         NewTradePointService(repos.TradePoint),
		Contact:            NewContactService(repos.Contact),
		Route:              NewRouteService(repos.Route),
		Log:                NewLogService(repos.Log),
		Order:              NewOrderService(repos.Order),
		Discount:           NewDiscountService(repos.Discount),
		Gift:               NewGiftService(repos.Gift),
		District:           NewDistrictService(repos.District),
		Note:               NewNoteService(repos.Note),
		User:               NewUserService(repos.User),
		Payment:            NewPaymentService(repos.Payment),
		Manual:             NewManualService(repos.Manual),
		Marketing:          NewMarketingService(repos.Marketing),
		Dealer:             NewDealerService(repos.Dealer),
		DealerProduct:      NewDealerProductService(repos.DealerProduct),
		DealerOrder:        NewDealerOrderService(repos.DealerOrder),
		Promotion:          NewPromotionService(repos.Promotion),
		ExtraTradePointDay: NewExtraTradePointDayService(repos.ExtraTradePointDay),
		ExtraOrder:         NewExtraOrderService(repos.ExtraOrder),
		Plan:               NewPlanService(repos.Plan),
		Furniture:          NewFurnitureService(repos.Furniture),
		Merchandiser:       NewMerchandiserService(repos.Merchandiser),
		TradePointDebt:     NewTradePointDebtService(repos.TradePointDebt),
		Promo:              NewPromoService(repos.Promo),
	}
}
