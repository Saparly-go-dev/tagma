package repository

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type DealerOrderPostgres struct {
	db *gorm.DB
}

func NewDealerOrderPostgres(db *gorm.DB) *DealerOrderPostgres {
	return &DealerOrderPostgres{db: db}
}

func (r *DealerOrderPostgres) Save(DealerId int, orders_list []models.CreateOrder) error {
	newDealerOrder := models.DealerOrder{
		DealerId:       DealerId,
		Sum:            0,
		SaleSum:        0,
		CreatedAt:      time.Now(),
		Status:         false,
		IsClosed:       false,
		Completed_date: time.Now(),
	}

	if err := r.db.Model(&models.DealerOrder{}).Create(&newDealerOrder).Error; err != nil {
		fmt.Println(err)
		return err
	}

	channelTypeId := 6

	// скидки
	var discounts []models.Discount
	if err := r.db.Model(&models.Discount{}).
		Where(`"id" > ?`, 1).
		Where(`"start" <= ?`, time.Now()).
		Where(`"end" >= ?`, time.Now()).
		Where(`"status" = ?`, true).
		Find(&discounts).Error; err != nil {
		fmt.Println("discount tapylmady", err)
	}

	// промо
	var promos []models.Promo
	if err := r.db.Model(&models.Promo{}).
		Where(`"start" <= ?`, time.Now()).
		Where(`"end" >= ?`, time.Now()).
		Where(`"status" = ?`, true).
		Find(&promos).Error; err != nil {
		fmt.Println("promo tapylmady", err)
	}

	// продукты пачкой (нужны DPrice + ProductType)
	ids := make([]int, 0, len(orders_list))
	for _, o := range orders_list {
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

	// list с ProductTypeId
	list := make([]Products_Count_Price, 0, len(orders_list))
	for _, element := range orders_list {
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

	// выбор акции по семье
	familyChoice := map[int]int{}

	// 1) скидки
	for _, discount := range discounts {
		if discount.IsDirect {
			list, _ = r.Type1Discount(discount.Id, channelTypeId, list, familyChoice)
		} else {
			list, _ = r.Type2Discount(discount.Id, channelTypeId, list, familyChoice)
		}
	}

	// 2) промо
	for _, p := range promos {
		list = r.CheckPromoDealer(p.Id, channelTypeId, list, familyChoice)
	}

	// сохраняем DealerOrderList
	var sum float64

	for _, item := range list {
		product, ok := prodMap[item.Id]
		if !ok {
			fmt.Println("product not found:", item.Id)
			continue
		}

		price := float64(item.OldCount) * product.DPrice

		var price2 float64
		if item.Count > 0 {
			if product.ProductTypeId == 4 || product.ProductTypeId == 5 || product.ProductTypeId == 6 || product.ProductTypeId == 8 {
				blockSize := product.ProductType.Count
				if blockSize <= 0 {
					blockSize = 1
				}
				bloc_baha := RoundToHalf(float64(blockSize) * product.DPrice)
				bloc_sany := item.Count / blockSize
				price2 = (bloc_baha * float64(bloc_sany)) + item.Price
			} else {
				price2 = item.Price + (float64(item.Count) * product.DPrice)
			}
		} else {
			price2 = item.Price
		}

		price = Round2(price)
		price2 = Round2(price2)

		newItem := models.DealerOrderList{
			ProductId:     item.Id,
			DiscountId:    item.DiscountId,
			DealerOrderId: newDealerOrder.Id,
			Price:         price,
			SalePrice:     price2,
			Coefficient:   item.Coefficient,
			Count:         item.OldCount,
			Status:        true,
		}

		sum += newItem.SalePrice

		var discount_Count int64
		isHaveDiscount := r.db.Model(&models.Discount{}).Where("id = ?", newItem.DiscountId).Count(&discount_Count)

		if isHaveDiscount.Error != nil || discount_Count == 0 {
			newItem.DiscountId = 1
		}

		if err := r.db.Model(&models.DealerOrderList{}).Create(&newItem).Error; err != nil {
			fmt.Println(err)
			continue
		}
	}

	sum = Round2(sum)

	updateData := map[string]interface{}{
		"sum":      sum,
		"sale_sum": sum,
	}

	return r.db.Model(&models.DealerOrder{}).Where("id = ?", newDealerOrder.Id).Updates(updateData).Error
}

func (r *DealerOrderPostgres) Update(DealerOrderId int, orders_list []models.CreateOrder) error {
	var dealer_order models.DealerOrder

	if err := r.db.Model(&models.DealerOrder{}).Where("id = ?", DealerOrderId).First(&dealer_order).Error; err != nil || dealer_order.Id == 0 {
		fmt.Println(err)
		return errors.New("Tapylmady")
	}

	if dealer_order.Status {
		return errors.New("Rugsat berilmeýär")
	}

	// чистим старый список
	if err := r.db.Where("dealer_order_id = ?", DealerOrderId).Delete(&models.DealerOrderList{}).Error; err != nil {
		return err
	}

	channelTypeId := 6

	// скидки
	var discounts []models.Discount
	if err := r.db.Model(&models.Discount{}).
		Where(`"id" > ?`, 1).
		Where(`"start" <= ?`, time.Now()).
		Where(`"end" >= ?`, time.Now()).
		Where(`"status" = ?`, true).
		Find(&discounts).Error; err != nil {
		fmt.Println("discount tapylmady", err)
	}

	// промо
	var promos []models.Promo
	if err := r.db.Model(&models.Promo{}).
		Where(`"start" <= ?`, time.Now()).
		Where(`"end" >= ?`, time.Now()).
		Where(`"status" = ?`, true).
		Find(&promos).Error; err != nil {
		fmt.Println("promo tapylmady", err)
	}

	// продукты пачкой
	ids := make([]int, 0, len(orders_list))
	for _, o := range orders_list {
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

	// list
	list := make([]Products_Count_Price, 0, len(orders_list))
	for _, element := range orders_list {
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

	familyChoice := map[int]int{}

	// 1) скидки
	for _, discount := range discounts {
		if discount.IsDirect {
			list, _ = r.Type1Discount(discount.Id, channelTypeId, list, familyChoice)
		} else {
			list, _ = r.Type2Discount(discount.Id, channelTypeId, list, familyChoice)
		}
	}

	// 2) промо
	for _, p := range promos {
		list = r.CheckPromoDealer(p.Id, channelTypeId, list, familyChoice)
	}

	// сохраняем
	var sum float64

	for _, item := range list {
		product, ok := prodMap[item.Id]
		if !ok {
			continue
		}

		price := float64(item.OldCount) * product.DPrice

		var price2 float64
		if item.Count > 0 {
			if product.ProductTypeId == 4 || product.ProductTypeId == 5 || product.ProductTypeId == 6 || product.ProductTypeId == 8 {
				blockSize := product.ProductType.Count
				if blockSize <= 0 {
					blockSize = 1
				}
				bloc_baha := RoundToHalf(float64(blockSize) * product.DPrice)
				bloc_sany := item.Count / blockSize
				price2 = (bloc_baha * float64(bloc_sany)) + item.Price
			} else {
				price2 = item.Price + (float64(item.Count) * product.DPrice)
			}
		} else {
			price2 = item.Price
		}

		price = Round2(price)
		price2 = Round2(price2)

		newItem := models.DealerOrderList{
			ProductId:     item.Id,
			DiscountId:    item.DiscountId,
			DealerOrderId: DealerOrderId,
			Price:         price,
			SalePrice:     price2,
			Coefficient:   item.Coefficient,
			Count:         item.OldCount,
			Status:        true,
		}

		sum += newItem.SalePrice

		var discount_Count int64
		isHaveDiscount := r.db.Model(&models.Discount{}).Where("id = ?", newItem.DiscountId).Count(&discount_Count)

		if isHaveDiscount.Error != nil || discount_Count == 0 {
			newItem.DiscountId = 1
		}

		if err := r.db.Model(&models.DealerOrderList{}).Create(&newItem).Error; err != nil {
			fmt.Println(err)
			continue
		}
	}

	sum = Round2(sum)

	updateData := map[string]interface{}{
		"sum":      sum,
		"sale_sum": sum,
	}

	return r.db.Model(&models.DealerOrder{}).Where("id = ?", DealerOrderId).Updates(updateData).Error
}

func (r *DealerOrderPostgres) CompleTheDealerOrder(dealerOrderId int) error {
	var dealer_order models.DealerOrder

	get_dealer_order_postgres := r.db.Model(&models.DealerOrder{}).Where("id = ?", dealerOrderId).First(&dealer_order)
	if get_dealer_order_postgres.Error != nil {
		return get_dealer_order_postgres.Error
	}

	update_data := map[string]interface{}{
		"completed_date": time.Now(),
	}

	update_request := r.db.Model(&models.DealerOrder{}).Where("id = ?", dealerOrderId).Updates(update_data)
	if update_request.Error != nil {
		return update_request.Error
	}

	return nil
}

func (r *DealerOrderPostgres) ChangeStatus(Id int) error {
	var data models.DealerOrder
	result := r.db.First(&data, Id)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	if data.Status == true {
		return nil

	} else {
		var dealer models.Dealer

		updateData := map[string]interface{}{
			"status": true,
		}
		update_request := r.db.Model(&models.DealerOrder{}).Where("id = ?", Id).Updates(updateData)
		if update_request.Error != nil {
			return update_request.Error
		}

		dealer_request := r.db.Model(&models.Dealer{}).Where("id =?", data.DealerId).Find(&dealer)

		if dealer_request.Error != nil {
			return dealer_request.Error
		}
		if dealer.Account >= data.Sum {
			// Eger sargydyň jemi tölegi dileriň hasabyndan kiçi bolsa onda sargyt tolenyar we dileriň hasabyndan degişli pul möçberi ayrylýar.
			update_dealer_account := map[string]interface{}{
				"account": Round2(dealer.Account - data.Sum),
			}

			update_dealer := r.db.Model(&models.Dealer{}).Where("id = ?", dealer.Id).Updates(update_dealer_account)
			if update_dealer.Error != nil {
				return update_dealer.Error
			}

			update_dealer_order := map[string]interface{}{
				"is_closed": true,
				"paid_date": time.Now(),
			}

			update_dealer_order_request := r.db.Model(&models.DealerOrder{}).Where("id = ?", data.Id).Updates(update_dealer_order)
			if update_dealer_order_request.Error != nil {
				return update_dealer_order_request.Error
			}

		}

	}

	return nil
}

func (r *DealerOrderPostgres) GetOrders(pageSize, pageNumber, DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) (models.DealerOrderPage, error) {
	var list []models.DealerOrder
	var items []models.ReadDealerOrder
	var result models.DealerOrderPage

	var start time.Time
	var end time.Time

	offset := (pageNumber - 1) * pageSize

	if endDay == 0 || endMonth == 0 || endYear == 0 {
		end = time.Now()
	} else {
		end = time.Date(endYear, end.Month(), endDay, 0, 0, 0, 0, time.UTC)
	}

	if startDay == 0 || startMonth == 0 || startYear == 0 {
		start = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		start = time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)
	}

	get_request := r.db.Model(&models.DealerOrder{}).Preload("Dealer").Order("id desc").
		Where(`"created_at" <= ?`, end).
		Where(`"created_at" >= ?`, start)

	if DealerId > 0 {
		get_request = get_request.Where("dealer_id = ?", DealerId)
	}

	var count int64
	get_request = get_request.Count(&count)
	get_request = get_request.Offset(offset).Limit(pageSize).Find(&list)

	for _, item := range list {
		var data models.ReadDealerOrder

		data.Id = item.Id
		data.DealerId = item.DealerId
		data.Status = item.Status
		data.IsClosed = item.IsClosed
		data.Sum = item.Sum
		data.CreatedAt = item.CreatedAt.Format("02.01.2006")
		data.PaidDate = item.Paid_date.Format("02.01.2006")
		data.CompletedDate = item.Completed_date.Format("02.01.2006")
		if data.PaidDate == "01.01.0001" {
			data.PaidDate = "-"
		}

		if data.CompletedDate == "01.01.0001" {
			data.CompletedDate = "-"
		}

		if language == "ru" {
			data.Dealer = item.Dealer.NameRu
		} else {
			data.Dealer = item.Dealer.NameTm
		}

		items = append(items, data)
	}

	if items == nil {
		items = []models.ReadDealerOrder{}
	}
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	result.Items = items
	result.PageNumber = pageNumber
	result.PageSize = pageSize
	result.PageCount = pageCount
	result.Total = int(count)

	return result, nil
}

func (r *DealerOrderPostgres) GetOrderList(Id int, language string) (*[]models.ReadOrderList, error) {
	var orderList []models.DealerOrderList
	var result []models.ReadOrderList

	var request = r.db.Model(&models.DealerOrderList{}).Preload("Discount").Where("dealer_order_id = ?", Id).Find(&orderList)

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
		data.Price = product.DPrice
		data.Status = item.Status
		data.Volume = product.ProductType.Volume
		data.Name = product.Name.Name
		data.ProductId = item.ProductId
		if item.DiscountId > 1 {
			data.IsProcent = item.Discount.IsProcent
			data.Start = item.Discount.Start.Format("02.01.2006")
			data.End = item.Discount.End.Format("02.01.2006")
			data.Coefficient = item.Coefficient

			if data.IsProcent {
				data.DiscountPrice = product.DPrice - ((product.DPrice / 100) * float64(data.Coefficient))
			} else {
				data.DiscountPrice = product.DPrice - float64(data.Coefficient)/100
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

func (r *DealerOrderPostgres) GetProduct(DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear int, status bool, language string) ([]models.ReadProductForTradePoint, error) {
	var dealer_orders []models.DealerOrder
	var result []models.ReadProductForTradePoint

	var start time.Time
	var end time.Time

	if endDay == 0 || endMonth == 0 || endYear == 0 {
		end = time.Now()
	} else {
		end = time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.UTC)
	}

	if startDay == 0 || startMonth == 0 || startYear == 0 {
		start = time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		start = time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)

	}

	get_request := r.db.Model(&models.DealerOrder{}).Preload("List").
		Where(`"created_at" <= ?`, end).
		Where(`"created_at" >= ?`, start)

	if DealerId > 0 {
		get_request = get_request.Where("dealer_id = ?", DealerId)
	}

	get_request = get_request.Find(&dealer_orders)
	fmt.Println(status)

	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Order("id").Where("id > 0")

	if status == true {

	} else {
		productsRequest = productsRequest.Where("status = ?", true)
	}

	productsRequest = productsRequest.Find(&products)

	if productsRequest.Error != nil {
		fmt.Println("Productlar tapylmady")
	}
	for _, element := range products {
		var data models.ReadProductForTradePoint
		data.Id = element.Id
		data.Name = element.Name.Name
		data.ProductTypeId = element.ProductType.Id
		data.Volume = element.ProductType.Volume
		if language == "ru" {
			data.Taste = element.TasteRu
			data.ProductType = element.ProductType.NameRu
		} else {
			data.Taste = element.TasteTm
			data.ProductType = element.ProductType.NameTm
		}

		result = append(result, data)
	}

	for _, order_element := range dealer_orders {
		for _, order_list_element := range order_element.List {
			for idx, product_element := range result {
				if order_list_element.ProductId == product_element.Id {
					result[idx].Count += order_list_element.Count
				}
			}
		}

	}

	if result == nil {
		result = []models.ReadProductForTradePoint{}
	}

	return result, nil

}

func (r *DealerOrderPostgres) DeleteDealerOrder(orderId int) error {
	var d_order models.DealerOrder

	get_dealer_order := r.db.Model(&models.DealerOrder{}).Where("id = ?", orderId).First(&d_order)

	if get_dealer_order.Error != nil {
		return errors.New("Not Found")
	}

	if d_order.Status == true {
		return errors.New("Rugsat berilmeýär")
	}

	delete_dealer_order_list := r.db.Model(&models.DealerOrderList{}).Where("dealer_order_id = ?", orderId).Delete(&[]models.DealerOrder{})

	if delete_dealer_order_list.Error != nil {
		return errors.New("Not Found")
	}

	delete_dealer_order := r.db.Model(&models.DealerOrder{}).Delete(&d_order)

	if delete_dealer_order.Error != nil {
		return errors.New("Not Found")
	}

	return nil
}

// ////Dealer-de ChannelTypeId 6-a deň
func normalizedDealerUnitPrice(p models.Product) float64 {
	price := p.DPrice

	if p.ProductTypeId != 4 && p.ProductTypeId != 5 && p.ProductTypeId != 6 {
		return price
	}

	bs := 1
	if p.ProductType.Count > 0 {
		bs = p.ProductType.Count
	}

	roundedBlock := RoundToHalf(float64(bs) * price) // например 53.04 -> 53.0
	return roundedBlock / float64(bs)                // 53.0/12 -> 4.416666...
}

// --- Dealer Type1Discount ---
func (r *DealerOrderPostgres) Type1Discount(
	DiscountId, ChannelTypeId int,
	list []Products_Count_Price,
	familyChoice map[int]int,
) ([]Products_Count_Price, bool) {

	var discount models.Discount
	var dproduct_lists []models.DProduct

	discountUsed := false
	isChannelTypeCorrect := false

	discountRequest := r.db.Model(&models.Discount{}).
		Preload("ChannelTypes").
		Where("id = ?", DiscountId).
		Find(&discount)
	if discountRequest.Error != nil {
		fmt.Println("type1-de discount tapylmady", discountRequest.Error)
		return list, false
	}

	for _, item := range discount.ChannelTypes {
		if item.Id == ChannelTypeId {
			isChannelTypeCorrect = true
			break
		}
	}
	if !isChannelTypeCorrect {
		return list, false
	}

	// важно: ProductType нужен для Count (block size)
	dproductRequest := r.db.Model(&models.DProduct{}).
		Preload("Product").
		Preload("Product.ProductType").
		Where("discount_id = ?", DiscountId).
		Find(&dproduct_lists)
	if dproductRequest.Error != nil {
		fmt.Println("type1-de dpproductlar tapylmady", dproductRequest.Error)
		return list, false
	}

	if discount.IsGroup {
		if ok := checkGroupByProductType(dproduct_lists, list); !ok {
			return list, false
		}
	}

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

			// ВАЖНО: нормализуем дилерскую цену за 1 ед. (для газировки) ДО расчёта скидки
			baseUnit := normalizedDealerUnitPrice(dp.Product)

			var adjustment float64
			if discount.IsProcent {
				adjustment = (baseUnit / 100) * float64(dp.Coefficient)
			} else {
				// оставил как у тебя: coefficient/100
				// если "фикс" должен быть в манатах — поменяй схему хранения coefficient
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

// --- Dealer Type2Discount ---
// / Dealer-de ÇhannelTypeId 6-a deň
func (r *DealerOrderPostgres) Type2Discount(
	DiscountId, ChannelTypeId int,
	list []Products_Count_Price,
	familyChoice map[int]int,
) ([]Products_Count_Price, bool) {

	var discount models.Discount
	var first []models.DProduct
	var second []models.DProduct

	discountUsed := false
	isChannelTypeCorrect := false

	listcopy := make([]Products_Count_Price, len(list))
	copy(listcopy, list)

	discountRequest := r.db.Model(&models.Discount{}).
		Preload("ChannelTypes").
		Where("id = ?", DiscountId).
		Find(&discount)
	if discountRequest.Error != nil {
		fmt.Println("type2-de discount tapylmady", discountRequest.Error)
		return list, false
	}

	for _, item := range discount.ChannelTypes {
		if item.Id == ChannelTypeId {
			isChannelTypeCorrect = true
			break
		}
	}
	if !isChannelTypeCorrect {
		return list, false
	}

	// важно: ProductType нужен для Count (block size)
	if err := r.db.Model(&models.DProduct{}).
		Preload("Product").
		Preload("Product.ProductType").
		Where("discount_id = ?", DiscountId).
		Where("type = ?", true).
		Find(&first).Error; err != nil {
		fmt.Println("type2-de first tapylmady", err)
		return list, false
	}

	if err := r.db.Model(&models.DProduct{}).
		Preload("Product").
		Preload("Product.ProductType").
		Where("discount_id = ?", DiscountId).
		Where("type = ?", false).
		Find(&second).Error; err != nil {
		fmt.Println("type2-de second tapylmady", err)
		return list, false
	}

	actionId := DiscountId

	// 1) обязательная часть
	if discount.IsGroup {
		if ok := checkGroupByProductType(first, list); !ok {
			return list, false
		}

		// семьи из first выбирают эту скидку (если ещё не выбрали)
		markFamiliesFromDProducts(first, familyChoice, actionId)

	} else {
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

				// ВАЖНО: нормализуем дилерскую цену за 1 ед. (для газировки) ДО расчёта обязательной части
				baseUnit := normalizedDealerUnitPrice(dp.Product)

				// обязательная часть по базовой цене
				list[idx].Price += float64(dp.Count) * baseUnit
				list[idx].Count = it.Count - dp.Count

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

	// 2) second — скидка на весь Count
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

			if chosen == 0 {
				familyChoice[it.ProductTypeId] = actionId
			}

			// ВАЖНО: нормализуем дилерскую цену за 1 ед. (для газировки) ДО расчёта скидки
			baseUnit := normalizedDealerUnitPrice(dp.Product)

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

func (r *DealerOrderPostgres) CheckPromoDealer(
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
	for _, ct := range promo.ChannelTypes {
		if ct.Id == channelTypeId {
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

	// ids needs+frees
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
		DPrice        float64
		BlockSize     int
		ProductTypeId int
	}
	priceMap := make(map[int]prodInfo, len(products))
	for _, p := range products {
		bs := p.ProductType.Count
		if bs <= 0 {
			bs = 1
		}
		priceMap[p.Id] = prodInfo{DPrice: p.DPrice, BlockSize: bs, ProductTypeId: p.ProductTypeId}
	}

	actionId := -promoId

	// totalBlocks по needs (семья свободна ИЛИ уже выбрала это промо)
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

	// применяем на frees (семья свободна ИЛИ уже выбрала это промо)
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
				blockPrice := RoundToHalf(float64(pi.BlockSize) * pi.DPrice)
				discountAmount = blockPrice * float64(giveBlocks)
			} else {
				discountAmount = float64(giveBlocks*pi.BlockSize) * pi.DPrice
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

func RoundToHalf(n float64) float64 {
	// Округляем до ближайшего кратного 0.5
	return math.Round(n*2) / 2
}
