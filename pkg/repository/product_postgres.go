package repository

import (
	"bytes"
	"math"
	"strconv"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ProductPostgres struct {
	db *gorm.DB
}

func NewProductPostgres(db *gorm.DB) *ProductPostgres {
	return &ProductPostgres{db: db}
}

// Create inserts a new city into the database and returns its ID
func (r *ProductPostgres) Create(product models.CreateProduct) (int, error) {

	var lastProduct models.Product
	var codeindex int
	var productPrice float64
	var dProductPrice float64

	lastresult := r.db.Model(models.Product{}).Last(&lastProduct).Order("id desc")

	if lastresult.Error != nil {
		codeindex = 1
	} else {
		codeindex = lastProduct.Id + 1
	}

	productPrice = float64(product.Price1) + float64(product.Price2)/100

	dProductPrice = float64(product.DPrice1) + float64(product.DPrice2)/100

	newProduct := models.Product{
		BrandId:       product.BrandId,
		TasteRu:       product.TasteRu,
		TasteTm:       product.TasteTm,
		Price:         productPrice,
		DPrice:        dProductPrice,
		Code:          strconv.Itoa(codeindex),
		Status:        product.Status,
		ProductTypeId: product.ProductTypeId,
		UpdatedAt:     time.Now(),
	}

	result := r.db.Create(&newProduct)
	if result.Error != nil {
		return 0, result.Error
	}

	return newProduct.Id, nil
}

func roundToTwoDecimals(num float64) float64 {
	return math.Round(num*100) / 100
}

// GetAll retrieves all cities from the database
func (r *ProductPostgres) GetAll(pageSize, pageNumber int, name, language string) (*models.ProductPage, error) {
	var products []models.Product
	var readProducts []models.ReadProduct

	// Calculate offset for pagination
	offset := (pageNumber - 1) * pageSize

	// Use correct join clause for product and product type, assuming a foreign key like `product_type_id` existsam
	result := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Where("taste_ru ILIKE ? OR taste_tm ILIKE ?", "%"+name+"%", "%"+name+"%").Order("id desc").
		Offset(offset).Limit(pageSize).Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	// Iterate through the retrieved products
	for _, product := range products {
		var newitem models.ReadProduct

		// Fill newitem with product details
		newitem.Id = product.Id
		newitem.Code = product.Code
		newitem.Name = product.Name.Name
		newitem.Price = product.Price
		newitem.Status = product.Status
		newitem.TasteRu = product.TasteRu
		newitem.TasteTm = product.TasteTm
		newitem.BrandId = product.BrandId
		newitem.DPrice = product.DPrice
		newitem.ProductTypeId = product.ProductTypeId
		newitem.BlockPrice = roundToTwoDecimals((float64(product.ProductType.Count) * product.Price))
		newitem.Volume = product.ProductType.Volume
		newitem.Updatedat = product.UpdatedAt.Format("02.01.2006")

		// Set ProductType based on language
		if language == "ru" {
			newitem.ProductType = product.ProductType.NameRu
		} else {
			newitem.ProductType = product.ProductType.NameTm
		}
		// Append to the result list
		readProducts = append(readProducts, newitem)
	}

	// Calculate total count for pagination
	var count int64
	result2 := r.db.Model(&models.Product{}).
		Where("taste_ru ILIKE ? OR taste_tm ILIKE ?", "%"+name+"%", "%"+name+"%").
		Count(&count)

	if result2.Error != nil {
		return nil, result2.Error
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if readProducts == nil {
		readProducts = []models.ReadProduct{}
	}

	// Create the page object
	page := models.ProductPage{
		Items:      readProducts,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
	}

	return &page, nil
}

func (r *ProductPostgres) ChangeStatus(Id int) error {
	var data models.Product
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.Product{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.Product{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

// GetById retrieves a city by its ID
func (r *ProductPostgres) GetById(ProductId int) (models.CreateProduct, error) {
	var product models.Product
	var data models.CreateProduct
	var price1 int
	var price2 int
	var dPrice1 int
	var dPrice2 int
	result := r.db.First(&product, ProductId)

	price1 = int(product.Price)
	price2 = (int(product.Price*100) % 100)
	dPrice1 = int(product.DPrice)
	dPrice2 = (int(product.DPrice*100) % 100)

	data.BrandId = product.BrandId
	data.TasteRu = product.TasteRu
	data.TasteTm = product.TasteTm
	data.Status = product.Status
	data.ProductTypeId = product.ProductTypeId
	data.Price1 = price1
	data.Price2 = price2
	data.DPrice1 = dPrice1
	data.DPrice2 = dPrice2

	return data, result.Error
}

func (r *ProductPostgres) GetTypes(language string) ([]models.SelectObject, error) {
	var productTypes []models.ProductType
	var result []models.SelectObject

	query := r.db.Where("id > ?", 0).Order("id").Find(&productTypes)

	if query.Error != nil {
		return nil, query.Error
	}

	for _, productType := range productTypes {
		var newitem models.SelectObject

		newitem.Id = productType.Id
		volume := productType.Volume
		volumeStr := ""

		if volume == float64(int(volume)) {
			volumeStr = strconv.Itoa(int(volume))
		} else {
			volumeStr = strconv.FormatFloat(volume, 'f', 1, 64)
		}
		volumeDisplay := volumeStr + "L"

		if language == "ru" {
			newitem.Name = productType.NameRu + " " + volumeDisplay
		} else {
			newitem.Name = productType.NameTm + " " + volumeDisplay
		}

		result = append(result, newitem)
	}

	return result, nil
}

func (r *ProductPostgres) GetProductBrands() ([]models.SelectObject, error) {
	var brands []models.Brand
	var result []models.SelectObject

	query := r.db.Model(&models.Brand{}).Where("id > ?", 0).Find(&brands)

	if query.Error != nil {
		return nil, query.Error
	}

	for _, brand := range brands {
		var newitem models.SelectObject
		newitem.Id = brand.Id
		newitem.Name = brand.Name
		result = append(result, newitem)
	}

	return result, nil
}

// Delete removes a city from the database by its ID
func (r *ProductPostgres) Delete(productId int) error {
	result := r.db.Delete(&models.Product{}, productId)
	return result.Error
}

// Update modifies a city's information based on its ID
func (r *ProductPostgres) Update(ProductId int, Product models.CreateProduct) error {
	var productPrice float64
	var dProductPrice float64

	productPrice = float64(Product.Price1) + float64(Product.Price2)/100

	dProductPrice = float64(Product.DPrice1) + float64(Product.DPrice2)/100

	updateData := map[string]interface{}{
		"brand_id":        Product.BrandId,
		"status":          Product.Status,
		"taste_ru":        Product.TasteRu,
		"taste_tm":        Product.TasteTm,
		"price":           productPrice,
		"d_price":         dProductPrice,
		"product_type_id": Product.ProductTypeId,
	}
	result := r.db.Model(&models.Product{}).Where("id = ?", ProductId).Updates(updateData)
	return result.Error
}

func (r *ProductPostgres) GetInformantionAboutSales(brandId, productTypeId, agentId, cityId, districtId, tradePointId, channelTypeId, startDay, startMonth, startYear, endDay, endMonth, endYear int, volume float64, product_taste, language string) ([]models.ReadProductForTradePoint, error) {
	var result []models.ReadProductForTradePoint
	var orders []models.Order
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

	request := r.db.Model(&models.Order{}).
		Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Joins("JOIN trade_channels ON trade_channels.id = trade_points.trade_channel_id").
		Preload("List")

	request = request.Order("orders.id desc").
		Where("orders.created_at <= ?", end).
		Where("orders.created_at >= ?", start)

	if agentId > 0 {
		request = request.Where("trade_points.trade_agent_id = ?", agentId)
	}

	if cityId > 0 {
		request = request.Where("trade_points.city_id = ?", cityId)
	}

	if districtId > 0 {
		request = request.Where("trade_points.district_id = ?", districtId)
	}

	if tradePointId > 0 {
		request = request.Where("trade_point_id = ?", tradePointId)
	}

	if channelTypeId > 0 {
		request = request.Where("trade_channels.type_id = ?", channelTypeId)
	}

	if err := request.Find(&orders).Error; err != nil {
		return nil, err
	}

	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}). //Joins("JOIN product_types ON product_types.id = products.product_type_id").
								Preload("ProductType").Preload("Name").Order("id").Where("id > 0")

	if brandId > 0 {
		productsRequest = productsRequest.Where("brand_id = ?", brandId)
	}

	if productTypeId > 0 {
		productsRequest = productsRequest.Where("product_type_id = ?", productTypeId)
	}

	if len(product_taste) > 0 {
		productsRequest = productsRequest.Where("taste_ru ILIKE ? OR taste_tm ILIKE ?", "%"+product_taste+"%", "%"+product_taste+"%")
	}

	if err := productsRequest.Find(&products).Error; err != nil {
		return nil, err
	}
	for _, element := range products {

		if volume > 0 {
			if volume != element.ProductType.Volume {
				continue
			}
		}

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

	for _, order_element := range orders {
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

func (r *ProductPostgres) GetDailySalesProductsInExcel(agentId, day, month, year int, language string) ([]byte, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Загружаем заказы
	var orders []models.Order
	orderQuery := r.db.Model(&models.Order{}).
		Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("List").
		Where("DATE(orders.created_at) = ?", date)

	if agentId > 0 {
		orderQuery = orderQuery.Where("trade_points.trade_agent_id = ?", agentId)
	}

	if err := orderQuery.Order("orders.id desc").Find(&orders).Error; err != nil {
		return nil, err
	}

	// Загружаем продукты
	var products []models.Product
	productQuery := r.db.Model(&models.Product{}).
		Preload("ProductType").
		Preload("Name").
		Order("id")

	if err := productQuery.Find(&products).Error; err != nil {
		return nil, err
	}

	// Заполняем структуру для отчета
	var reportData []models.ReadProductForTradePoint
	for _, p := range products {
		item := models.ReadProductForTradePoint{
			Id:            p.Id,
			Name:          p.Name.Name,
			ProductTypeId: p.ProductType.Id,
			Volume:        p.ProductType.Volume,
		}

		if language == "ru" {
			item.Taste = p.TasteRu
			item.ProductType = p.ProductType.NameRu
		} else {
			item.Taste = p.TasteTm
			item.ProductType = p.ProductType.NameTm
		}
		reportData = append(reportData, item)
	}

	// Считаем продажи
	for _, order := range orders {
		for _, ol := range order.List {
			for i, prod := range reportData {
				if ol.ProductId == prod.Id {
					reportData[i].Count += ol.Count
				}
			}
		}
	}

	if reportData == nil {
		reportData = []models.ReadProductForTradePoint{}
	}

	// Генерация Excel
	f := excelize.NewFile()
	sheet := "Продажи"
	f.NewSheet(sheet)

	headers := []string{"ID", "Название", "Тип продукта", "Вкус", "Объем", "Количество"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, item := range reportData {
		row := rowIdx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Id)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.ProductType)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), item.Taste)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), item.Volume)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), item.Count)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
