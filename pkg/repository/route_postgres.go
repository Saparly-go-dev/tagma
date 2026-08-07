package repository

import (
	"bytes"
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type RoutePostgres struct {
	db *gorm.DB
}

func NewRoutePostgres(db *gorm.DB) *RoutePostgres {
	return &RoutePostgres{db: db}
}

func (r *RoutePostgres) Create(route models.CreateRoute) (int, error) {

	newRoute := models.Route{
		TradeAgentId: route.TradeAgentId,
		TradePointId: route.TradePointId,
		EkspeditorId: route.EkspeditorId,
		DayId:        route.DayId,
		Status:       route.Status,
		UpdatedAt:    time.Now(),
	}

	result := r.db.Create(&newRoute)
	if err := result.Error; err != nil {
		return 0, err
	}

	return newRoute.Id, nil
}

func (r *RoutePostgres) GetById(Id int, language string) (*models.ReadRoute, error) {
	var element models.Route
	result := r.db.Model(models.Route{}).Preload("TradeAgent").Preload("Day").Preload("Ekspeditor").First(&element, Id)

	var tradePoint models.TradePoint
	var item models.ReadRoute
	pointResult := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").Where("id = ?", element.TradePointId).First(&tradePoint)

	if pointResult.Error != nil {

		return nil, pointResult.Error
	}
	item.Id = element.Id
	item.TradePointId = element.TradePointId
	item.TradeAgentId = element.TradeAgentId
	item.DayId = element.DayId
	item.Code = tradePoint.Code
	item.Status = element.Status
	item.EkspeditorId = element.EkspeditorId
	item.Updatedat = element.UpdatedAt.Format("02.01.2006")
	item.CityId = tradePoint.CityId
	item.DistrictId = tradePoint.DistrictId

	if language == "ru" {
		item.TradeAgent = element.TradeAgent.NameRu
		item.Day = element.Day.NameRu
		item.City = tradePoint.City.NameRu
		item.Location = tradePoint.LocationRu
		item.Orientir = tradePoint.OrientirRu
		item.TradePoint = tradePoint.NameRu
		item.Ekspeditor = element.Ekspeditor.NameRu
		item.District = tradePoint.District.NameRu
	} else {
		item.TradeAgent = element.TradeAgent.NameTm
		item.TradePoint = element.TradePoint.NameTm
		item.Day = element.Day.NameTm
		item.City = tradePoint.City.NameTm
		item.Location = tradePoint.LocationTm
		item.Orientir = tradePoint.OrientirTm
		item.TradePoint = tradePoint.NameTm
		item.Ekspeditor = element.Ekspeditor.NameTm
		item.District = tradePoint.District.NameTm
	}

	return &item, result.Error
}

func (r *RoutePostgres) GetAll(pageSize, pageNumber, agentId, cityId, dayId int, language string) (*models.RoutePage, error) {
	var routes []models.Route
	var readRoutes []models.ReadRoute

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.Route{}).
		Joins("JOIN trade_points ON trade_points.id = routes.trade_point_id").
		Joins("JOIN cities ON cities.id = trade_points.city_id").
		Preload("Day").
		Preload("TradeAgent").Preload("Ekspeditor").
		Preload("TradePoint.City")

	if cityId != 0 {
		request = request.Where("cities.id = ?", cityId)
	}

	if agentId != 0 {
		request = request.Where("routes.trade_agent_id = ?", agentId)
	}

	if dayId != 0 {
		request = request.Where("routes.day_id = ?", dayId)
	}

	var count int64

	request = request.Offset(offset).Limit(pageSize).Find(&routes)

	request = request.Count(&count)

	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, element := range routes {
		var tradePoint models.TradePoint
		var item models.ReadRoute

		pointResult := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").Where("id = ?", element.TradePointId).First(&tradePoint)

		if pointResult.Error != nil {
			return nil, pointResult.Error
		}

		item.Id = element.Id
		item.TradePointId = element.TradePointId
		item.TradeAgentId = element.TradeAgentId
		item.DayId = element.DayId
		item.Status = element.Status
		item.Code = tradePoint.Code
		item.Updatedat = element.UpdatedAt.Format("02.01.2006")
		item.EkspeditorId = element.EkspeditorId
		if language == "ru" {
			item.TradeAgent = element.TradeAgent.NameRu
			item.Day = element.Day.NameRu
			item.City = tradePoint.City.NameRu
			item.Location = tradePoint.LocationRu
			item.Orientir = tradePoint.OrientirRu
			item.TradePoint = tradePoint.NameRu
			item.Ekspeditor = element.Ekspeditor.NameRu
			item.District = tradePoint.District.NameRu
		} else {
			item.TradeAgent = element.TradeAgent.NameTm
			item.TradePoint = element.TradePoint.NameTm
			item.Day = element.Day.NameTm
			item.City = tradePoint.City.NameTm
			item.Location = tradePoint.LocationTm
			item.Orientir = tradePoint.OrientirTm
			item.TradePoint = tradePoint.NameTm
			item.Ekspeditor = element.Ekspeditor.NameTm
			item.District = tradePoint.District.NameTm
		}

		readRoutes = append(readRoutes, item)
	}

	if readRoutes == nil {
		readRoutes = []models.ReadRoute{}
	}

	var page = models.RoutePage{
		Items:      readRoutes,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		PageCount:  int(pageCount),
	}

	return &page, nil

}

