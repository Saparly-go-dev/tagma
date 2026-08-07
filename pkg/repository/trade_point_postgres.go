package repository

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/nfnt/resize"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type TradePointPostgres struct {
	db *gorm.DB
}

func NewTradePointPostgres(db *gorm.DB) *TradePointPostgres {
	return &TradePointPostgres{db: db}
}

func (r *TradePointPostgres) Create(point models.CreateTradePoint, image1, image2 *multipart.FileHeader) (int, error) {
	//fmt.Println("geldi")
	var lastdata models.TradePoint
	var city tagma.City
	var code string
	var sany int
	//fmt.Println("geldi1234")

	cityresult := r.db.Model(tagma.City{}).Where("id = ?", point.CityId).Find(&city)
	lastresult := r.db.Model(&models.TradePoint{}).
		Where("city_id = ?", point.CityId).
		Order("id DESC").
		Last(&lastdata)
	//fmt.Println("geld12341234i")

	if cityresult.Error != nil {
		return 0, cityresult.Error
	}

	//fmt.Printf(lastdata.Code)
	//fmt.Printf(strconv.Itoa(lastdata.CityId))

	sany, err := strconv.Atoi(lastdata.Code)
	if err != nil || lastresult.Error != nil {
		code = city.Code + "0001"
	} else {
		code = strconv.Itoa(sany + 1)
	}
	//fmt.Println(code)

	if len(code) < 6 {
		code = "0" + code
	}

	newPoint := models.TradePoint{
		Code:            code,
		NameRu:          point.NameRu,
		NameTm:          point.NameTm,
		LocationRu:      point.LocationRu,
		LocationTm:      point.LocationTm,
		OrientirRu:      point.OrientirRu,
		OrientirTm:      point.OrientirTm,
		CityId:          point.CityId,
		DistrictId:      point.DistrictId,
		TradeAgentId:    point.TradeAgentId,
		TradeChannelId:  point.TradeChannelId,
		DayId:           point.DayId,
		EkspeditorId:    point.EkspeditorId,
		Status:          point.Status,
		TradeCategoryId: 1,
		UpdatedAt:       time.Now(),
	}

	result := r.db.Create(&newPoint)
	//fmt.Println(result)
	if result.Error != nil {
		return 0, result.Error
	}

	var imageNames []string

	if (image1 != nil) && (image2 != nil) {
		imageNames, _ = SaveImages(newPoint.Id, *image1, *image2)
	}

	if len(imageNames) == 0 {
		// surat yok bolsa onda hic hili hereket etmeyaris
	} else {
		for _, element := range imageNames {
			var imagedata = models.Image{
				Name:         element,
				TradePointId: newPoint.Id,
			}

			image1result := r.db.Model(&models.Image{}).Create(&imagedata)
			if image1result.Error != nil {
				log.Printf("Failed to save image1: %v", image1result.Error)
			}
		}

	}

	return newPoint.Id, nil
}

