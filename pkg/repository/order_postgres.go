package repository

import (
	"bytes"
	"errors"

	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type OrderPostgres struct {
	db *gorm.DB
}

func NewOrderPostgres(db *gorm.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) GetTradePointsForMarketing(day, month, year, agentid int, language string) (*[]models.TradePointMobile, error) {
	var tradePoints []models.TradePoint
	var result []models.TradePointMobile

	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	request := r.db.Model(models.TradePoint{}).Preload("City").
		Preload("District").
		Where("trade_agent_id = ?", agentid).
		Where("day_id = ?", dayOfWeek).
		Where("status = ?", true).
		Find(&tradePoints)
	if request.Error != nil {
		return nil, request.Error
	}

	var routes []models.RouteOrder

	route_order_get_request := r.db.Model(&models.RouteOrder{}).Order("order_number").
		Where("trade_agent_id = ?", agentid).Where("day_id = ?", dayOfWeek).Find(&routes)

	if route_order_get_request.Error != nil {
		fmt.Println(route_order_get_request.Error)
	}

	for _, tradePoint := range tradePoints {
		var data models.TradePointMobile
		var contact contact_object

		contact = r.GetTredePointContact(tradePoint.Id, language)

		var count int64
		count = 0
		request3 := r.db.Model(&models.Marketing{}).Where("trade_point_id = ?", tradePoint.Id).
			Where("DATE(created_at) = ?", date.Format("2006-01-02")).Count(&count)

		if count == 0 {
			data.Status = false
		} else {
			data.Status = true
		}

		if request3.Error != nil {
			data.Status = false
		}

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

func (r *OrderPostgres) GetTradePoints(day, month, year, agentid int, language string) (*[]models.TradePointMobile, error) {
	var tradePoints []models.TradePoint
	var result []models.TradePointMobile
	paymentRepo := PaymentPostgres{db: r.db}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}
	//today := time.Now().Truncate(24 * time.Hour)
	request := r.db.Model(models.TradePoint{}).Preload("City").Preload("District").Where("status = ?", true).
		Where("trade_agent_id = ?", agentid).Where("day_id", dayOfWeek).Find(&tradePoints)

	if request.Error != nil {
		return nil, request.Error
	}

	tradePoints = r.ADDExtraTradePoints(agentid, dayOfWeek, tradePoints)
	tradePoints = r.ADDExtraOrderPoints(day, month, year, agentid, tradePoints)

	var routes []models.RouteOrder

	route_order_get_request := r.db.Model(&models.RouteOrder{}).Order("order_number").
		Where("trade_agent_id = ?", agentid).Where("day_id = ?", dayOfWeek).Find(&routes)

	if route_order_get_request.Error != nil {
		fmt.Println(route_order_get_request.Error)
	}

	for _, tradePoint := range tradePoints {
		var data models.TradePointMobile
		var count int64
		var contact contact_object
		var note models.PointNote
		var current_trade_point_debt float64
		current_trade_point_debt = paymentRepo.GetTradePointAllDebtSum(tradePoint.Id)

		get_note_request := r.db.Model(&models.PointNote{}).Where("trade_point_id = ?", tradePoint.Id).First(&note)

		if get_note_request.Error != nil {

		}

		contact = r.GetTredePointContact(tradePoint.Id, language)
		count = 0
		request3 := r.db.Model(&models.Order{}).Where("trade_point_id = ?", tradePoint.Id).
			Where("DATE(created_at) = ?", date.Format("2006-01-02")).Count(&count)

		if count == 0 {
			data.Status = false
		} else {
			data.Status = true
		}

		if request3.Error != nil {
			data.Status = false
		}

		data.Id = tradePoint.Id
		data.Code = tradePoint.Code
		data.Contact = contact.Name
		data.Number = contact.Number
		//data.Status = dataStatus
		data.Debt = current_trade_point_debt

		if len(note.Description) > 0 {
			data.Note = note.Description
		}

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

func (r *OrderPostgres) GetProducts(language string) (*[]models.ProductDistinct, error) {
	var result []models.ProductDistinct
	var products []models.Product

	request := r.db.Model(&models.Product{}).Preload("Name").Preload("ProductType").
		Where("id > ?", 0).Where("status = ?", true).
		Order("id").Find(&products)
	if request.Error != nil {
		return nil, request.Error
	}
	groupIndex := 1

	for idx, product := range products {
		if product.Id > 1 {
			if product.ProductTypeId != products[idx-1].ProductTypeId || product.BrandId != products[idx-1].BrandId {
				groupIndex++
			}
		}

		var data models.ProductDistinct
		data.Id = product.Id
		data.Code = product.Code
		data.Volume = product.ProductType.Volume
		data.Name = product.Name.Name
		data.GroupId = groupIndex
		data.Block_Count = product.ProductType.Count
		data.Price = product.Price
		if language == "ru" {
			data.Taste = product.TasteRu
			data.ProductTypeName = product.ProductType.NameRu
		} else {
			data.Taste = product.TasteTm
			data.ProductTypeName = product.ProductType.NameTm
		}

		result = append(result, data)
	}

	return &result, nil
}

func (r *OrderPostgres) GetBrand(language string) (*[]string, error) {
	result := make([]string, 0)
	var products []models.Product
	request := r.db.Model(&models.Product{}).Preload("Name").
		Where("id > ?", 0).
		Distinct("name_ru", "name_tm").
		Find(&products)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, product := range products {
		result = append(result, product.Name.Name)
	}

	return &result, nil
}

func (r *OrderPostgres) GetModel(brandId, ProductTypeId int, language string) (*[]string, error) {
	result := make([]string, 0)
	var products []models.Product
	request := r.db.Model(&models.Product{}).
		Where("brand_id = ? AND product_type_id = ?", brandId, ProductTypeId).
		Distinct("taste_ru", "taste_tm").
		Find(&products)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, product := range products {
		if language == "ru" {
			result = append(result, product.TasteRu)
		} else {
			result = append(result, product.TasteTm)
		}
	}

	return &result, nil
}

func (r *OrderPostgres) GetTypes(language string) (*[]string, error) {
	result := make([]string, 0)
	var productTypes []models.ProductType
	request := r.db.Model(&models.ProductType{}).
		Where("id > ?", 0).
		Distinct("name_ru", "name_tm").
		Find(&productTypes)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, productType := range productTypes {
		if language == "ru" {
			result = append(result, productType.NameRu)
		} else {
			result = append(result, productType.NameTm)
		}
	}
	return &result, nil
}

func (r *OrderPostgres) GetVolume(brand, model, tip string) (*[]models.SelectObject, error) {
	var result []models.SelectObject

	var products []models.Product
	request :=
		r.db.Model(&models.Product{}).Preload("ProductType").
			Where("name_ru LIKE ? OR name_tm LIKE ?", "%"+brand+"%", "%"+brand+"%").
			Where("taste_ru LIKE ? Or taste_tm LIKE ?", "%"+model+"%", "%"+model+"%").
			Find(&products)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, product := range products {
		if product.ProductType.NameRu == tip {
			var data models.SelectObject
			data.Id = product.Id
			data.Name = fmt.Sprintf("%.2f", product.ProductType.Volume)

			result = append(result, data)
		} else {
			fmt.Println("+++")
		}
	}

	if result == nil {
		result = make([]models.SelectObject, 0)
	}

	return &result, nil
}

func (r *OrderPostgres) GetSum(Id, Count int) (float64, error) {
	var product models.Product
	var sum float64

	request := r.db.Model(&models.Product{}).First(&product, Id)
	if request.Error != nil {
		return 0, request.Error
	}

	sum = roundToTwoDecimals((float64(Count) * product.Price))

	return sum, nil
}

func (r *OrderPostgres) saveOrder(
	day, month, year, tradePointId, paymentTypeId int,
	first, is_credit bool,
	orders []models.CreateOrder,
) error {

	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	var sum float64
	var order models.Order
	var tradePoint models.TradePoint

	// Торговая точка + канал
	request := r.db.Model(&models.TradePoint{}).
		Preload("TradeChannel").
		First(&tradePoint, tradePointId)

	if request.Error != nil || tradePoint.TradeAgentId <= 0 {
		if request.Error == nil {
			return errors.New("trade point has no trade agent")
		}
		return request.Error
	}

	// Ищем заказ на дату
	orderRequest := r.db.Model(&models.Order{}).
		Where("trade_point_id = ? AND created_at = ?", tradePointId, crDate).
		Find(&order)
	if orderRequest.Error != nil {
		return orderRequest.Error
	}

	if order.Id > 0 && order.Status == true {
		return errors.New("Rugsat berilmeyar")
	}

	// Если нет заказа — создаём
	if order.Id == 0 {
		order.TradePointId = tradePointId
		order.TradeAgentId = tradePoint.TradeAgentId
		order.Sum = 0
		order.SaleSum = 0
		order.CreatedAt = crDate
		order.Status = false
		order.IsClosed = false
		order.PaymentTypeId = paymentTypeId
		order.IsCredit = is_credit

		if err := r.db.Model(&models.Order{}).Create(&order).Error; err != nil {
			return err
		}
	}

	// Чистим старый OrderList
	if err := r.db.Where("order_id = ?", order.Id).Delete(&models.OrderList{}).Error; err != nil {
		return err
	}

	// Загружаем скидки + промо
	var discounts []models.Discount
	var promo []models.Promo

	if err := r.db.Model(&models.Discount{}).
		Order("sort_order").
		Where(`"id" > ?`, 1).
		Where(`"start" <= ?`, time.Now()).
		Where(`"end" >= ?`, time.Now()).
		Where(`"status" = ?`, true).
		Find(&discounts).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.Promo{}).
		Where(`"start" <= ?`, time.Now()).
		Where(`"end" >= ?`, time.Now()).
		Where(`"status" = ?`, true).
		Find(&promo).Error; err != nil {
		return err
	}

	// Загружаем продукты пачкой (цена + ProductTypeId)
	ids := make([]int, 0, len(orders))
	for _, o := range orders {
		ids = append(ids, o.Id)
	}

	var products []models.Product
	if len(ids) > 0 {
		if err := r.db.Model(&models.Product{}).
			Preload("ProductType").
			Where("id IN ?", ids).
			Find(&products).Error; err != nil {
			return err
		}
	}

	prodMap := make(map[int]models.Product, len(products))
	for _, p := range products {
		prodMap[p.Id] = p
	}

	// Формируем list
	list := make([]Products_Count_Price, 0, len(orders))
	for _, element := range orders {
		p := prodMap[element.Id]
		list = append(list, Products_Count_Price{
			Id:            element.Id,
			Count:         element.Count,
			OldCount:      element.Count,
			Price:         0,
			DiscountId:    1,
			Coefficient:   0,
			ProductTypeId: p.ProductTypeId,
		})
	}

	// Семейная блокировка: одна акция на один ProductTypeId
	familyChoice := map[int]int{}

	// 1) СНАЧАЛА СКИДКИ
	for _, discount := range discounts {
		if discount.IsDirect {
			list, _ = r.Type1Discount(discount.Id, tradePoint.TradeChannel.TypeId, tradePointId, tradePoint.TradeAgentId, list, familyChoice)
		} else {
			list, _ = r.Type2Discount(discount.Id, tradePoint.TradeChannel.TypeId, tradePointId, tradePoint.TradeAgentId, list, familyChoice)
		}
		// break НЕ делаем — разные скидки могут уйти на разные семьи
	}

	// 2) ПОТОМ PROMO (только на свободные семьи)
	for _, promoElement := range promo {
		list = r.CheckPromo(promoElement.Id, tradePoint.TradeChannel.TypeId, list, familyChoice)
	}

	// Сохраняем OrderList и считаем sum
	for _, element := range list {
		product, ok := prodMap[element.Id]
		if !ok {
			return fmt.Errorf("product not found for id %d", element.Id)
		}

		price := float64(element.OldCount) * product.Price

		var price2 float64
		if element.Count > 0 {
			if product.ProductTypeId == 4 || product.ProductTypeId == 5 || product.ProductTypeId == 6 || product.ProductTypeId == 8 {
				blockSize := product.ProductType.Count
				if blockSize <= 0 {
					blockSize = 1
				}
				blockPrice := math.Round(float64(blockSize) * product.Price)
				blocks := element.Count / blockSize
				price2 = (blockPrice * float64(blocks)) + element.Price
			} else {
				price2 = element.Price + (float64(element.Count) * product.Price)
			}
		} else {
			price2 = element.Price
		}

		price2 = Round2(price2)

		newOrderList := models.OrderList{
			ProductId:   element.Id,
			DiscountId:  element.DiscountId,
			OrderId:     order.Id,
			Count:       element.OldCount,
			Status:      true,
			Coefficient: element.Coefficient,
			Price:       price,
			SalePrice:   price2,
		}

		sum += newOrderList.SalePrice

		var discount_Count int64
		isHaveDiscount := r.db.Model(&models.Discount{}).Where("id = ?", newOrderList.DiscountId).Count(&discount_Count)

		if isHaveDiscount.Error != nil || discount_Count == 0 {
			newOrderList.DiscountId = 1
		}

		if err := r.db.Model(&models.OrderList{}).Create(&newOrderList).Error; err != nil {
			return err
		}
	}

	sum = Round2(sum)

	updateData := map[string]interface{}{
		"sum":             sum,
		"sale_sum":        sum,
		"is_credit":       is_credit,
		"payment_type_id": paymentTypeId,
	}

	return r.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(updateData).Error
}

func Round2(number float64) float64 {
	result := math.Floor(number*100) / 100

	return result
}

func (r *OrderPostgres) Marketing(day, month, year, tradePointId int, list []models.MarketingPost) error {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	for _, item := range list {
		var data models.Marketing

		request := r.db.Model(&models.Marketing{}).Where("trade_point_id = ? AND product_id = ? AND created_at = ?", tradePointId, item.ProductId, date).Find(&data)

		if request.Error != nil {
			fmt.Println("Marketing tapylmady ", request.Error)
		}

		if data.Id != 0 {
			updateData := map[string]interface{}{
				"have":     item.Have,
				"exposure": item.Exposure,
				"saled":    item.Saled,
			}

			updateRequest := r.db.Model(&models.Marketing{}).Where("id = ?", data.Id).Updates(updateData)

			if updateRequest.Error != nil {
				fmt.Println("Marketing maglumatlary uytgedip bolmady", updateRequest.Error)
			}

		} else {
			newMarketingData := models.Marketing{
				TradePointId: tradePointId,
				ProductId:    item.ProductId,
				Have:         item.Have,
				Exposure:     item.Exposure,
				Saled:        item.Saled,
				CreatedAt:    date,
			}

			savaReqeust := r.db.Model(&models.Marketing{}).Create(&newMarketingData)

			if savaReqeust.Error != nil {

			}
		}
	}
	return nil
}

func (r *OrderPostgres) OrderExcel(tradeAgentId, cityId, districtId int, name, language string) ([]byte, error) {
	var data []models.Order
	var list []models.ReadOrder

	request := r.db.Model(&models.Order{}).Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("TradePoint")

	if tradeAgentId > 0 {
		request = request.Where("trade_points.trade_agent_id = ?", tradeAgentId)
	}

	if len(name) > 0 {
		request = request.Where("trade_points.name_ru ILIKE ?", "%"+name+"%").Or("trade_points.name_tm ILIKE ?", "%"+name+"%")
	}

	if cityId > 0 {
		request = request.Where("trade_points.city_id = ?", cityId)
	}

	if districtId > 0 {
		request = request.Where("trade_points.district_id = ?", districtId)
	}

	var count int64

	request = request.Count(&count)

	request = request.Order("id desc").Find(&data)

	for _, element := range data {
		var item models.ReadOrder
		var tradePoint models.TradePoint
		item.Id = element.Id
		item.TradePointId = element.TradePointId
		item.TradeAgentId = element.TradePoint.TradeAgentId
		item.Sum = element.Sum
		item.CreatedAt = element.CreatedAt.Format("02.01.2006")

		pointReqeust := r.db.Model(&models.TradePoint{}).Preload("TradeAgent").
			Preload("City").Preload("District").Where("id = ?", element.TradePointId).Find(&tradePoint)

		if pointReqeust.Error != nil {
		}

		if language == "ru" {
			item.TradePoint = element.TradePoint.NameRu
			item.TradeAgent = tradePoint.TradeAgent.NameRu
			item.Location = tradePoint.LocationRu
			item.City = tradePoint.City.NameRu
			item.District = tradePoint.District.NameRu
		} else {
			item.TradePoint = element.TradePoint.NameTm
			item.TradeAgent = tradePoint.TradeAgent.NameTm
			item.Location = tradePoint.LocationTm
			item.City = tradePoint.City.NameTm
			item.District = tradePoint.District.NameTm
		}

		list = append(list, item)
	}

	//	fmt.Println(list)

	f := excelize.NewFile()
	sheet := "Заказы"
	f.NewSheet(sheet)

	// Define headers
	headers := []string{"Наименование торговой точки", "Торговый представитель", "Город", "Район", "Местонахождение", "Сумма заказа", "Оплачено", "Не оплачено", "Дата"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	// Populate data rows
	for rowIdx, order := range list {
		row := rowIdx + 2 // Start from the second row
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), order.TradePoint)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), order.TradeAgent)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), order.City)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), order.District)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), order.Location)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), order.Sum)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), order.Paid)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), order.NotPaid)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), order.CreatedAt)
	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *OrderPostgres) GetOrderPage(pageSize, pageNumber, tradeAgentId, cityId, districtId, day, month, year int, name, language, location, code, order string, status, is_status_active bool) (*models.OrderPage, error) {
	var data []models.Order
	var list []models.ReadOrder
	var paid_all, notPaid_all float64

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.Order{}).Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("TradePoint")

	if len(order) > 0 {
		sortOrder, err := allowedOrderSort(order)
		if err != nil {
			return nil, err
		}
		request = request.Order(sortOrder)
	}

	if is_status_active {
		request = request.Where("orders.status = ?", status)
	}

	if tradeAgentId > 0 {
		request = request.Where("trade_points.trade_agent_id = ?", tradeAgentId)
	}

	if len(name) > 0 {
		request = request.Where("trade_points.name_ru ILIKE ?", "%"+name+"%").Or("trade_points.name_tm ILIKE ?", "%"+name+"%")
	}

	if len(location) > 0 {
		request = request.Where("trade_points.location_ru ILIKE ?", "%"+location+"%").Or("trade_points.location_tm ILIKE ?", "%"+location+"%")
	}

	if cityId > 0 {
		request = request.Where("trade_points.city_id = ?", cityId)
	}

	if districtId > 0 {
		request = request.Where("trade_points.district_id = ?", districtId)
	}
	fmt.Println(code)
	if len(code) > 0 {
		request = request.Where("trade_points.code ILIKE ?", "%"+code+"%")
	}

	if day > 0 && month > 0 && year > 0 {
		crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		request = request.Where("created_at = ?", crDate)
	}

	var count int64

	request = request.Count(&count)

	request = request.Order("id desc").Offset(offset).Limit(pageSize).Find(&data)

	for _, element := range data {
		var item models.ReadOrder
		var tradePoint models.TradePoint
		var paid float64
		var payments []models.Payment

		get_payments_request := r.db.Model(&models.Payment{}).Where("order_id = ?", element.Id).Find(&payments)

		if get_payments_request.Error != nil {
			fmt.Println(get_payments_request.Error)
		}

		for _, payment_element := range payments {
			paid += payment_element.Currency
		}

		not_paid := element.Sum - paid
		paid = Round2(paid)
		not_paid = Round2(not_paid)

		item.Id = element.Id
		item.TradePointId = element.TradePointId
		item.TradeAgentId = element.TradePoint.TradeAgentId
		item.Sum = element.Sum
		item.CreatedAt = element.CreatedAt.Format("02.01.2006")
		item.Status = element.Status
		item.Paid = paid
		item.NotPaid = not_paid

		pointReqeust := r.db.Model(&models.TradePoint{}).Preload("TradeAgent").Preload("TradeCategory").
			Preload("City").Preload("District").Where("id = ?", element.TradePointId).Find(&tradePoint)

		if pointReqeust.Error != nil {
		}

		if language == "ru" {
			item.TradePoint = element.TradePoint.NameRu
			item.TradeAgent = tradePoint.TradeAgent.NameRu
			item.Location = tradePoint.LocationRu
			item.Orientir = tradePoint.OrientirRu
			item.City = tradePoint.City.NameRu
			item.District = tradePoint.District.NameRu
			item.TradeCategory = tradePoint.TradeCategory.NameRu
		} else {
			item.TradePoint = element.TradePoint.NameTm
			item.TradeAgent = tradePoint.TradeAgent.NameTm
			item.Location = tradePoint.LocationTm
			item.Orientir = tradePoint.OrientirTm
			item.City = tradePoint.City.NameTm
			item.District = tradePoint.District.NameTm
			item.TradeCategory = tradePoint.TradeCategory.NameTm
		}

		list = append(list, item)
		paid_all += paid
		notPaid_all += not_paid
	}

	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if list == nil {
		list = []models.ReadOrder{}
	}

	page := models.OrderPage{
		Items:      list,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
		Sum:        paid_all + notPaid_all,
		Paid:       paid_all,
		NotPaid:    notPaid_all,
	}

	return &page, nil
}

