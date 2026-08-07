package repository

import (
	"bytes"
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"strconv"
)

type MarketingPostgres struct {
	db *gorm.DB
}

func NewMarketingPostgres(db *gorm.DB) *MarketingPostgres {
	return &MarketingPostgres{db: db}
}

func (r *MarketingPostgres) GetMarketingMainPage(cityId, districtId, tradePointId, agentId, brand int, volume float64, taste, language string) ([]models.MarketingMainPage, error) {
	var result []models.MarketingMainPage
	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}). //Joins("JOIN product_types ON product_types.id = products.product_type_id").
								Preload("ProductType").Preload("Name").Order("id").Where("id > 0")

	if volume > 0 {
		productsRequest = productsRequest.Where("product_type_id = ?", volume)
	}

	if len(taste) > 0 {
		productsRequest = productsRequest.Where("taste_ru ILIKE ? ", "%"+taste+"%").Or("taste_tm ILIKE ?", "%"+taste+"%")
	}

	if brand > 0 {
		productsRequest = productsRequest.Where("brand_id = ?", brand)
	}

	productsRequest.Find(&products)

	if productsRequest.Error != nil {
		fmt.Println("Productlar tapylmady")
	}

	for _, element := range products {
		var data models.MarketingMainPage
		data.ProductId = element.Id
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

	var trade_point_count int64

	get_trade_point_count_request := r.db.Model(&models.TradePoint{}).Where("id >0")

	if cityId > 0 {
		get_trade_point_count_request = get_trade_point_count_request.Where("city_id = ?", cityId)
	}

	if districtId > 0 {
		get_trade_point_count_request = get_trade_point_count_request.Where("district_id = ?", districtId)
	}

	if tradePointId > 0 {
		get_trade_point_count_request = get_trade_point_count_request.Where("id = ?", tradePointId)
	}

	get_trade_point_count_request.Count(&trade_point_count)

	for idx, element := range result {
		var list []models.Marketing

		getRequest := r.db.Model(&models.Marketing{}).Joins("JOIN trade_points ON trade_points.id = marketings.trade_point_id").
			Preload("TradePoint").Where("product_id = ?", element.ProductId)

		if cityId > 0 {
			getRequest = getRequest.Where("trade_points.city_id = ?", cityId)
		}

		if districtId > 0 {
			getRequest = getRequest.Where("trade_points.district_id = ?", districtId)
		}

		if tradePointId > 0 {
			getRequest = getRequest.Where("trade_points.id = ?", tradePointId)
		}

		if agentId > 0 {
			getRequest = getRequest.Where("trade_points.trade_agent_id = ?", agentId)
		}

		var count int64
		getRequest.Count(&count)
		getRequest = getRequest.Find(&list)

		distribution := 0
		have := 0
		wystawlaymost := 0
		rasprodano := 0
		for _, item := range list {
			if item.Saled {
				rasprodano++
			}

			if item.Exposure {
				wystawlaymost++
			}

			if item.Have {
				have++
			}

			if item.Have || item.Exposure {
				distribution++
			}
		}
		if distribution > 0 {
			result[idx].Distribution = int(float64(distribution*100) / float64(count))
		} else {
			result[idx].Distribution = 0
		}

		if have > 0 {
			result[idx].Have = int(float64(have*100) / float64(count))
		}

		if wystawlaymost > 0 {
			result[idx].Wystawlaymost = int(float64(wystawlaymost*100) / float64(count))
		} else {
			result[idx].Wystawlaymost = 0
		}

		if rasprodano > 0 {
			result[idx].Rasprodano = int(float64(rasprodano*100) / float64(count))
		} else {
			result[idx].Rasprodano = 0
		}

		result[idx].Trade_point_Count = int(trade_point_count)
	}

	if result == nil {
		result = []models.MarketingMainPage{}
	}

	return result, nil
}

func (r *MarketingPostgres) MarketingMainPageExport(cityId, districtId, tradePointId, agentId, brand int, volume float64, taste, language string) ([]byte, error) {

	list, err := r.GetMarketingMainPage(cityId, districtId, tradePointId, agentId, brand, volume, taste, language)

	if err != nil {
		return []byte{}, err
	}

	f := excelize.NewFile()
	sheet := "Filtrs"
	f.NewSheet(sheet)

	headers := []string{"Бренд", "Тип продукции", "Объём (в литрах)", "Вкус", "Дистрибуция", "Выставленность", "Распродано"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, item := range list {
		row := rowIdx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.ProductType)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.Volume)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), item.Taste)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), item.Distribution)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), item.Wystawlaymost)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), item.Rasprodano)
	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