func (r *TradePointPostgres) ReSaveImage(PointId int, file *multipart.FileHeader) error {

	var newimage = multipart.FileHeader{}

	fileExt := filepath.Ext(file.Filename)
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	if !contains(allowedExtensions, fileExt) {
		return errors.New("Invalid file type. Please upload an image file.")
	}

	var imageNames, _ = SaveImages(PointId, *file, newimage)

	if len(imageNames) == 0 {
		return errors.New("No images uploaded")
	}

	for _, element := range imageNames {

		if len(element) == 0 {
			continue
		}

		var imagedata = models.Image{
			Name:         element,
			TradePointId: PointId,
		}

		image1result := r.db.Model(&models.Image{}).Create(&imagedata)
		if image1result.Error != nil {
			log.Printf("Failed to save image1: %v", image1result.Error)
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (r *TradePointPostgres) GetExcel(cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId int, name, language string) ([]byte, error) {
	var rawdata []models.TradePoint
	var pageresult []models.ReadTradePoint

	result := r.db.Model(&models.TradePoint{}).
		Joins("JOIN trade_channels ON trade_channels.id = trade_points.trade_channel_id").
		Preload("City").Preload("District").Preload("TradeAgent").
		Preload("TradeChannel").Preload("TradeCategory").Preload("Day").Preload("Ekspeditor")

	if len(name) > 0 {
		result = result.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%")
	}

	if cityId > 0 {
		result = result.Where("city_id = ?", cityId)
	}

	if agentId > 0 {
		result = result.Where("trade_agent_id = ?", agentId)
	}

	if districtId > 0 {
		//fmt.Println("distictId = ", districtId)
		result = result.Where("district_id = ?", districtId)
	}

	if ekspeditorId > 0 {
		result = result.Where("ekspeditor_id = ?", ekspeditorId)
	}

	if dayId > 0 {
		result = result.Where("day_id = ?", dayId)
	}

	if typeId > 0 {
		result = result.Where("trade_channels.type_id = ?", typeId)
	}

	if structureId > 0 {
		result = result.Where("trade_channels.structure_id = ?", structureId)
	}

	if sizeId > 0 {
		result = result.Where("trade_channels.size_id = ?", sizeId)
	}

	if managementId > 0 {
		result = result.Where("trade_channels.management_id = ?", managementId)
	}

	var count int64

	result = result.Count(&count)

	result = result.Order("id desc").Find(&rawdata)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range rawdata {
		var newItem models.ReadTradePoint
		var tradechannel models.TradeChannel

		channelRequest := r.db.Model(&models.TradeChannel{}).Preload("Type").Preload("Structure").Preload("Size").
			Preload("Management").Where("id = ?", element.TradeChannelId).First(&tradechannel)

		if channelRequest.Error != nil {
			//fmt.Println(channelRequest.Error)
			//continue
		}

		newItem.Id = element.Id
		newItem.Code = element.Code
		newItem.CityId = element.CityId
		newItem.DistrictId = element.DistrictId
		newItem.TradeAgentId = element.TradeAgentId
		newItem.TradeChannelId = element.TradeChannelId
		newItem.TradeCategoryId = element.TradeCategoryId
		newItem.DayId = element.DayId
		newItem.EkspeditorId = element.EkspeditorId
		newItem.Status = element.Status
		newItem.UpdatedAt = element.UpdatedAt.Format("02.01.2006")
		if language == "ru" {
			newItem.Name = element.NameRu
			newItem.Location = element.LocationRu
			newItem.Orientir = element.OrientirRu
			newItem.City = element.City.NameRu
			newItem.District = element.District.NameRu
			newItem.TradeAgent = element.TradeAgent.NameRu
			newItem.TradeChannel = tradechannel.Type.NameRu + "-" + tradechannel.Structure.NameRu + "-" + tradechannel.Size.NameRu + "-" + tradechannel.Management.NameRu
			newItem.TradCategory = element.TradeCategory.NameRu
			newItem.Day = element.Day.NameRu
			newItem.Ekspeditor = element.Ekspeditor.NameRu
		} else {
			newItem.Name = element.NameTm
			newItem.Location = element.LocationTm
			newItem.Orientir = element.OrientirTm
			newItem.City = element.City.NameTm
			newItem.District = element.District.NameTm
			newItem.TradeAgent = element.TradeAgent.NameTm
			newItem.TradeChannel = tradechannel.Type.NameTm + "-" + tradechannel.Structure.NameTm + "-" + tradechannel.Size.NameTm + "-" + tradechannel.Management.NameTm
			newItem.TradCategory = element.TradeCategory.NameTm
			newItem.Day = element.Day.NameTm
			newItem.Ekspeditor = element.Ekspeditor.NameTm
		}

		pageresult = append(pageresult, newItem)
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheet := "Торговые точки"
	f.NewSheet(sheet)

	// Define headers
	headers := []string{"ID", "Код торговой точки", "Наименование", "Местонахождение", "Ориентир", "Город", "Район", "Торговый представитель", "Тип канала торговли", "Категория продаж", "Статус", "День недели", "Экспедитор"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	// Populate data rows
	for rowIdx, item := range pageresult {
		row := rowIdx + 2 // Start from the second row
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Id)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.Code)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), item.Name)
		f.SetCellValue(sheet, "D"+strconv.Itoa(row), item.Location)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), item.Orientir)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), item.City)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), item.District)
		f.SetCellValue(sheet, "H"+strconv.Itoa(row), item.TradeAgent)
		f.SetCellValue(sheet, "I"+strconv.Itoa(row), item.TradeChannel)
		f.SetCellValue(sheet, "J"+strconv.Itoa(row), item.TradCategory)
		f.SetCellValue(sheet, "K"+strconv.Itoa(row), item.Status)
		f.SetCellValue(sheet, "L"+strconv.Itoa(row), item.Day)
		f.SetCellValue(sheet, "M"+strconv.Itoa(row), item.Ekspeditor)

	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *TradePointPostgres) GetAll(pageSize, pageNumber, cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, status int, name, code, contact, language string) (*models.TradePointPage, error) {
	var rawdata []models.TradePoint
	var pageresult []models.ReadTradePoint
	var db_furniture_list []models.Furniture
	var page_furniture_list []models.ReadFurnitureForMobile

	get_all_furniture_request := r.db.Model(&models.Furniture{}).Where("id >0").Find(&db_furniture_list)

	if get_all_furniture_request.Error != nil {

	}

	for _, furniture := range db_furniture_list {
		raw_data := models.ReadFurnitureForMobile{}
		raw_data.Count = 0
		if language == "ru" {
			raw_data.Name = furniture.NameRu
		} else {
			raw_data.Name = furniture.NameTm
		}

		page_furniture_list = append(page_furniture_list, raw_data)
	}

	offset := (pageNumber - 1) * pageSize
	fmt.Println(cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId)
	fmt.Println(name, "-", code, "-", contact)

	result := r.db.Model(&models.TradePoint{}).
		Joins("JOIN trade_channels ON trade_channels.id = trade_points.trade_channel_id").
		Preload("City").Preload("District").Preload("TradeAgent").
		Preload("TradeChannel").Preload("TradeCategory").Preload("Day").Preload("Ekspeditor")

	if len(name) > 0 {
		result = result.Where("name_ru ILIKE ? OR name_tm = ?", "%"+name+"%", "%"+name+"%")
	}

	if len(code) > 0 {
		result = result.Where("code ILIKE ?", "%"+code+"%")
	}

	if cityId > 0 {
		result = result.Where("city_id = ?", cityId)
	}

	if agentId > 0 {
		result = result.Where("trade_agent_id = ?", agentId)
	}

	if districtId > 0 {
		//fmt.Println("distictId = ", districtId)

		result = result.Where("district_id = ?", districtId)
	}

	if ekspeditorId > 0 {
		result = result.Where("ekspeditor_id = ?", ekspeditorId)
	}

	if dayId > 0 {
		result = result.Where("day_id = ?", dayId)
	}

	if status == 1 {
		result = result.Where("trade_points.status = ?", true)
	}

	if status == 0 {
		result = result.Where("trade_points.status = ?", false)
	}

	if typeId > 0 {
		result = result.Where("trade_channels.type_id = ?", typeId)
	}

	if structureId > 0 {
		result = result.Where("trade_channels.structure_id = ?", structureId)
	}

	if sizeId > 0 {
		result = result.Where("trade_channels.size_id = ?", sizeId)
	}

	if managementId > 0 {
		result = result.Where("trade_channels.management_id = ?", managementId)
	}

	if len(contact) > 0 {
		result = result.Where("EXISTS (SELECT 1 FROM contacts WHERE contacts.trade_point_id = trade_points.id AND (contacts.name_ru ILIKE ? OR contacts.name_tm ILIKE ? OR contacts.number ILIKE ?))", "%"+contact+"%", "%"+contact+"%", "%"+contact+"%")
	}

	var count int64

	result = result.Count(&count)

	result = result.Order("id desc").Offset(offset).Limit(pageSize).Find(&rawdata)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range rawdata {
		var newItem models.ReadTradePoint
		var tradechannel models.TradeChannel

		var furniture_list []models.TradePointsFurniture

		get_furnitures_request := r.db.Model(&models.TradePointsFurniture{}).Preload("Furniture").
			Where("trade_point_id = ?", element.Id).Find(&furniture_list)

		if get_furnitures_request.Error != nil {

		}

		for _, furniture := range furniture_list {
			var raw_furniture_data models.ReadFurnitureForMobile
			raw_furniture_data.Count = furniture.Count

			if language == "ru" {
				raw_furniture_data.Name = furniture.Furniture.NameRu
			} else {
				raw_furniture_data.Name = furniture.Furniture.NameTm
			}

			newItem.FurnitureList = append(newItem.FurnitureList, raw_furniture_data)

		}

		channelRequest := r.db.Model(&models.TradeChannel{}).Preload("Type").Preload("Structure").Preload("Size").
			Preload("Management").Where("id = ?", element.TradeChannelId).First(&tradechannel)

		if channelRequest.Error != nil {
			//fmt.Println(channelRequest.Error)
			//continue
		}

		trade_point_debt_information := r.GetTradePoinstDebt(element.Id)

		newItem.Id = element.Id
		newItem.Code = element.Code
		newItem.CityId = element.CityId
		newItem.DistrictId = element.DistrictId
		newItem.TradeAgentId = element.TradeAgentId
		newItem.TradeChannelId = element.TradeChannelId
		newItem.TradeCategoryId = element.TradeCategoryId
		newItem.DayId = element.DayId
		newItem.EkspeditorId = element.EkspeditorId
		newItem.Status = element.Status
		newItem.UpdatedAt = element.UpdatedAt.Format("02.01.2006")
		newItem.Debt = trade_point_debt_information.Debt
		if language == "ru" {
			newItem.Name = element.NameRu
			newItem.Location = element.LocationRu
			newItem.Orientir = element.OrientirRu
			newItem.City = element.City.NameRu
			newItem.District = element.District.NameRu
			newItem.TradeAgent = element.TradeAgent.NameRu
			newItem.TradeChannel = tradechannel.Type.NameRu + "-" + tradechannel.Structure.NameRu + "-" + tradechannel.Size.NameRu + "-" + tradechannel.Management.NameRu
			newItem.TradCategory = element.TradeCategory.NameRu
			newItem.Day = element.Day.NameRu
			newItem.Ekspeditor = element.Ekspeditor.NameRu
		} else {
			newItem.Name = element.NameTm
			newItem.Location = element.LocationTm
			newItem.Orientir = element.OrientirTm
			newItem.City = element.City.NameTm
			newItem.District = element.District.NameTm
			newItem.TradeAgent = element.TradeAgent.NameTm
			newItem.TradeChannel = tradechannel.Type.NameTm + "-" + tradechannel.Structure.NameTm + "-" + tradechannel.Size.NameTm + "-" + tradechannel.Management.NameTm
			newItem.TradCategory = element.TradeCategory.NameTm
			newItem.Day = element.Day.NameTm
			newItem.Ekspeditor = element.Ekspeditor.NameTm
		}

		pageresult = append(pageresult, newItem)

		for i, p_furniture := range page_furniture_list {
			for _, item_furniture := range newItem.FurnitureList {
				if p_furniture.Name == item_furniture.Name {
					page_furniture_list[i].Count += item_furniture.Count
				}
			}
		}
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if pageresult == nil {
		pageresult = []models.ReadTradePoint{}
	}

	// Create the page object
	page := models.TradePointPage{
		Items:         pageresult,
		PageCount:     pageCount,
		PageNumber:    pageNumber,
		PageSize:      pageSize,
		Total:         int(count),
		FurnitureList: page_furniture_list,
	}

	return &page, nil
}

func (r *TradePointPostgres) ChangeStatus(Id int) error {
	var data models.TradePoint
	result := r.db.Model(&models.TradePoint{}).First(&data, Id)

	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.TradePoint{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.TradePoint{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *TradePointPostgres) GetGeneralInformation(Id int, language string) (*models.GeneralTradePoint, error) {
	var data models.TradePoint
	var resultData models.GeneralTradePoint

	result := r.db.Model(&models.TradePoint{}).Preload("City").Preload("TradeAgent").Preload("TradeChannel").Preload("TradeCategory").Preload("District").
		Preload("Day").Preload("Ekspeditor").Where("id = ?", Id).First(&data)

	if result.Error != nil {
		return nil, result.Error
	}
	var tradechannel models.TradeChannel
	channelRequest := r.db.Model(&models.TradeChannel{}).Preload("Type").Preload("Structure").Preload("Size").
		Preload("Management").Where("id = ?", data.TradeChannelId).First(&tradechannel)

	trade_point_debts := r.GetTradePoinstDebt(Id)

	if channelRequest.Error != nil {
		fmt.Println(channelRequest.Error)
	}
	resultData.Id = data.Id
	resultData.Code = data.Code
	resultData.CityId = data.CityId
	resultData.DistrictId = data.DistrictId
	resultData.DayId = data.DayId
	resultData.EkspeditorId = data.EkspeditorId
	resultData.Debt = trade_point_debts.Debt
	resultData.Cash = trade_point_debts.Cash
	resultData.Terminal = trade_point_debts.Terminal
	resultData.Transfer = trade_point_debts.Transfer

	if language == "ru" {
		resultData.Name = data.NameRu
		resultData.City = data.City.NameRu
		resultData.District = data.District.NameRu
		resultData.Location = data.LocationRu
		resultData.Orientir = data.OrientirRu
		resultData.TradeChannel = tradechannel.Type.NameRu + "-" + tradechannel.Structure.NameRu + "-" + tradechannel.Size.NameRu + "-" + tradechannel.Management.NameRu
		resultData.Day = data.Day.NameRu
		resultData.Ekspeditor = data.Ekspeditor.NameRu
		resultData.TradeAgent = data.TradeAgent.NameRu
		resultData.TradeCategory = data.TradeCategory.NameRu
	} else {
		resultData.Name = data.NameTm
		resultData.City = data.City.NameTm
		resultData.District = data.District.NameTm
		resultData.Location = data.LocationTm
		resultData.Orientir = data.OrientirTm
		resultData.TradeChannel = tradechannel.Type.NameTm + "-" + tradechannel.Structure.NameTm + "-" + tradechannel.Size.NameTm + "-" + tradechannel.Management.NameTm
		resultData.Day = data.Day.NameTm
		resultData.Ekspeditor = data.Ekspeditor.NameTm
		resultData.TradeAgent = data.TradeAgent.NameTm
		resultData.TradeCategory = data.TradeCategory.NameTm
	}

	var images []models.Image

	resultImage := r.db.Model(models.Image{}).Where("trade_point_id = ?", Id).Find(&images)

	if result.Error != nil {
		return nil, resultImage.Error
	}

	for _, image := range images {
		var item models.ReadImage

		item.Id = image.Id
		item.Path = "images/" + strconv.Itoa(image.TradePointId) + "/" + image.Name

		resultData.Images = append(resultData.Images, item)
	}

	return &resultData, nil
}

func (r *TradePointPostgres) GetPointImages(Id int) (*[]models.ReadImage, error) {
	var images []models.Image
	var readImage []models.ReadImage

	resultImage := r.db.Model(models.Image{}).Where("trade_point_id = ?", Id).Find(&images)

	if resultImage.Error != nil {
		return nil, resultImage.Error
	}

	for _, image := range images {
		var item models.ReadImage

		item.Id = image.Id
		item.Path = "images/" + strconv.Itoa(image.TradePointId) + "/" + image.Name

		readImage = append(readImage, item)
	}

	return &readImage, nil
}

func (r *TradePointPostgres) DeleteImage(Id int) error {
	var image models.Image

	result := r.db.Model(&models.Image{}).Where("id = ?", Id).Find(&image)

	if result.Error != nil {
		return result.Error
	}

	outputDir := filepath.Join("images", strconv.Itoa(image.TradePointId))
	outputDir = outputDir + "/" + image.Name
	//fmt.Println(outputDir)
	err := deleteFile(outputDir)

	if err != nil {
		return err
	}

	resultDelete := r.db.Model(&models.Image{}).Where("id = ?", Id).Delete(&image)

	if resultDelete.Error != nil {
		return resultDelete.Error
	}

	return nil
}

func (r *TradePointPostgres) GetById(Id int) (models.TradePoint, error) {
	var data models.TradePoint

	result := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").Preload("TradeAgent").Preload("TradeChannel").Where("id = ?", Id).First(&data)

	return data, result.Error
}

func (r *TradePointPostgres) Delete(Id int) error {
	var images []models.Image
	var path string

	imageResult := r.db.Model(&models.Image{}).Where("trade_point_id = ?", Id).Find(&images)

	if imageResult.Error != nil {
		return imageResult.Error
	}

	path = "images/" + strconv.Itoa(Id)

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	imageDelete := r.db.Model(models.Image{}).Delete(&images)

	if imageDelete.Error != nil {
		return imageDelete.Error
	}

	result := r.db.Model(&models.TradePoint{}).Delete(&models.TradePoint{}, Id)

	return result.Error
}

func (r *TradePointPostgres) Update(Id int, point models.CreateTradePoint) error {
	updateData := map[string]interface{}{
		"name_ru":          point.NameRu,
		"name_tm":          point.NameTm,
		"location_ru":      point.LocationRu,
		"location_tm":      point.LocationTm,
		"orientir_ru":      point.OrientirRu,
		"orientir_tm":      point.OrientirTm,
		"city_id":          point.CityId,
		"district_id":      point.DistrictId,
		"trade_agent_id":   point.TradeAgentId,
		"trade_channel_id": point.TradeChannelId,
		"day_id":           point.DayId,
		"ekspeditor_id":    point.EkspeditorId,
	}

	result := r.db.Model(&models.TradePoint{}).Where("id = ?", Id).Updates(updateData)

	return result.Error
}

func (r *TradePointPostgres) ChangeTradePointAgent(agent_id, ekspeditor_id, day_id int, ids []int) error {

	for _, id := range ids {
		var count int64
		is_have_request := r.db.Model(&models.TradePoint{}).Where("id = ?", id).Count(&count)
		if is_have_request.Error != nil || count == 0 {
			continue
		}

		updateData := map[string]interface{}{
			"trade_agent_id": agent_id,
			"ekspeditor_id":  ekspeditor_id,
			"day_id":         day_id,
		}

		update_request := r.db.Model(&models.TradePoint{}).Where("id = ?", id).Updates(updateData)
		if update_request.Error != nil {
		}
	}

	return nil
}

func (r *TradePointPostgres) GetCities(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []tagma.City

	result := r.db.Model(&tagma.City{}).Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range listOfData {
		var item models.SelectObject

		item.Id = element.Id
		if language == "ru" {
			item.Name = element.NameRu
		} else {
			item.Name = element.NameTm
		}

		responsedata = append(responsedata, item)
	}

	return &responsedata, nil
}

func (r *TradePointPostgres) GetDistrictsByCity(cityId int, language string) (*[]models.SelectObject, error) {
	var districts []models.District
	var result []models.SelectObject

	request := r.db.Model(&models.District{}).Where("city_id = ?", cityId).Find(&districts)
	if request.Error != nil {
		return nil, request.Error
	}

	for _, element := range districts {
		var item models.SelectObject
		item.Id = element.Id
		if language == "ru" {
			item.Name = element.NameRu
		} else {
			item.Name = element.NameTm
		}
		result = append(result, item)
	}
	if result == nil {
		result = []models.SelectObject{}
	}

	return &result, nil
}

func (r *TradePointPostgres) GetAgents(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.TradeAgent

	result := r.db.Model(&models.TradeAgent{}).Order("id").Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range listOfData {
		var item models.SelectObject

		item.Id = element.Id
		if language == "ru" {
			item.Name = element.NameRu
		} else {
			item.Name = element.NameTm
		}

		responsedata = append(responsedata, item)
	}
	return &responsedata, nil
}

func (r *TradePointPostgres) GetChannels(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.TradeChannel

	result := r.db.Model(&models.TradeChannel{}).Preload("Type").Preload("Structure").Preload("Size").
		Preload("Management").Find(&listOfData)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, tradechannel := range listOfData {
		var item models.SelectObject

		item.Id = tradechannel.Id
		if language == "ru" {
			item.Name = tradechannel.Type.NameRu + "-" + tradechannel.Structure.NameRu + "-" + tradechannel.Size.NameRu + "-" + tradechannel.Management.NameRu
		} else {
			item.Name = tradechannel.Type.NameTm + "-" + tradechannel.Structure.NameTm + "-" + tradechannel.Size.NameTm + "-" + tradechannel.Management.NameTm

		}

		responsedata = append(responsedata, item)
	}
	return &responsedata, nil
}

func (r *TradePointPostgres) GetInformationAboutOrder(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) ([]models.ReadOrder, error) {

	var data []models.Order
	var list []models.ReadOrder
	var start time.Time
	var end time.Time

	if endDay == 0 || endMonth == 0 || endYear == 0 {
		end = time.Now()
	} else {
		end = time.Date(endYear, time.Month(endMonth), endDay, 0, 0, 0, 0, time.UTC)
	}

	if startDay == 0 || startMonth == 0 || startYear == 0 {
		start = end.Add(time.Hour * 24 * (-7))
	} else {
		start = time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)

	}

	request := r.db.Model(&models.Order{}).Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("TradePoint")

	var tradePoint models.TradePoint

	pointReqeust := r.db.Model(&models.TradePoint{}).Preload("TradeAgent").
		Preload("City").Preload("District").Where("id = ?", tradePointId).Find(&tradePoint)

	if pointReqeust.Error != nil {
	}

	var count int64

	request = request.Count(&count)

	request = request.Order("id desc").
		Where("trade_point_id = ?", tradePointId).
		Where(`"created_at" <= ?`, end).
		Where(`"created_at" >= ?`, start).Find(&data)

	for _, element := range data {
		var paid float64
		var payments []models.Payment

		get_payments_request := r.db.Model(&models.Payment{}).Where("order_id = ?", element.Id).Find(&payments)

		if get_payments_request.Error != nil {
			//fmt.Println(get_payments_request.Error)
		}

		for _, payment_element := range payments {
			paid += payment_element.Currency
		}

		not_paid := element.Sum - paid
		paid = Round2(paid)
		not_paid = Round2(not_paid)

		var item models.ReadOrder
		item.Id = element.Id
		item.TradePointId = element.TradePointId
		item.TradeAgentId = element.TradePoint.TradeAgentId
		item.Sum = element.Sum
		item.NotPaid = not_paid
		item.Paid = paid
		item.CreatedAt = element.CreatedAt.Format("02.01.2006")
		item.Status = element.Status

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

	if list == nil {
		list = []models.ReadOrder{}
	}

	return list, nil
}

func (r *TradePointPostgres) GetInformantionAboutProducts(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) ([]models.ReadProductForTradePoint, error) {
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
		start = end.Add(time.Hour * 24 * (-7))
	} else {
		start = time.Date(startYear, time.Month(startMonth), startDay, 0, 0, 0, 0, time.UTC)

	}

	request := r.db.Model(&models.Order{}).Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("List")

	var tradePoint models.TradePoint

	pointReqeust := r.db.Model(&models.TradePoint{}).Preload("TradeAgent").
		Preload("City").Preload("District").Where("id = ?", tradePointId).Find(&tradePoint)

	if pointReqeust.Error != nil {
	}

	request = request.Order("id desc").
		Where("trade_point_id = ?", tradePointId).
		Where(`"created_at" <= ?`, end).
		Where(`"created_at" >= ?`, start).Find(&orders)

	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Order("id").Where("id > 0").Find(&products)
	if productsRequest.Error != nil {
		//fmt.Println("Productlar tapylmady")
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

func (r *TradePointPostgres) GetTradePointSaleHistoryForMonth(tradePointId int, language string) ([]models.SaleHistoryTradePoint, error) {
	var result []models.SaleHistoryTradePoint
	var firstorders, secondOrders, thirdOrders, fourthOrders []models.Order

	now := time.Now()
	firstWeekDay := now.Add(time.Hour * 24 * (-7))
	secondWeekDay := firstWeekDay.Add(time.Hour * 24 * (-7))
	thirdWeekDay := secondWeekDay.Add(time.Hour * 24 * (-7))
	fourthWeekDay := thirdWeekDay.Add(time.Hour * 24 * (-7))

	request := r.db.Model(&models.Order{}).Joins("JOIN trade_points ON trade_points.id = orders.trade_point_id").
		Preload("List")

	var tradePoint models.TradePoint

	pointReqeust := r.db.Model(&models.TradePoint{}).Preload("TradeAgent").
		Preload("City").Preload("District").Where("id = ?", tradePointId).Find(&tradePoint)

	if pointReqeust.Error != nil {
	}

	request = request.Order("id desc").
		Where("trade_point_id = ?", tradePointId).
		Where(`"created_at" <= ?`, now).
		Where(`"created_at" >= ?`, firstWeekDay).Find(&firstorders)

	request = request.Order("id desc").
		Where("trade_point_id = ?", tradePointId).
		Where(`"created_at" <= ?`, firstWeekDay).
		Where(`"created_at" >= ?`, secondWeekDay).Find(&secondOrders)

	request = request.Order("id desc").
		Where("trade_point_id = ?", tradePointId).
		Where(`"created_at" <= ?`, secondWeekDay).
		Where(`"created_at" >= ?`, thirdWeekDay).Find(&thirdOrders)

	request = request.Order("id desc").
		Where("trade_point_id = ?", tradePointId).
		Where(`"created_at" <= ?`, thirdWeekDay).
		Where(`"created_at" >= ?`, fourthWeekDay).Find(&fourthOrders)

	var products []models.Product
	productsRequest := r.db.Model(&models.Product{}).Preload("ProductType").Preload("Name").Order("id").Where("id > 0").Find(&products)
	if productsRequest.Error != nil {
		//fmt.Println("Productlar tapylmady")
	}
	for _, element := range products {
		var data models.SaleHistoryTradePoint
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

	for _, order_element := range firstorders {
		for _, order_list_element := range order_element.List {
			for idx, product_element := range result {
				if order_list_element.ProductId == product_element.Id {
					result[idx].First += order_list_element.Count
				}
			}
		}

	}

	for _, order_element := range secondOrders {
		for _, order_list_element := range order_element.List {
			for idx, product_element := range result {
				if order_list_element.ProductId == product_element.Id {
					result[idx].Second += order_list_element.Count
				}
			}
		}

	}

	for _, order_element := range thirdOrders {
		for _, order_list_element := range order_element.List {
			for idx, product_element := range result {
				if order_list_element.ProductId == product_element.Id {
					result[idx].Third += order_list_element.Count
				}
			}
		}

	}

	for _, order_element := range fourthOrders {
		for _, order_list_element := range order_element.List {
			for idx, product_element := range result {
				if order_list_element.ProductId == product_element.Id {
					result[idx].Fourth += order_list_element.Count
				}
			}
		}

	}

	if result == nil {
		result = []models.SaleHistoryTradePoint{}
	}

	return result, nil
}

func SaveImages(Id int, image1, image2 multipart.FileHeader) ([]string, error) {
	outputDir := filepath.Join("images", strconv.Itoa(Id))
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	//fmt.Printf("Saving images to %s\n", outputDir)

	var firstImageName = RandomString(6) + ".webp"
	if err := processAndSaveImage(image1, filepath.Join(outputDir, firstImageName), "png"); err != nil {
		log.Printf("Failed to process image1: %v", err)
		return []string{}, err
	}

	//fmt.Println(" 12341234")

	var secondImageName = ""
	if image2.Size > 0 {
		secondImageName = RandomString(6) + ".webp"
		if err := processAndSaveImage(image2, filepath.Join(outputDir, secondImageName), "png"); err != nil {
			log.Printf("Failed to process image2: %v", err)
			return []string{}, err
		}
	}

	return []string{firstImageName, secondImageName}, nil
}

func processAndSaveImage(fileHeader multipart.FileHeader, outputPath string, format string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := decodeImage(fileHeader)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	for width > 450 {
		width = width / 2
		height = height / 2
	}

	resizedImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	switch format {
	case "png":
		err = png.Encode(outFile, resizedImg)
	case "jpeg":
		err = jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 80})
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return err
}

func decodeImage(fileHeader multipart.FileHeader) (image.Image, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var img image.Image
	var format string
	switch ext {
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(file)
		format = "jpeg"
	case ".png":
		img, err = png.Decode(file)
		format = "png"
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}

	return img, format, err
}

const harplar = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = harplar[rand.Intn(len(harplar))]
	}
	return string(b)
}

func deleteFile(filePath string) error {
	// Use os.Remove to delete the file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

type TradePointDebtObject struct {
	Debt     float64
	Cash     float64
	Terminal float64
	Transfer float64
}

func (r *TradePointPostgres) GetTradePoinstDebt(Id int) TradePointDebtObject {

	var orders []models.Order
	get_order_request := r.db.Model(&models.Order{}).Where("trade_point_id = ? AND is_closed = ?", Id, false).Find(&orders)

	if get_order_request.Error != nil {
		//fmt.Println(get_order_request.Error)
	}

	var sum float64
	var cash float64
	var terminal float64
	var transfer float64
	var not_paid float64
	var not_paid_cash float64

	for _, order := range orders {
		var payments []models.Payment
		var payment_sum float64
		var this_cash float64
		var this_terminal float64
		var this_transfer float64

		get_payments_request := r.db.Model(&models.Payment{}).Where("order_id = ?", order.Id).Find(&payments)

		if get_payments_request.Error != nil {
			fmt.Println(get_payments_request.Error)
		}

		for _, payment_element := range payments {

			if payment_element.Status == true {
				payment_sum += payment_element.Currency
			} else {
				if payment_element.PaymentTypeId == 1 {
					this_cash += payment_element.Currency
					//fmt.Println("added nalich", payment_element.Currency)
				}

				if payment_element.PaymentTypeId == 2 {
					this_terminal += payment_element.Currency
					//fmt.Println("added terminal", payment_element.Currency)

				}

				if payment_element.PaymentTypeId == 3 {
					this_transfer += payment_element.Currency
					//fmt.Println("added transfer", payment_element.Currency)
				}
			}
		}

		if payment_sum >= order.SaleSum {
			continue
		}

		var tradePoint models.TradePoint

		pointRequest := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").Where("id = ?", order.TradePointId).Find(&tradePoint)
		if pointRequest.Error != nil {
			log.Println(pointRequest.Error)
		}

		var paid float64

		for _, payment_element := range payments {
			if payment_element.Status == true {
				paid += payment_element.Currency
			}
		}

		not_paid = order.SaleSum - paid
		paid = Round2(paid)
		not_paid = Round2(not_paid)

		not_paid_cash = not_paid_cash + (not_paid - this_cash - this_terminal - this_transfer)
		//fmt.Println("not_paid cash", (not_paid - this_cash - this_terminal - this_transfer))
		sum = sum + not_paid
		cash += this_cash
		this_cash = 0
		terminal += this_terminal
		this_terminal = 0
		transfer += this_transfer
		this_transfer = 0
	}

	result := TradePointDebtObject{}
	//fmt.Println("cash", cash)
	//fmt.Println("terminal", terminal)
	//fmt.Println("transfer", transfer)
	//fmt.Println("not_paid", not_paid_cash)

	result.Debt = Round2(cash + terminal + transfer + not_paid_cash)
	result.Cash = Round2(cash + not_paid_cash)
	result.Terminal = Round2(terminal)
	result.Transfer = Round2(transfer)

	return result

}
