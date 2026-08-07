package repository

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ExtraOrderPostgres struct {
	db *gorm.DB
}

func NewExtraOrderPostgres(db *gorm.DB) *ExtraOrderPostgres {
	return &ExtraOrderPostgres{db: db}
}

func (r *ExtraOrderPostgres) GetTradePoint(day_id, trade_agent_id int, language string) (*[]models.TradePointMobile, error) {
	var tradePoints []models.TradePoint
	var result []models.TradePointMobile

	get_trade_points_request := r.db.Model(&models.TradePoint{}).Preload("City").
		Where("trade_agent_id = ?", trade_agent_id).
		Where("day_id = ?", day_id).
		Where("status = ?", true).
		Find(&tradePoints)

	if get_trade_points_request.Error != nil {
		return nil, get_trade_points_request.Error
	}

	orderRepo := OrderPostgres{db: r.db}

	tradePoints = orderRepo.ADDExtraTradePoints(trade_agent_id, day_id, tradePoints)

	var routes []models.RouteOrder

	route_order_get_request := r.db.Model(&models.RouteOrder{}).Order("order_number").
		Where("trade_agent_id = ?", trade_agent_id).Where("day_id = ?", day_id).Find(&routes)

	if route_order_get_request.Error != nil {
		fmt.Println(route_order_get_request.Error)
	}

	for _, tradePoint := range tradePoints {
		var data models.TradePointMobile
		var contact contact_object

		contact = orderRepo.GetTredePointContact(tradePoint.Id, language)

		data.Id = tradePoint.Id
		data.Code = tradePoint.Code
		data.Contact = contact.Name
		data.Number = contact.Number
		data.Status = false

		if language == "ru" {
			data.City = tradePoint.City.NameRu
			data.Name = tradePoint.NameRu
			data.Location = tradePoint.LocationRu
			data.Orientir = tradePoint.OrientirRu
		} else {
			data.City = tradePoint.City.NameTm
			data.Name = tradePoint.NameTm
			data.Location = tradePoint.LocationTm
			data.Orientir = tradePoint.OrientirTm
		}

		result = append(result, data)
	}

	var new_result []models.TradePointMobile

	for _, route := range routes {
		for _, data := range result {
			if route.TradePointId == data.Id {
				new_result = append(new_result, data)
			}
		}
	}

	for _, data := range result {
		ishave := false
		for _, route := range routes {
			if route.TradePointId == data.Id {
				ishave = true
			}
		}

		if ishave == false {
			new_result = append(new_result, data)
		}
	}

	if new_result == nil {
		new_result = []models.TradePointMobile{}
	}

	return &new_result, nil
}

func (r *ExtraOrderPostgres) GetOrder(day, month, year, tradePointId int, language string) (models.GetOrderObject, error) {
	var result models.GetOrderObject
	var order models.Order
	var orderslist []models.OrderList
	//fmt.Println(day, ";", month, ";", year, "++", day)
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	var productsList []models.ProductOrder
	var price float64
	fmt.Println(crDate)
	orderRequest := r.db.Model(&models.Order{}).Where("trade_point_id = ? AND created_at = ?", tradePointId, crDate).Find(&order)

	if orderRequest.Error != nil {

	}
	fmt.Println(order)
	//emptyOrder := models.Order{}

	if order.Id > 0 {
		listRequest := r.db.Model(&models.OrderList{}).Where("order_id = ?", order.Id).Find(&orderslist)
		if listRequest.Error != nil {
			//
		}

		for _, element := range orderslist {
			var data models.ProductOrder

			var product models.Product

			productRequest := r.db.Model(&models.Product{}).Preload("Name").Preload("ProductType").Where("id = ?", element.ProductId).Find(&product)

			if productRequest == nil {

			}
			data.Id = element.Id
			data.ProductId = element.ProductId
			data.Count = element.Count
			data.Status = element.Status
			data.Price = element.Price
			price = price + element.SalePrice
			data.Code = product.Code
			data.Volume = product.ProductType.Volume
			data.Name = product.Name.Name

			if language == "ru" {
				data.Taste = product.TasteRu
				data.ProductTypeName = product.ProductType.NameRu
			} else {
				data.Taste = product.TasteTm
				data.ProductTypeName = product.ProductType.NameTm
			}
			productsList = append(productsList, data)
		}

	}
	if productsList == nil {
		productsList = []models.ProductOrder{}
		price = 0
	}

	var images []models.Image
	var readImage []models.ReadImage

	resultImage := r.db.Model(models.Image{}).Where("trade_point_id = ?", tradePointId).Find(&images)

	if resultImage.Error != nil {
		fmt.Println(resultImage.Error)
	}

	for _, image := range images {
		var item models.ReadImage

		item.Id = image.Id
		item.Path = "images/" + strconv.Itoa(image.TradePointId) + "/" + image.Name

		readImage = append(readImage, item)
	}

	if readImage == nil {
		readImage = []models.ReadImage{}
	}

	result.Products = productsList
	result.Price = price
	result.Images = readImage

	return result, nil
}