func (r *OrderPostgres) GetOrderList(OrderId int, language string) (*[]models.ReadOrderList, error) {
	var orderList []models.OrderList
	var result []models.ReadOrderList

	var request = r.db.Model(&models.OrderList{}).Preload("Discount").Where("order_id = ?", OrderId).Find(&orderList)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, item := range orderList {
		var data models.ReadOrderList
		var product models.Product

		var productRequest = r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Where("id = ?", item.ProductId).Find(&product)
		if productRequest.Error != nil {
		}
		data.Id = item.Id
		data.Count = item.Count
		data.Price = product.Price
		data.Status = item.Status
		data.Volume = product.ProductType.Volume
		data.Name = product.Name.Name
		if item.DiscountId > 1 {
			data.IsProcent = item.Discount.IsProcent
			data.Start = item.Discount.Start.Format("02.01.2006")
			data.End = item.Discount.End.Format("02.01.2006")
			data.Coefficient = item.Coefficient

			if data.IsProcent {
				data.DiscountPrice = product.Price - ((product.Price / 100) * float64(data.Coefficient))
			} else {
				data.DiscountPrice = product.Price - float64(data.Coefficient)/100
			}
		}

		data.DiscountPrice = Round2(data.DiscountPrice)

		if language == "ru" {
			data.Taste = product.TasteRu
		} else {
			data.Taste = product.TasteTm
		}
		result = append(result, data)
	}

	if result == nil {
		result = []models.ReadOrderList{}
	}

	return &result, nil
}