func (r *RoutePostgres) ChangeStatus(Id int) error {
	var route models.Route
	result := r.db.Model(models.Route{}).First(&route, Id)
	if err := result.Error; err != nil {
		return err
	}

	// Toggle status
	newStatus := !route.Status
	updateData := map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	// Use "=" in the Where clause instead of "=="
	result = r.db.Model(&models.Route{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *RoutePostgres) Delete(Id int) error {
	result := r.db.Delete(&models.Route{}, Id)
	return result.Error
}

func (r *RoutePostgres) Update(Id int, route models.CreateRoute) error {
	updateData := map[string]interface{}{
		"trade_agent_id": route.TradeAgentId,
		"trade_point_id": route.TradePointId,
		"day_id":         route.DayId,
		"status":         route.Status,
		"ekspeditor_id":  route.EkspeditorId,
	}

	result := r.db.Model(models.Route{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *RoutePostgres) GetDays(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.Day

	result := r.db.Model(&models.Day{}).Where("name_tm ILIKE ?", "%"+name+"%").Or("name_ru ILIKE ?", "%"+name+"%").Find(&listOfData)

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

func (r *RoutePostgres) GetPoints(cityId, districtId int, name, language string) (*[]models.PointSelect, error) {
	var responsedata []models.PointSelect
	var listOfData []models.TradePoint

	result := r.db.Model(&models.TradePoint{}).Preload("TradeAgent").Preload("City")

	if cityId > 0 {
		result = result.Where("city_id = ?", cityId)
	}

	if districtId > 0 {
		result = result.Where("district_id = ?", districtId)

	}

	if len(name) > 0 {
		result = result.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%")
	}

	result = result.Find(&listOfData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range listOfData {
		var item models.PointSelect

		item.Id = element.Id
		item.Code = element.Code
		item.TradeAgentId = element.TradeAgent.Id
		if language == "ru" {
			item.Name = element.NameRu
			item.City = element.City.NameRu
			item.Location = element.LocationRu
			item.TradeAgent = element.TradeAgent.NameRu
		} else {
			item.Name = element.NameTm
			item.City = element.City.NameTm
			item.Location = element.LocationTm
			item.TradeAgent = element.TradeAgent.NameTm
		}

		responsedata = append(responsedata, item)
	}

	return &responsedata, nil
}

func (r *RoutePostgres) GetEkspeditors(TradeAgentId int, language string) (*[]models.SelectObject, error) {
	var agent models.TradeAgent
	var result []models.SelectObject

	request := r.db.Model(&models.TradeAgent{}).Preload("Ekspeditor").Where("id = ?", TradeAgentId).Find(&agent)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, item := range agent.Ekspeditor {
		var data models.SelectObject
		data.Id = item.Id
		if language == "ru" {
			data.Name = item.NameRu
		} else {
			data.Name = item.NameTm
		}

		result = append(result, data)
	}

	if result == nil {
		result = []models.SelectObject{}
	}

	return &result, nil
}

func (r *RoutePostgres) GetRouteOrder(agent_id, day_id int, language string) ([]models.SelectObject, error) {
	var result []models.SelectObject
	var trade_points []models.TradePoint

	get_trade_points_request := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").
		Where("trade_agent_id = ?", agent_id).
		Where("day_id = ?", day_id).Find(&trade_points)

	orderRepo := OrderPostgres{db: r.db}
	trade_points = orderRepo.ADDExtraTradePoints(agent_id, day_id, trade_points)

	if get_trade_points_request.Error != nil {
	}

	for index, element := range trade_points {
		if element.Status == false {
			trade_points[index].Status = true
		}
	}

	var route_orders []models.RouteOrder

	get_route_orders_request := r.db.Model(&models.RouteOrder{}).Where("trade_agent_id = ?", agent_id).
		Order("order_number").
		Where("day_id = ?", day_id).Find(&route_orders)

	if get_route_orders_request.Error != nil {
	}
	var order_index int

	for _, route_element := range route_orders {
		for point_index, point_element := range trade_points {
			if point_element.Id == route_element.TradePointId {
				order_index++
				var data models.SelectObject
				data.Id = point_element.Id
				if language == "ru" {
					data.Name = point_element.Code + "| " + point_element.NameRu + ";" + point_element.LocationRu + ";" + point_element.OrientirRu
				} else {
					data.Name = point_element.Code + "| " + point_element.NameTm + ";" + point_element.LocationTm + ";" + point_element.OrientirTm
				}
				result = append(result, data)
				trade_points[point_index].Status = false
				break
			}
		}
	}

	for point_index, point_element := range trade_points {
		if point_element.Status == true {
			var data models.SelectObject
			order_index++
			data.Id = point_element.Id
			if language == "ru" {
				data.Name = point_element.Code + "| " + point_element.NameRu + ";" + point_element.LocationRu + ";" + point_element.OrientirRu
			} else {
				data.Name = point_element.Code + "| " + point_element.NameTm + ";" + point_element.LocationTm + ";" + point_element.OrientirTm
			}
			result = append(result, data)
			trade_points[point_index].Status = false
		}
	}

	if result == nil {
		result = []models.SelectObject{}
	}

	return result, nil
}

func (r *RoutePostgres) PostRouteOrder(agent_id, day_id int, list []models.SelectObject) error {
	var route_orders []models.RouteOrder

	delete_route_order_request := r.db.Model(&models.RouteOrder{}).
		Where("trade_agent_id = ?", agent_id).
		Where("day_id = ?", day_id).Delete(&route_orders)

	if delete_route_order_request.Error != nil {
	}

	var order_number int
	for _, item := range list {
		order_number++
		new_route_order := models.RouteOrder{
			TradePointId: item.Id,
			TradeAgentId: agent_id,
			DayId:        day_id,
			Order_Number: order_number,
		}

		save_reuqest := r.db.Model(&models.RouteOrder{}).Create(&new_route_order)
		if save_reuqest.Error != nil {
			fmt.Println(save_reuqest.Error)
			return save_reuqest.Error
		}

	}

	return nil
}

func (r *RoutePostgres) ExcelRouteOrder(agent_id, day_id int, language string) ([]byte, error) {
	var list []models.SelectObject

	list, err := r.GetRouteOrder(agent_id, day_id, language)

	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "маршруты"
	f.NewSheet(sheet)

	headers := []string{"№", "Наименование торговой точки"}

	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for ronIdx, item := range list {
		row := ronIdx + 2 // Start from the second row
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), item.Id)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), item.Name)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
