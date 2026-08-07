package repository

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type EkspeditorPostgres struct {
	db *gorm.DB
}

func NewEkspeditorPostgres(db *gorm.DB) *EkspeditorPostgres {
	return &EkspeditorPostgres{db: db}
}

func (r *EkspeditorPostgres) Create(ekspeditor models.CreateEkspeditor) (int, error) {
	newEkspeditor := models.Ekspeditor{
		NameRu:    ekspeditor.NameRu,
		NameTm:    ekspeditor.NameTm,
		Status:    ekspeditor.Status,
		Number:    ekspeditor.Number,
		UpdatedAt: time.Now(),
	}

	result := r.db.Create(&newEkspeditor)
	if err := result.Error; err != nil {
		return 0, err
	}

	return newEkspeditor.Id, nil
}

func (r *EkspeditorPostgres) GetAll(pageNumber, pageSize int, name, language string) (*models.EkspeditorPage, error) {
	var ekspeditors []models.Ekspeditor
	var pageData []models.ReadEkspeditor

	offset := (pageNumber - 1) * pageSize

	result := r.db.Model(&models.Ekspeditor{}).Preload("TradeAgent").Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Order("id").Offset(offset).Limit(pageSize).Find(&ekspeditors)

	if result.Error != nil {
		return nil, result.Error
	}

	var count int64
	result2 := r.db.Model(&models.Ekspeditor{}).Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)

	if result2.Error != nil {
		return nil, result.Error
	}

	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, ekspeditor := range ekspeditors {
		var data models.ReadEkspeditor
		var tradeAgent models.TradeAgent

		//fmt.Println(ekspeditor)
		if len(ekspeditor.TradeAgent) > 0 {
			tradeAgent = ekspeditor.TradeAgent[0]
		}

		//request := r.db.Model(&models.TradeAgent{}).Where("ekspeditor_id = ?", ekspeditor.Id).First(&tradeAgent)
		//if request.Error != nil {
		//	fmt.Println(request.Error)
		//}

		data.Id = ekspeditor.Id
		data.Code = tradeAgent.Code
		data.NameRu = ekspeditor.NameRu
		data.NameTm = ekspeditor.NameTm
		data.Number = ekspeditor.Number
		data.Status = ekspeditor.Status
		data.UpdatedAt = ekspeditor.UpdatedAt.Format("02.01.2006")
		if language == "ru" {
			data.TradeAgent = tradeAgent.NameRu
		} else {
			data.TradeAgent = tradeAgent.NameTm
		}

		pageData = append(pageData, data)
	}

	var page = models.EkspeditorPage{
		Items:      pageData,
		PageCount:  int(pageCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	//Creating page with calculate variables

	return &page, nil
}

func (r *EkspeditorPostgres) ChangeStatus(Id int) error {
	var ekspeditor models.Ekspeditor
	result := r.db.First(&ekspeditor, Id)

	if err := result.Error; err != nil {
		return err
	}

	// Toggle status
	newStatus := !ekspeditor.Status
	updateData := map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	// Use "=" in the Where clause instead of "=="
	result = r.db.Model(&models.Ekspeditor{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *EkspeditorPostgres) GetById(Id int) (models.Ekspeditor, error) {
	var ekspeditor models.Ekspeditor
	result := r.db.First(&ekspeditor, Id)
	return ekspeditor, result.Error
}

func (r *EkspeditorPostgres) Delete(Id int) error {
	result := r.db.Delete(&models.Ekspeditor{}, Id)
	return result.Error
}

func (r *EkspeditorPostgres) Update(Id int, ekspeditor models.CreateEkspeditor) error {
	updateData := map[string]interface{}{
		"name_ru": ekspeditor.NameRu,
		"name_tm": ekspeditor.NameTm,
		"status":  ekspeditor.Status,
		"number":  ekspeditor.Number,
	}

	result := r.db.Model(&models.Ekspeditor{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *EkspeditorPostgres) GetProductCount(day, month, year, ekspeditorId int) (*map[int]int, error) {
	var tradePoints []models.TradePoint
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	var result = make(map[int]int)

	dayOfWeek := int(crDate.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	request := r.db.Model(&models.TradePoint{}).Where("ekspeditor_id = ?", ekspeditorId).Where("day_id = ?", dayOfWeek).
		Find(&tradePoints)

	tradePoints = r.ADDExtraTradePointsForEkspeditor(ekspeditorId, dayOfWeek, tradePoints)
	tradePoints = r.ADDExtraOrderPointsForEkpeditor(day, month, year, ekspeditorId, tradePoints)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, item := range tradePoints {
		var order models.Order

		OrderRequest := r.db.Model(&models.Order{}).Preload("List").
			Where(`"trade_point_id" = ?`, item.Id).
			Where(`"created_at" = ?`, crDate).
			Find(&order)

		if OrderRequest.Error != nil {
			fmt.Println("Order tapylmady", OrderRequest.Error)
		}

		for _, element := range order.List {
			_, exists := result[element.ProductId]
			if exists {
				result[element.ProductId] += element.Count
			} else {
				result[element.ProductId] = element.Count
			}

		}

	}

	return &result, nil
}

func (r *EkspeditorPostgres) GetTradePoints(day, month, year, ekspeditorId int, language string) (*[]models.TradePointMobile, error) {
	var tradePoints []models.TradePoint
	var result []models.TradePointMobile
	var agent_id int

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	request := r.db.Model(models.TradePoint{}).Preload("City").
		Preload("District").
		Where("ekspeditor_id = ?", ekspeditorId).
		Where("day_id = ?", dayOfWeek).
		Where("status = ?", true).
		Find(&tradePoints)

	tradePoints = r.ADDExtraTradePointsForEkspeditor(ekspeditorId, dayOfWeek, tradePoints)
	tradePoints = r.ADDExtraOrderPointsForEkpeditor(day, month, year, ekspeditorId, tradePoints)

	if request.Error != nil {
		return nil, request.Error
	}
	orderRepo := OrderPostgres{db: r.db}

	for _, tradePoint := range tradePoints {
		agent_id = tradePoint.TradeAgentId

		var data models.TradePointMobile

		var order models.Order

		request3 := r.db.Model(&models.Order{}).
			Where(`"trade_point_id" = ?`, tradePoint.Id).
			Where(`"created_at" = ?`, date).
			Find(&order)

		if request3.Error != nil {

			fmt.Println("Order tapylmady", request3.Error)
		}
		// Если на указанную дату нет заказа, исключаем торговую точку из результата
		if order.Id == 0 {
			continue
		}

		data.Status = order.Status

		contact := orderRepo.GetTredePointContact(tradePoint.Id, language)

		data.Id = tradePoint.Id
		data.Code = tradePoint.Code
		data.Contact = contact.Name
		data.Number = contact.Number
		if language == "ru" {
			data.City = tradePoint.City.NameRu + ":" + tradePoint.District.NameRu
			data.Name = tradePoint.NameRu
			data.Location = tradePoint.LocationRu
			data.Orientir = tradePoint.OrientirRu
		} else {
			data.City = tradePoint.City.NameTm + ":" + tradePoint.District.NameTm
			data.Name = tradePoint.NameTm
			data.Location = tradePoint.LocationTm
			data.Orientir = tradePoint.OrientirTm
		}

		result = append(result, data)
	}

	var routes []models.RouteOrder

	route_order_get_request := r.db.Model(&models.RouteOrder{}).Order("order_number").
		Where("trade_agent_id = ?", agent_id).Where("day_id = ?", dayOfWeek).Find(&routes)

	if route_order_get_request.Error != nil {
		fmt.Println(route_order_get_request.Error)
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

func (r *EkspeditorPostgres) GetPointProductList(day, month, year, tradePointId int, language string) (models.GetOrderObject, error) {
	var result models.GetOrderObject
	var order models.Order
	var orderslist []models.OrderList
	var promo_sum float64
	//fmt.Println(day, ";", month, ";", year, "++", day)
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	var productsList []models.ProductOrder
	var price, sale_price float64
	//fmt.Println(crDate)
	orderRequest := r.db.Model(&models.Order{}).Preload("PaymentType").Where("trade_point_id = ? AND created_at = ?", tradePointId, crDate).Find(&order)
	//orderRequest := r.db.Model(&models.Order{}).Where("trade_point_id = ? ", tradePointId).Find(&order)

	if orderRequest.Error != nil {
		//order = models.Order{}
		//return models.GetOrderObject{}, orderRequest.Error
	}
	//fmt.Println(order)
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
			data.SalePrice = element.SalePrice

			price = price + element.Price
			sale_price = sale_price + element.SalePrice

			//Eger discountId ==1 bolsa onda product öz bahasynda ýa-da promo boýunça satylan
			if element.DiscountId == 1 {
				promo_sum += element.Price - element.SalePrice
			}
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

	if language == "ru" {
		result.PaymentType = order.PaymentType.NameRu
	} else {
		result.PaymentType = order.PaymentType.NameTm
	}

	result.Products = productsList
	result.Price = price
	result.SalePrice = sale_price
	result.Images = readImage
	result.Promo = promo_sum
	result.Discount = result.Price - result.SalePrice - result.Promo
	result.Is_Credit = order.IsCredit

	return result, nil
}

func (r *EkspeditorPostgres) SaveProductList(day, month, year, ekspeditorId int, orders []models.CreateOrder) error {
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	//list_checked := true
	//
	//for _, item := range orders {
	//	var product models.Product
	//
	//	get_product_information := r.db.Model(&models.Product{}).Preload("ProductType").Where("id = ?", item.Id).Find(&product)
	//
	//	if get_product_information.Error != nil {
	//		list_checked = false
	//	}
	//
	//	if item.Count%product.ProductType.Count != 0 {
	//		list_checked = false
	//	}
	//}
	//
	//if list_checked == false {
	//	return errors.New("Rugsat berilmeýär")
	//}

	for _, data := range orders {
		var eksData models.EkspData

		ishaveRequest := r.db.Model(&models.EkspData{}).Where("ekspeditor_id = ? AND product_id = ?", ekspeditorId, data.Id).Find(&eksData)

		if ishaveRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy", eksData)
		}
		if eksData.Id > 0 && eksData.CreatedAt == crDate {
			return errors.New("Rugsat berilmeýär")
		}

		if eksData.Id == 0 {
			newdata := models.EkspData{
				EkspeditorId: ekspeditorId,
				ProductId:    data.Id,
				Count:        data.Count,
				CreatedAt:    time.Now(),
			}

			request := r.db.Model(&models.EkspData{}).Create(&newdata)
			if request.Error != nil {
				return request.Error
			}
		} else {
			updateData := map[string]interface{}{
				"count":      eksData.Count + data.Count,
				"created_at": time.Now(),
			}
			//fmt.Println(updateData)

			updateRequest := r.db.Model(&models.EkspData{}).Where("id = ?", eksData.Id).Updates(updateData)
			if updateRequest.Error != nil {
				fmt.Println("maglumat uytgedilende yalnyshlyk yuze cykdy", eksData)
			}
		}
	}

	return nil
}

func (r *EkspeditorPostgres) OrderStatus(day, month, year, tradePointId, ekspeditorId int) error {
	// ekspeditoryň söwda nokadyndaky goýup gaýdan produktlaryny sürüjiniň ulagyndaky produktlaryndan aýyrmaly
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	var order models.Order
	getRequest := r.db.Model(&models.Order{}).Preload("List").Where("trade_point_id = ? AND created_at = ?", tradePointId, crDate).Find(&order)

	if getRequest.Error != nil || order.Id == 0 {
		return getRequest.Error
	}

	is_ekspeditor_have_enough_products := true

	for _, item := range order.List {
		var eksData models.EkspData

		ekspRequest := r.db.Model(&models.EkspData{}).Where("ekspeditor_id = ? AND product_id = ?", ekspeditorId, item.ProductId).Find(&eksData)

		if ekspRequest.Error != nil || eksData.Id == 0 {
			is_ekspeditor_have_enough_products = false
		}

		if (eksData.Count - item.Count) < 0 {
			is_ekspeditor_have_enough_products = false
		}
	}

	if is_ekspeditor_have_enough_products == false {
		return errors.New("Rugsat berilmeýär")
	}

	for _, item := range order.List {
		var eksData models.EkspData

		ekspRequest := r.db.Model(&models.EkspData{}).Where("ekspeditor_id = ? AND product_id = ?", ekspeditorId, item.ProductId).Find(&eksData)

		if ekspRequest.Error != nil {
			fmt.Println(ekspRequest.Error)
		}

		updateData := map[string]interface{}{
			"count": eksData.Count - item.Count,
		}

		updateRequest := r.db.Model(&models.EkspData{}).Where("id = ?", eksData.Id).Updates(updateData)
		if updateRequest.Error != nil {
			fmt.Println(updateRequest.Error)
		}
	}
	// Berilen zakazy tapyp onuň statusyny true etmeli

	//fmt.Println(day, ".", month, ".", year)

	updateData := map[string]interface{}{
		"status": true,
	}
	updateRequest := r.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(updateData)
	return updateRequest.Error
}

func (r *EkspeditorPostgres) EkspeditorProducts(day, month, year, ekspeditorId int) (*[]models.EkspData, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	var data []models.EkspData

	request := r.db.Model(&models.EkspData{}).Where("ekspeditor_id = ? AND DATE(created_at) = ?", ekspeditorId, date).Find(&data)

	if request.Error != nil {
		//return []models.EkspData{}, nil
		data = []models.EkspData{}
	}

	if data == nil {
		data = []models.EkspData{}
	}

	return &data, nil
}

func (r *EkspeditorPostgres) GetGalyndy(ekspeditorId int) (*[]models.CreateOrder, error) {
	var list []models.EkspData
	var result []models.CreateOrder
	today := time.Now()
	day := today.Day()
	month := today.Month()
	year := today.Year()

	if result == nil {
		result = []models.CreateOrder{}
	}

	request := r.db.Model(&models.EkspData{}).Where("ekspeditor_id = ?", ekspeditorId).Order("product_id").Find(&list)
	if request.Error != nil {
		fmt.Println(request.Error)
	}

	var metka int

	for _, list_item := range list {
		cryear, crmonth, crday := list_item.CreatedAt.Date()

		if year == cryear && month == crmonth && day == crday {
			metka++
		}
	}

	//
	//if metka == 0 {
	//	return &result, nil
	//}
	//
	//for _, item := range list {
	//	var data models.CreateOrder
	//	data.Id = item.ProductId
	//	data.Count = item.Count
	//
	//	result = append(result, data)
	//}

	return &result, nil
}

func (r *EkspeditorPostgres) GetEkspeditorGalyndy(ekspeditorId int, language string) ([]models.ReadProductForTradePoint, error) {
	var products []models.Product
	var result []models.ReadProductForTradePoint

	getRequest := r.db.Model(&models.Product{}).Preload("Name").Preload("ProductType").Where("id > 0").Find(&products)
	if getRequest.Error != nil {
		fmt.Println(getRequest.Error)
	}

	for _, item := range products {
		var data models.EkspData
		request := r.db.Model(&models.EkspData{}).Where("product_id = ?", item.Id).Where("ekspeditor_id = ?", ekspeditorId).Find(&data)
		if request.Error != nil {
			fmt.Println(request.Error)
		}

		var result_item models.ReadProductForTradePoint
		result_item.Id = item.Id
		result_item.Volume = item.ProductType.Volume
		result_item.Count = data.Count
		result_item.Name = item.Name.Name

		if language == "ru" {
			result_item.Taste = item.TasteRu
			result_item.ProductType = item.ProductType.NameRu
		} else {
			result_item.Taste = item.TasteTm
			result_item.ProductType = item.ProductType.NameTm
		}

		result = append(result, result_item)
	}

	if result == nil {
		result = []models.ReadProductForTradePoint{}
	}

	return result, nil
}

func (r *EkspeditorPostgres) DeleEkspData() error {
	delete_ekspdata_request := r.db.Model(&models.EkspData{}).Where("id > 0").Delete(&models.EkspData{})

	if delete_ekspdata_request.Error != nil {
		return delete_ekspdata_request.Error
	}

	return nil
}

func (r *EkspeditorPostgres) GetProductCountWithMinus(day, month, year, ekspeditorId int) (*map[int]int, error) {
	var tradePoints []models.TradePoint
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	var result = make(map[int]int)

	dayOfWeek := int(crDate.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	request := r.db.Model(&models.TradePoint{}).Where("ekspeditor_id = ?", ekspeditorId).Where("day_id = ?", dayOfWeek).
		Find(&tradePoints)

	tradePoints = r.ADDExtraTradePointsForEkspeditor(ekspeditorId, dayOfWeek, tradePoints)
	tradePoints = r.ADDExtraOrderPointsForEkpeditor(day, month, year, ekspeditorId, tradePoints)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, item := range tradePoints {
		var order models.Order

		OrderRequest := r.db.Model(&models.Order{}).Preload("List").
			Where(`"trade_point_id" = ?`, item.Id).
			Where(`"created_at" = ?`, crDate).
			Find(&order)

		if OrderRequest.Error != nil {
			fmt.Println("Order tapylmady", OrderRequest.Error)
		}

		for _, element := range order.List {
			_, exists := result[element.ProductId]
			if exists {
				result[element.ProductId] += element.Count
			} else {
				result[element.ProductId] = element.Count
			}

		}

	}

	var eksp_data []models.EkspData

	get_eksp_data_request := r.db.Model(&models.EkspData{}).Where("ekspeditor_id = ?", ekspeditorId).Find(&eksp_data)
	if get_eksp_data_request.Error != nil {
		return &result, nil
	}

	for idx, _ := range result {
		for _, eksp_data_item := range eksp_data {
			if idx == eksp_data_item.ProductId {
				result[idx] -= eksp_data_item.Count
				if result[idx] <= 0 {
					result[idx] = 0
				}
			}
		}
	}

	return &result, nil
}

func (r *EkspeditorPostgres) GetProductInTheCarWithDifference(ekspeditorId int, language string) ([]models.ReadProductForTradePoint, error) {
	type ProductWithCount struct {
		models.Product
		Count float64
	}

	var productsData []ProductWithCount

	err := r.db.Table("products").
		Select("products.*, eksp_data.count").
		Joins("LEFT JOIN eksp_data ON eksp_data.product_id = products.id AND eksp_data.ekspeditor_id = ?", ekspeditorId).
		Preload("Name").
		Preload("ProductType").
		Where("products.id > 0").
		Find(&productsData).Error

	if err != nil {
		return nil, err
	}

	var result []models.ReadProductForTradePoint

	productMap := make(map[int]int)

	for _, item := range productsData {
		//if item.Count <= 0 {
		//	continue
		//}

		var res models.ReadProductForTradePoint
		res.Id = item.Id
		res.Volume = item.ProductType.Volume
		res.Count = int(item.Count)
		res.Free = int(item.Count)
		res.Name = item.Name.Name

		if language == "ru" {
			res.Taste = item.TasteRu
			res.ProductType = item.ProductType.NameRu
		} else {
			res.Taste = item.TasteTm
			res.ProductType = item.ProductType.NameTm
		}

		result = append(result, res)
		productMap[item.Id] = len(result) - 1
	}

	if len(result) == 0 {
		return []models.ReadProductForTradePoint{}, nil
	}

	date := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC) // Лучше использовать Local или UTC согласованно
	date = date.AddDate(0, 0, -1)

	if date.Weekday() == time.Sunday {
		date = date.AddDate(0, 0, -1)
	}

	type SoldItem struct {
		ProductId int
		TotalSold int
	}

	var soldItems []SoldItem

	var orders []models.Order

	get_orders_request := r.db.Model(&models.Order{}).
		Preload("List").
		Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Where("orders.created_at = ?", date).
		Where("trade_points.ekspeditor_id = ?", ekspeditorId).
		Find(&orders)

	if get_orders_request.Error != nil {

	}

	for _, order := range orders {
		if order.Status == true {
			continue
		}
		for _, list_element := range order.List {
			is_have := false
			for idx, soldItems_element := range soldItems {
				if soldItems_element.ProductId == list_element.ProductId {
					soldItems[idx].TotalSold += list_element.Count
					is_have = true
					break
				}
			}

			if is_have == false {
				var data SoldItem
				data.ProductId = list_element.ProductId
				data.TotalSold = list_element.Count
				soldItems = append(soldItems, data)

			}
		}
	}

	if err != nil {
		return nil, err
	}

	for _, sold := range soldItems {
		if idx, found := productMap[sold.ProductId]; found {
			result[idx].Free = result[idx].Count - int(sold.TotalSold)
		}
	}

	return result, nil
}

func (r *EkspeditorPostgres) GetProductPriceInTheCarWithDifference(ekspeditorId int) (float64, error) {
	// Берём все позиции из eksp_data с ценой и типом продукта
	type productPriceData struct {
		Id            int
		Price         float64
		Count         int
		ProductTypeId int
		BlockCount    int
	}

	var productsData []productPriceData

	err := r.db.Table("eksp_data").
		Select(`eksp_data.product_id AS id,
		        eksp_data.count,
		        products.price,
		        products.product_type_id,
		        product_types.count AS block_count`).
		Joins("JOIN products ON products.id = eksp_data.product_id").
		Joins("JOIN product_types ON product_types.id = products.product_type_id").
		Where("eksp_data.ekspeditor_id = ?", ekspeditorId).
		Find(&productsData).Error
	if err != nil {
		return 0, err
	}

	var totalPrice float64

	for _, item := range productsData {
		if item.Count <= 0 {
			continue
		}

		// Для типов 4/5/6 используем округление блоков, как в SaveOrder
		if item.ProductTypeId == 4 || item.ProductTypeId == 5 || item.ProductTypeId == 6 || item.ProductTypeId == 8 {
			blockSize := item.BlockCount
			if blockSize <= 0 {
				blockSize = 1
			}
			blockPrice := RoundToHalf(float64(blockSize) * item.Price)
			blocks := item.Count / blockSize
			rest := item.Count % blockSize

			totalPrice += float64(blocks)*blockPrice + float64(rest)*item.Price
		} else {
			totalPrice += float64(item.Count) * item.Price
		}
	}

	// Округляем до двух знаков, как Round2
	totalPrice = Round2(totalPrice)

	return totalPrice, nil
}

func (r *EkspeditorPostgres) GetProductInTheCarWithDifferenceInArray(language string) ([]models.GetGalyndyInArray, error) {
	var ekspeditors []models.Ekspeditor

	var result []models.GetGalyndyInArray

	get_ekspeditors := r.db.Model(&models.Ekspeditor{}).Where("status = true").Find(&ekspeditors)

	if get_ekspeditors.Error != nil {

	}
	var products []models.Product

	products_request := r.db.Model(&models.Product{}).Order("id").Preload("Name").Preload("ProductType").
		Where("status = true").Find(&products)

	if products_request.Error != nil {
		fmt.Println("Product gelmedi", products_request.Error)
	}

	for _, product_item := range products {
		var temporary_product models.GetGalyndyInArray
		temporary_product.Id = product_item.Id
		temporary_product.ProductTypeId = product_item.ProductTypeId
		temporary_product.Volume = product_item.ProductType.Volume
		temporary_product.Name = product_item.Name.Name

		if language == "ru" {
			temporary_product.Taste = product_item.TasteRu
			temporary_product.ProductType = product_item.ProductType.NameRu
		} else {
			temporary_product.Taste = product_item.TasteTm
			temporary_product.ProductType = product_item.ProductType.NameTm
		}

		result = append(result, temporary_product)
	}

	for _, item := range ekspeditors {
		data, err := r.GetProductInTheCarWithDifference(item.Id, language)
		if err == nil {
			for _, data_item := range data {
				for result_idx, result_item := range result {
					if data_item.Id == result_item.Id {
						var temporary_data models.EkspeditorsGalyndyObject
						temporary_data.EkspeditorId = item.Id
						temporary_data.Count = data_item.Count
						temporary_data.Free = data_item.Free
						result[result_idx].List = append(result[result_idx].List, temporary_data)
					}
				}
			}
		}
	}

	if len(result) == 0 {
		result = []models.GetGalyndyInArray{}
	}

	return result, nil
}

func (r *EkspeditorPostgres) ADDExtraTradePointsForEkspeditor(ekspeditor_id, day int, trade_points []models.TradePoint) []models.TradePoint {
	var extradDays []models.ExtraTradePointDay

	get_extra_day_request := r.db.Model(&models.ExtraTradePointDay{}).
		Joins("JOIN trade_points ON trade_points.id = extra_trade_point_days.trade_point_id").
		Preload("TradePoint").
		Where("extra_trade_point_days.day_id = ?", day).
		Where("trade_points.ekspeditor_id = ?", ekspeditor_id).
		Where("trade_points.status = ?", true).
		Find(&extradDays)

	if get_extra_day_request.Error != nil {
		return trade_points
	}

	for _, item := range extradDays {
		var trade_point models.TradePoint

		get_trade_point_request := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").
			Where("id = ?", item.TradePointId).Find(&trade_point)

		if get_trade_point_request.Error != nil {
			return trade_points
		}

		exists := false

		for _, element := range trade_points {
			if element.Id == trade_point.Id {
				exists = true
				break
			}
		}

		if exists == false {
			trade_points = append(trade_points, trade_point)
		}

	}

	return trade_points
}

func (r *EkspeditorPostgres) ADDExtraOrderPointsForEkpeditor(day, month, year, ekspeditor_id int, trade_points []models.TradePoint) []models.TradePoint {
	var orders []models.Order
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	get_orders_request := r.db.Model(&models.Order{}).
		Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("TradePoint").
		Where("created_at = ?", crDate).
		Where("trade_points.ekspeditor_id = ?", ekspeditor_id).
		Find(&orders)

	if get_orders_request.Error != nil {
		return trade_points
	}
	var new_trade_point_ids []int

	for _, order := range orders {
		count := 0
		for _, trade_point := range trade_points {
			if order.TradePointId == trade_point.Id {
				count++
			}
		}

		if count == 0 {
			new_trade_point_ids = append(new_trade_point_ids, order.TradePointId)
		}
	}

	if len(new_trade_point_ids) > 0 {

		for _, item := range new_trade_point_ids {
			var trade_point models.TradePoint

			get_trade_point_request := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").
				Where("id = ?", item).Find(&trade_point)

			if get_trade_point_request.Error != nil {
				return trade_points
			}

			exists := false

			for _, element := range trade_points {
				if element.Id == trade_point.Id {
					exists = true
					break
				}
			}

			if exists == false {
				trade_points = append(trade_points, trade_point)
			}

		}

	}

	return trade_points
}