func (r *OrderPostgres) GetOrders(day, month, year, tradePointId int, language string) (models.GetOrderObject, error) {
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

func (r *OrderPostgres) CloseTheOrder(orderId int) error {
	var order models.Order

	getRequest := r.db.Model(&models.Order{}).Where("id = ?", orderId).Find(&order)
	if getRequest.Error != nil {
		return getRequest.Error
	}

	updateDate := map[string]interface{}{
		"is_closed": true,
	}

	updateRequest := r.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(updateDate)
	if updateRequest.Error != nil {
		return updateRequest.Error
	}

	return nil
}

type Products_Count_Price struct {
	Id            int
	Count         int
	OldCount      int
	Price         float64
	DiscountId    int
	Coefficient   int
	ProductTypeId int // семья = ProductTypeId
}

type contact_object struct {
	Name   string
	Number string
}

func (r *OrderPostgres) GetTredePointContact(trade_point_id int, language string) contact_object {
	var result contact_object
	result.Name = ""
	result.Number = ""
	var contact models.Contact

	get_contact_request := r.db.Model(&models.Contact{}).Where("trade_point_id = ?", trade_point_id).Order("kind_id").Find(&contact)

	if get_contact_request.Error != nil {
	} else {
		if language == "ru" {
			result.Name = contact.NameRu
		} else {
			result.Name = contact.NameTm
		}

		result.Number = contact.Number
	}

	return result
}

func checkGroupByProductType(dproducts []models.DProduct, list []Products_Count_Price) bool {
	if len(dproducts) == 0 {
		return false
	}

	groups := make(map[int][]models.DProduct)
	for _, dp := range dproducts {
		typeID := dp.Product.ProductTypeId
		groups[typeID] = append(groups[typeID], dp)
	}

	for _, group := range groups {
		if len(group) == 0 {
			continue
		}

		required := group[0].Count
		have := 0

		for _, dp := range group {
			for _, item := range list {
				if item.Id == dp.ProductId {
					have += item.Count
				}
			}
		}

		if have < required {
			return false
		}
	}

	return true
}

func markFamiliesFromDProducts(first []models.DProduct, familyChoice map[int]int, actionId int) {
	for _, dp := range first {
		typeId := dp.Product.ProductTypeId
		if typeId <= 0 {
			continue
		}

		// Назначаем выбор, только если семья ещё ничего не выбирала.
		// Если уже выбрала другую акцию — не трогаем.
		if chosen := familyChoice[typeId]; chosen == 0 {
			familyChoice[typeId] = actionId
		}
	}
}

func normalizedUnitPrice(p models.Product) float64 {
	price := p.Price

	if p.ProductTypeId != 4 && p.ProductTypeId != 5 && p.ProductTypeId != 6 {
		return price
	}

	bs := 1
	if p.ProductType.Count > 0 {
		bs = p.ProductType.Count
	}

	roundedBlock := RoundToHalf(float64(bs) * price) // 53.04 -> 53.0
	return roundedBlock / float64(bs)                // 53.0/12 -> 4.416666...
}

func (r *OrderPostgres) Type1Discount(
	DiscountId, ChannelTypeId, trade_point_id, trade_agent_id int,
	list []Products_Count_Price,
	familyChoice map[int]int,
) ([]Products_Count_Price, bool) {

	var discount models.Discount
	var dproduct_lists []models.DProduct

	isChannelTypeCorrect := false
	isTradePointCorrect := false
	isTradeAgentCorrect := false
	discountUsed := false

	discountRequest := r.db.Model(&models.Discount{}).
		Preload("ChannelTypes").
		Preload("TradePoints").
		Preload("TradeAgents").
		Where("id = ?", DiscountId).
		Find(&discount)
	if discountRequest.Error != nil {
		return list, false
	}

	for _, item := range discount.ChannelTypes {
		if item.Id == ChannelTypeId {
			isChannelTypeCorrect = true
			break
		}
	}

	if len(discount.TradeAgents) > 0 {
		for _, element := range discount.TradeAgents {
			if element.Id == trade_agent_id {
				isTradeAgentCorrect = true
				break
			}
		}
	} else {
		isTradeAgentCorrect = true
	}

	if len(discount.TradePoints) > 0 {
		for _, element := range discount.TradePoints {
			if element.Id == trade_point_id {
				isTradePointCorrect = true
				break
			}
		}
	} else {
		isTradePointCorrect = true
	}

	if !isChannelTypeCorrect || !isTradeAgentCorrect || !isTradePointCorrect {
		return list, false
	}

	// важно: ProductType нужен для Count (block size)
	dproductRequest := r.db.Model(&models.DProduct{}).
		Preload("Product").
		Preload("Product.ProductType").
		Where("discount_id = ?", DiscountId).
		Find(&dproduct_lists)
	if dproductRequest.Error != nil {
		return list, false
	}

	if discount.IsGroup {
		if ok := checkGroupByProductType(dproduct_lists, list); !ok {
			return list, false
		}
	}

	// isIncrease: если true — "увеличение", иначе скидка
	is_increase_element := 1
	if discount.IsIncrease {
		is_increase_element = -1
	}

	actionId := DiscountId

	for _, dp := range dproduct_lists {
		for idx, item := range list {
			if dp.ProductId != item.Id {
				continue
			}
			if item.Count <= 0 {
				continue
			}

			// семья выбрала другую акцию -> пропуск
			chosen := familyChoice[item.ProductTypeId]
			if chosen != 0 && chosen != actionId {
				continue
			}

			// не-group: порог по конкретному продукту, но скидка на весь Count
			if !discount.IsGroup && item.Count < dp.Count {
				continue
			}

			// семья свободна -> выбирает эту скидку
			if chosen == 0 {
				familyChoice[item.ProductTypeId] = actionId
			}

			// ВАЖНО: нормализуем базовую цену за 1 ед. (для газировки) ДО расчёта скидки
			baseUnit := normalizedUnitPrice(dp.Product)

			var adjustment float64
			if discount.IsProcent {
				adjustment = (baseUnit / 100) * float64(dp.Coefficient)
			} else {
				// оставил как у тебя: coefficient/100
				// если "фикс" должен быть в манатах — тут нужно менять логику хранения coefficient
				adjustment = float64(dp.Coefficient) / 100
			}

			coefficientPrice := baseUnit - (adjustment * float64(is_increase_element))

			// скидка на ВЕСЬ item.Count
			list[idx].Price += float64(item.Count) * coefficientPrice
			list[idx].Count = 0
			list[idx].DiscountId = DiscountId
			list[idx].Coefficient = dp.Coefficient

			discountUsed = true
		}
	}

	return list, discountUsed
}

func (r *OrderPostgres) Type2Discount(
	DiscountId, ChannelTypeId, trade_point_id, trade_agent_id int,
	list []Products_Count_Price,
	familyChoice map[int]int,
) ([]Products_Count_Price, bool) {

	var discount models.Discount
	var first []models.DProduct
	var second []models.DProduct

	listcopy := make([]Products_Count_Price, len(list))
	copy(listcopy, list)

	isChannelTypeCorrect := false
	isTradePointCorrect := false
	isTradeAgentCorrect := false
	discountUsed := false

	discountRequest := r.db.Model(&models.Discount{}).
		Preload("ChannelTypes").
		Preload("TradeAgents").
		Preload("TradePoints").
		Where("id = ?", DiscountId).
		Find(&discount)
	if discountRequest.Error != nil {
		return list, false
	}

	for _, item := range discount.ChannelTypes {
		if item.Id == ChannelTypeId {
			isChannelTypeCorrect = true
			break
		}
	}

	if len(discount.TradeAgents) > 0 {
		for _, element := range discount.TradeAgents {
			if element.Id == trade_agent_id {
				isTradeAgentCorrect = true
				break
			}
		}
	} else {
		isTradeAgentCorrect = true
	}

	if len(discount.TradePoints) > 0 {
		for _, element := range discount.TradePoints {
			if element.Id == trade_point_id {
				isTradePointCorrect = true
				break
			}
		}
	} else {
		isTradePointCorrect = true
	}

	if !isChannelTypeCorrect || !isTradeAgentCorrect || !isTradePointCorrect {
		return list, false
	}

	// важно: ProductType нужен для Count (block size)
	if err := r.db.Model(&models.DProduct{}).
		Preload("Product").
		Preload("Product.ProductType").
		Where("discount_id = ?", DiscountId).
		Where("type = ?", true).
		Find(&first).Error; err != nil {
		fmt.Println("Type2Discount: first dproducts not found:", err)
		return list, false
	}

	if err := r.db.Model(&models.DProduct{}).
		Preload("Product").
		Preload("Product.ProductType").
		Where("discount_id = ?", DiscountId).
		Where("type = ?", false).
		Find(&second).Error; err != nil {
		return list, false
	}

	actionId := DiscountId

	// 1) обязательная часть
	if discount.IsGroup {
		ok := checkGroupByProductType(first, list)
		if !ok {
			return list, false
		}

		// Семьи из first выбирают эту скидку, если ещё не выбирали
		markFamiliesFromDProducts(first, familyChoice, actionId)

	} else {
		// не-group: списываем first поштучно (по базовой цене), при этом семья должна быть свободна или уже выбрать эту скидку
		breakForLoop := false
		firstMarker := 0

		for _, dp := range first {
			found := false

			for idx, it := range list {
				if dp.ProductId != it.Id {
					continue
				}

				chosen := familyChoice[it.ProductTypeId]
				if chosen != 0 && chosen != actionId {
					breakForLoop = true
					break
				}

				if dp.Count > it.Count {
					breakForLoop = true
					break
				}

				// ВАЖНО: нормализуем базовую цену за 1 ед. (для газировки) ДО расчёта обязательной части
				baseUnit := normalizedUnitPrice(dp.Product)

				// обязательная часть по базовой цене
				list[idx].Price += float64(dp.Count) * baseUnit
				list[idx].Count = it.Count - dp.Count

				// семья выбирает эту скидку, если ещё не выбрала
				if chosen == 0 {
					familyChoice[it.ProductTypeId] = actionId
				}

				firstMarker++
				found = true
				break // критично: один dp должен матчиться ровно в одну позицию
			}

			if breakForLoop || !found {
				breakForLoop = true
				break
			}
		}

		if breakForLoop || firstMarker != len(first) {
			return listcopy, false
		}
	}

	// 2) скидочная часть second (скидка на весь Count)
	is_increase_element := 1
	if discount.IsIncrease {
		is_increase_element = -1
	}

	for _, dp := range second {
		for idx, it := range list {
			if dp.ProductId != it.Id {
				continue
			}
			if it.Count <= 0 {
				continue
			}

			chosen := familyChoice[it.ProductTypeId]
			if chosen != 0 && chosen != actionId {
				continue
			}

			// семья свободна -> выбирает эту скидку
			if chosen == 0 {
				familyChoice[it.ProductTypeId] = actionId
			}

			// ВАЖНО: нормализуем базовую цену за 1 ед. (для газировки) ДО расчёта скидки
			baseUnit := normalizedUnitPrice(dp.Product)

			var adjustment float64
			if discount.IsProcent {
				adjustment = (baseUnit / 100) * float64(dp.Coefficient)
			} else {
				adjustment = float64(dp.Coefficient) / 100
			}

			coefficientPrice := baseUnit - (adjustment * float64(is_increase_element))

			list[idx].Price += float64(it.Count) * coefficientPrice
			list[idx].Count = 0
			list[idx].DiscountId = DiscountId
			list[idx].Coefficient = dp.Coefficient

			discountUsed = true
		}
	}

	return list, discountUsed
}

func (r *OrderPostgres) CheckPromo(
	promoId, channelTypeId int,
	list []Products_Count_Price,
	familyChoice map[int]int,
) []Products_Count_Price {

	var promo models.Promo
	if err := r.db.Model(&models.Promo{}).
		Preload("ChannelTypes").
		Where("id = ?", promoId).
		First(&promo).Error; err != nil {
		return list
	}

	isChannelTypeCorrect := false
	for _, channelType := range promo.ChannelTypes {
		if channelTypeId == channelType.Id {
			isChannelTypeCorrect = true
			break
		}
	}
	if !isChannelTypeCorrect {
		return list
	}

	var needs []models.PromoProduct
	var frees []models.PromoProduct

	if err := r.db.Where("promo_id = ? AND is_needed = ?", promoId, true).Find(&needs).Error; err != nil {
		return list
	}
	if err := r.db.Where("promo_id = ? AND is_needed = ?", promoId, false).Find(&frees).Error; err != nil {
		return list
	}

	// собираем ids needs+frees
	idSet := map[int]struct{}{}
	for _, n := range needs {
		idSet[n.ProductId] = struct{}{}
	}
	for _, f := range frees {
		idSet[f.ProductId] = struct{}{}
	}

	ids := make([]int, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}

	var products []models.Product
	if len(ids) > 0 {
		if err := r.db.Model(&models.Product{}).
			Preload("ProductType").
			Where("id IN ?", ids).
			Find(&products).Error; err != nil {
			return list
		}
	}

	type prodInfo struct {
		Price         float64
		BlockSize     int
		ProductTypeId int
	}
	priceMap := make(map[int]prodInfo, len(products))
	for _, p := range products {
		bs := p.ProductType.Count
		if bs <= 0 {
			bs = 1
		}
		priceMap[p.Id] = prodInfo{
			Price:         p.Price,
			BlockSize:     bs,
			ProductTypeId: p.ProductTypeId,
		}
	}

	actionId := -promoId

	// 1) totalBlocks по needs: семья должна быть свободна ИЛИ уже выбрать это промо
	totalBlocks := 0
	for _, need := range needs {
		pi, ok := priceMap[need.ProductId]
		if !ok {
			continue
		}

		for _, item := range list {
			if item.Id != need.ProductId || item.Count <= 0 {
				continue
			}

			chosen := familyChoice[item.ProductTypeId]
			if chosen != 0 && chosen != actionId {
				continue
			}

			totalBlocks += item.Count / pi.BlockSize
		}
	}

	if totalBlocks < promo.NeedCount+promo.FreeCount {
		return list
	}

	promoSetSize := promo.NeedCount + promo.FreeCount
	if promoSetSize <= 0 {
		return list
	}

	times := totalBlocks / promoSetSize
	totalFreeBlocks := times * promo.FreeCount
	if totalFreeBlocks <= 0 {
		return list
	}

	// 2) применяем на frees: семья свободна ИЛИ уже выбрала это промо
	for i := 0; i < len(list) && totalFreeBlocks > 0; i++ {
		if list[i].Count <= 0 {
			continue
		}

		chosen := familyChoice[list[i].ProductTypeId]
		if chosen != 0 && chosen != actionId {
			continue
		}

		for _, free := range frees {
			if list[i].Id != free.ProductId {
				continue
			}

			pi, ok := priceMap[list[i].Id]
			if !ok || pi.BlockSize <= 0 {
				continue
			}

			blocksInProduct := list[i].Count / pi.BlockSize
			if blocksInProduct <= 0 {
				continue
			}

			giveBlocks := blocksInProduct
			if giveBlocks > totalFreeBlocks {
				giveBlocks = totalFreeBlocks
			}

			var discountAmount float64
			if pi.ProductTypeId == 4 || pi.ProductTypeId == 5 || pi.ProductTypeId == 6 || pi.ProductTypeId == 8 {
				blockPrice := math.Round(float64(pi.BlockSize) * pi.Price)
				discountAmount = blockPrice * float64(giveBlocks)
			} else {
				discountAmount = float64(giveBlocks*pi.BlockSize) * pi.Price
			}

			// семья свободна -> выбирает это промо
			if chosen == 0 {
				familyChoice[list[i].ProductTypeId] = actionId
				chosen = actionId
			}

			list[i].Price -= discountAmount
			list[i].DiscountId = actionId
			list[i].Coefficient = 0

			totalFreeBlocks -= giveBlocks
			if totalFreeBlocks <= 0 {
				break
			}
		}
	}

	return list
}

func (r *OrderPostgres) ADDExtraTradePoints(trade_agent_id, day int, trade_points []models.TradePoint) []models.TradePoint {
	var extradDays []models.ExtraTradePointDay

	get_extra_day_request := r.db.Model(&models.ExtraTradePointDay{}).
		Joins("JOIN trade_points ON trade_points.id = extra_trade_point_days.trade_point_id").
		Preload("TradePoint").
		Where("extra_trade_point_days.day_id = ?", day).
		Where("trade_points.trade_agent_id = ?", trade_agent_id).
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

func (r *OrderPostgres) ADDExtraOrderPoints(day, month, year, trade_agent_id int, trade_points []models.TradePoint) []models.TradePoint {
	var orders []models.Order
	crDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	get_orders_request := r.db.Model(&models.Order{}).Where("trade_agent_id = ? AND created_at = ?", trade_agent_id, crDate).Find(&orders)

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