func (r *MarketingPostgres) GetDistribution(tradePointId int, language string) ([]models.ReadMarketingById, error) {
	if tradePointId == 0 {
		return []models.ReadMarketingById{}, nil
	}

	var result []models.ReadMarketingById
	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Order("id").Where("id > 0").Find(&products)
	if productsRequest.Error != nil {
		fmt.Println("Productlar tapylmady")
	}

	var list []models.Marketing

	getRequest := r.db.Model(&models.Marketing{}).Joins("JOIN trade_points ON trade_points.id = marketings.trade_point_id").
		Preload("TradePoint").Where("trade_point_id = ?", tradePointId).Find(&list)

	if getRequest.Error != nil {
		fmt.Println(getRequest.Error)
	}

	for _, element := range products {
		var data models.ReadMarketingById
		data.ProductId = element.Id
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

	for _, list_element := range list {
		for idx, result_element := range result {
			if list_element.ProductId == result_element.ProductId {
				result[idx].Status = list_element.Have
			}
		}
	}

	if result == nil {
		result = []models.ReadMarketingById{}
	}

	return result, nil
}

func (r *MarketingPostgres) DistributionExport(tradePointId int, language string) ([]byte, error) {
	list, err := r.GetDistribution(tradePointId, language)
	if err != nil {
		return []byte{}, err
	}

	f := excelize.NewFile()
	sheet := "Filtrs"

	f.NewSheet(sheet)

	headers := []string{"Бренд", "Тип продукции", "Объём (в литрах)", "Вкус", "Статус"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, item := range list {
		row := rowIdx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.ProductType)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.Volume)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), item.Taste)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), item.Status)
	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *MarketingPostgres) GetRasprodano(tradePointId int, language string) ([]models.ReadMarketingById, error) {
	if tradePointId == 0 {
		return []models.ReadMarketingById{}, nil
	}

	var result []models.ReadMarketingById
	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Order("id").Where("id > 0").Find(&products)
	if productsRequest.Error != nil {
		fmt.Println("Productlar tapylmady")
	}

	var list []models.Marketing

	getRequest := r.db.Model(&models.Marketing{}).Joins("JOIN trade_points ON trade_points.id = marketings.trade_point_id").
		Preload("TradePoint").Where("trade_point_id = ?", tradePointId).Find(&list)

	if getRequest.Error != nil {
		fmt.Println(getRequest.Error)
	}

	for _, element := range products {
		var data models.ReadMarketingById
		data.ProductId = element.Id
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

	for _, list_element := range list {
		for idx, result_element := range result {
			if list_element.ProductId == result_element.ProductId {
				result[idx].Status = list_element.Saled
			}
		}
	}

	if result == nil {
		result = []models.ReadMarketingById{}
	}

	return result, nil
}

func (r *MarketingPostgres) RasprodanoExport(tradePointId int, language string) ([]byte, error) {
	list, err := r.GetRasprodano(tradePointId, language)
	if err != nil {
		return []byte{}, err
	}

	f := excelize.NewFile()
	sheet := "Filtrs"
	f.NewSheet(sheet)

	headers := []string{"Бренд", "Тип продукции", "Объём (в литрах)", "Вкус", "Статус"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, item := range list {
		row := rowIdx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.ProductType)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.Volume)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), item.Taste)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), item.Status)
	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *MarketingPostgres) GetWystawlaymost(tradePointId int, language string) ([]models.ReadMarketingById, error) {
	if tradePointId == 0 {
		return []models.ReadMarketingById{}, nil
	}

	var result []models.ReadMarketingById
	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Order("id").Where("id > 0").Find(&products)
	if productsRequest.Error != nil {
		fmt.Println("Productlar tapylmady")
	}

	var list []models.Marketing

	getRequest := r.db.Model(&models.Marketing{}).Joins("JOIN trade_points ON trade_points.id = marketings.trade_point_id").
		Preload("TradePoint").Where("trade_point_id = ?", tradePointId).Find(&list)

	if getRequest.Error != nil {
		fmt.Println(getRequest.Error)
	}

	for _, element := range products {
		var data models.ReadMarketingById
		data.ProductId = element.Id
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

	for _, list_element := range list {
		for idx, result_element := range result {
			if list_element.ProductId == result_element.ProductId {
				result[idx].Status = list_element.Exposure
			}
		}
	}

	if result == nil {
		result = []models.ReadMarketingById{}
	}

	return result, nil
}

func (r *MarketingPostgres) WystawlaymostExport(tradePointId int, language string) ([]byte, error) {
	list, err := r.GetWystawlaymost(tradePointId, language)
	if err != nil {
		return []byte{}, err
	}

	f := excelize.NewFile()
	sheet := "Filtrs"
	f.NewSheet(sheet)

	headers := []string{"Бренд", "Тип продукции", "Объём (в литрах)", "Вкус", "Статус"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, item := range list {
		row := rowIdx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.ProductType)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.Volume)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), item.Taste)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), item.Status)
	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
