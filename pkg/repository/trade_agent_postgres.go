package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type TradeAgentPostgres struct {
	db *gorm.DB
}

func NewTradeAgentPostgres(db *gorm.DB) *TradeAgentPostgres {
	return &TradeAgentPostgres{db: db}
}

func (r *TradeAgentPostgres) Create(agent models.CreateTradeAgent) (int, error) {
	var lastdata models.TradeAgent
	var codeindex int

	lastresult := r.db.Model(models.TradeAgent{}).Last(&lastdata).Order("id desc")

	if lastresult.Error != nil {
		codeindex = 1
	} else {
		codeindex = lastdata.Id + 1
	}

	newAgent := models.TradeAgent{
		Code:      strconv.Itoa(codeindex),
		NameRu:    agent.NameRu,
		NameTm:    agent.NameTm,
		Number:    agent.Number,
		Status:    agent.Status,
		UpdatedAt: time.Now(),
	}

	result := r.db.Create(&newAgent)
	if result.Error != nil {
		return 0, result.Error
	}

	for _, driverId := range agent.DriverIds {
		var data models.AgentsDrivers
		data.TradeAgentId = newAgent.Id
		data.DriverId = driverId

		request := r.db.Model(&models.AgentsDrivers{}).Create(&data)
		if request.Error != nil {
		}
	}

	for _, ekspeditorId := range agent.EkspeditorIds {
		var data models.AgentsEkspeditors
		data.TradeAgentId = newAgent.Id
		data.EkspeditorId = ekspeditorId

		request := r.db.Model(&models.AgentsEkspeditors{}).Create(&data)
		if request.Error != nil {
		}
	}

	for _, merchanId := range agent.MerchandiserIds {
		var data models.AgentsMerchandisers
		data.TradeAgentId = newAgent.Id
		data.MerchandiserId = merchanId

		request := r.db.Model(&models.AgentsMerchandisers{}).Create(&data)
		if request.Error != nil {

		}
	}

	return newAgent.Id, nil
}

func (r *TradeAgentPostgres) AgentExcel(name, language string) ([]byte, error) {
	var agents []models.TradeAgent
	var readAgents []models.ReadTradeAgent

	result := r.db.Model(&models.TradeAgent{}).Preload("Driver").
		Preload("Ekspeditor").Preload("Merchandiser").
		Where("name_ru ILIKE ? ", "%"+name+"%").Or("name_tm ILIKE ? ", "%"+name+"%").Order("id").
		Find(&agents)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, agent := range agents {
		var newitem models.ReadTradeAgent
		var driverlist []models.SelectObject
		var ekspeditorslist []models.SelectObject
		var merchandisers []models.SelectObject

		newitem.Id = agent.Id
		newitem.Code = agent.Code
		newitem.NameRu = agent.NameRu
		newitem.NameTm = agent.NameTm
		newitem.Number = agent.Number
		newitem.Status = agent.Status

		for _, driver := range agent.Driver {
			var rawdata = models.SelectObject{}
			rawdata.Id = driver.Id
			if language == "ru" {
				rawdata.Name = driver.NameRu
			} else {
				rawdata.Name = driver.NameTm
			}
			driverlist = append(driverlist, rawdata)
		}

		for _, ekspeditor := range agent.Ekspeditor {
			var rawdata = models.SelectObject{}
			rawdata.Id = ekspeditor.Id
			if language == "ru" {
				rawdata.Name = ekspeditor.NameRu
			} else {
				rawdata.Name = ekspeditor.NameTm
			}
			ekspeditorslist = append(ekspeditorslist, rawdata)
		}

		for _, merchandiser := range agent.Merchandiser {
			var rawdata = models.SelectObject{}
			rawdata.Id = merchandiser.Id
			if language == "ru" {
				rawdata.Name = merchandiser.NameRu
			} else {
				rawdata.Name = merchandiser.NameTm
			}

			merchandisers = append(merchandisers, rawdata)

		}
		if driverlist == nil {
			driverlist = []models.SelectObject{}
		}

		if ekspeditorslist == nil {
			ekspeditorslist = []models.SelectObject{}
		}

		if merchandisers == nil {
			merchandisers = []models.SelectObject{}
		}

		newitem.UpdatedAt = agent.UpdatedAt.Format("02.01.2006")
		newitem.Drivers = driverlist
		newitem.Ekspeditors = ekspeditorslist
		newitem.Merchandisers = merchandisers

		readAgents = append(readAgents, newitem)
	}

	f := excelize.NewFile()
	sheet := "Торговые представители"
	f.NewSheet(sheet)

	headers := []string{"ФИО на русском языке", "ФИО на туркменском языке", "Номер телефона", "Экспедиторы", "Водители", "Мерчендайзер", "Долги"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, agent := range readAgents {
		row := rowIdx + 2
		f.SetCellValue(sheet, "A"+strconv.Itoa(row), agent.NameRu)
		f.SetCellValue(sheet, "B"+strconv.Itoa(row), agent.NameTm)
		f.SetCellValue(sheet, "C"+strconv.Itoa(row), agent.Number)

		var ekspeditors string
		var drivers string
		var merchandisers string
		for _, item := range agent.Ekspeditors {
			ekspeditors = ekspeditors + item.Name + ";"
		}

		for _, item := range agent.Drivers {
			drivers = drivers + item.Name + ";"
		}

		for _, item := range agent.Merchandisers {
			merchandisers = merchandisers + item.Name + ";"
		}

		f.SetCellValue(sheet, "D"+strconv.Itoa(row), ekspeditors)
		f.SetCellValue(sheet, "E"+strconv.Itoa(row), drivers)
		f.SetCellValue(sheet, "F"+strconv.Itoa(row), merchandisers)
		f.SetCellValue(sheet, "G"+strconv.Itoa(row), agent.Debt)
	}

	// Save the file to a byte buffer
	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *TradeAgentPostgres) GetAll(pageSize, pageNumber int, name, language string) (*models.TradeAgentPage, error) {
	var agents []models.TradeAgent
	var readAgents []models.ReadTradeAgent

	offset := (pageNumber - 1) * pageSize

	result := r.db.Model(&models.TradeAgent{}).Preload("Driver").Preload("Ekspeditor").Preload("Merchandiser").
		Where("name_ru ILIKE ? ", "%"+name+"%").Or("name_tm ILIKE ? ", "%"+name+"%").Order("id").
		Offset(offset).Limit(pageSize).Find(&agents)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, agent := range agents {
		var newitem models.ReadTradeAgent
		var driverlist []models.SelectObject
		var ekspeditorslist []models.SelectObject
		var merchandisers []models.SelectObject
		var debt models.Debt

		get_debt_request := r.db.Model(&models.Debt{}).Where("trade_agent_id = ?", agent.Id).Find(&debt)

		if get_debt_request.Error != nil {
			newitem.Debt = 0
		} else {
			newitem.Debt = debt.Total
		}

		newitem.Id = agent.Id
		newitem.Code = agent.Code
		newitem.NameRu = agent.NameRu
		newitem.NameTm = agent.NameTm
		newitem.Number = agent.Number
		newitem.Status = agent.Status

		for _, driver := range agent.Driver {
			var rawdata = models.SelectObject{}
			rawdata.Id = driver.Id
			if language == "ru" {
				rawdata.Name = driver.NameRu
			} else {
				rawdata.Name = driver.NameTm
			}
			driverlist = append(driverlist, rawdata)
		}

		for _, ekspeditor := range agent.Ekspeditor {
			var rawdata = models.SelectObject{}
			rawdata.Id = ekspeditor.Id
			if language == "ru" {
				rawdata.Name = ekspeditor.NameRu
			} else {
				rawdata.Name = ekspeditor.NameTm
			}
			ekspeditorslist = append(ekspeditorslist, rawdata)
		}

		for _, merchandiser := range agent.Merchandiser {
			var rawdata = models.SelectObject{}
			rawdata.Id = merchandiser.Id
			if language == "ru" {
				rawdata.Name = merchandiser.NameRu
			} else {
				rawdata.Name = merchandiser.NameTm
			}

			merchandisers = append(merchandisers, rawdata)

		}
		if driverlist == nil {
			driverlist = []models.SelectObject{}
		}

		if ekspeditorslist == nil {
			ekspeditorslist = []models.SelectObject{}
		}

		if merchandisers == nil {
			merchandisers = []models.SelectObject{}
		}

		newitem.UpdatedAt = agent.UpdatedAt.Format("02.01.2006")
		newitem.Drivers = driverlist
		newitem.Ekspeditors = ekspeditorslist
		newitem.Merchandisers = merchandisers

		readAgents = append(readAgents, newitem)
	}
	var count int64

	result2 := r.db.Model(&models.TradeAgent{}).Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)

	if result2.Error != nil {
		//fmt.Println("result2 islemedi")
		return nil, result.Error
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if readAgents == nil {
		readAgents = []models.ReadTradeAgent{}
	}

	// Create the page object
	page := models.TradeAgentPage{
		Items:      readAgents,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
	}

	return &page, nil
}

func (r *TradeAgentPostgres) ChangeStatus(Id int) error {
	var data models.TradeAgent
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.TradeAgent{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.TradeAgent{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *TradeAgentPostgres) GetById(Id int) (models.TradeAgent, error) {
	var data models.TradeAgent

	result := r.db.Preload("Driver").Preload("Ekspeditor").Preload("Merchandiser").First(&data, Id)

	return data, result.Error
}

func (r *TradeAgentPostgres) Delete(Id int) error {
	result := r.db.Delete(&models.TradeAgent{}, Id)
	return result.Error
}

func (r *TradeAgentPostgres) GetDrivers(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.Driver

	result := r.db.Model(&models.Driver{}).Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

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

func (r *TradeAgentPostgres) GetEkspeditors(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.Ekspeditor

	result := r.db.Model(&models.Ekspeditor{}).Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

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

// Update modifies a city's information based on its ID
func (r *TradeAgentPostgres) Update(Id int, agent models.CreateTradeAgent) error {
	var agentdata models.TradeAgent

	resultdata := r.db.First(&agentdata, Id).Preload("Ekspeditor").Preload("Driver")

	if resultdata.Error != nil {
		return resultdata.Error
	}

	var oldDrivers []int
	var oldEkspeditors []int

	for _, item := range agentdata.Driver {
		oldDrivers = append(oldDrivers, item.Id)
	}

	for _, item := range agentdata.Ekspeditor {
		oldEkspeditors = append(oldEkspeditors, item.Id)
	}

	// Log-a ýazmak üçin öňki bilen täze torgowylar ýa-da sürüjiler deňeşdirilýär tapawut bar bolsa bellenilýär
	r.CheckDrivers(agentdata.Id, oldDrivers, agent.DriverIds)
	r.CheckEkspeditor(agentdata.Id, oldEkspeditors, agent.EkspeditorIds)

	updateData := map[string]interface{}{
		"name_ru": agent.NameRu,
		"name_tm": agent.NameTm,
		"number":  agent.Number,
		"status":  agent.Status,
	}

	result := r.db.Model(&models.TradeAgent{}).Where("id = ?", Id).Updates(updateData)

	deleteEkspeditors := r.db.Model(&models.AgentsEkspeditors{}).Where("trade_agent_id = ?", Id).Delete(&models.AgentsEkspeditors{})
	if deleteEkspeditors.Error != nil {
	}
	deleteDrivers := r.db.Model(&models.AgentsDrivers{}).Where("trade_agent_id = ?", Id).Delete(&models.AgentsDrivers{})
	if deleteDrivers.Error != nil {
	}

	deleteMerchandirers := r.db.Model(&models.AgentsMerchandisers{}).Where("trade_agent_id = ?", Id).Delete(&models.AgentsMerchandisers{})
	if deleteMerchandirers.Error != nil {

	}

	for _, driverId := range agent.DriverIds {
		var data models.AgentsDrivers
		data.TradeAgentId = Id
		data.DriverId = driverId

		request := r.db.Model(&models.AgentsDrivers{}).Create(&data)
		if request.Error != nil {
		}
	}

	for _, ekspeditorId := range agent.EkspeditorIds {
		var data models.AgentsEkspeditors
		data.TradeAgentId = Id
		data.EkspeditorId = ekspeditorId

		request := r.db.Model(&models.AgentsEkspeditors{}).Create(&data)
		if request.Error != nil {
		}
	}

	for _, merchandiserId := range agent.MerchandiserIds {
		var data models.AgentsMerchandisers
		data.TradeAgentId = Id
		data.MerchandiserId = merchandiserId

		request := r.db.Model(&models.AgentsMerchandisers{}).Create(&data)
		if request.Error != nil {

		}
	}

	return result.Error
}

type CheckingObject struct {
	Id    int
	Count int
}

func (r *TradeAgentPostgres) Save(infoRu, infoTm, logType string) {
	data := models.Log{
		InfoRu:   infoRu,
		InfoTm:   infoTm,
		LogType:  logType,
		User:     "admin",
		CreateAt: time.Now(),
	}

	result := r.db.Model(&models.Log{}).Create(&data)

	if result.Error != nil {
		log.Println(result.Error)
	}
}

func (r *TradeAgentPostgres) GetAgentDebt(AgentId int, language string) (models.AgentDebt, error) {
	var orders []models.Order
	var agent models.TradeAgent
	var result models.AgentDebt
	var list []models.AgentDebtItem
	var sum float64
	var cash float64
	var terminal float64
	var transfer float64
	var not_paid_cash float64
	var not_paid float64

	get_agent_request := r.db.Model(&models.TradeAgent{}).Where("id = ?", AgentId).Find(&agent)
	if get_agent_request.Error != nil {
		log.Println(get_agent_request.Error)
		return result, get_agent_request.Error
	}

	get_orders_request := r.db.Model(&models.Order{}).Where("trade_agent_id = ? AND is_closed = ?", AgentId, false).Find(&orders)
	if get_orders_request.Error != nil {
		log.Println(get_orders_request.Error)
	}

	for _, order := range orders {
		var payments []models.Payment
		var payment_sum float64

		var order_cash float64
		var order_transfer float64
		var order_terminal float64

		get_payments_request := r.db.Model(&models.Payment{}).Where("order_id = ?", order.Id).Find(&payments)

		if get_payments_request.Error != nil {
			fmt.Println(get_payments_request.Error)
		}

		for _, payment_element := range payments {
			if payment_element.Status == true {
				payment_sum += payment_element.Currency
			} else {
				if payment_element.PaymentTypeId == 1 {
					order_cash += payment_element.Currency
				}

				if payment_element.PaymentTypeId == 2 {
					order_terminal += payment_element.Currency
				}

				if payment_element.PaymentTypeId == 3 {
					order_transfer += payment_element.Currency
				}
			}
		}

		if payment_sum >= order.Sum {
			continue
		}

		cash += order_cash
		terminal += order_terminal
		transfer += order_transfer

		var tradePoint models.TradePoint

		pointRequest := r.db.Model(&models.TradePoint{}).Preload("City").Preload("District").Where("id = ?", order.TradePointId).Find(&tradePoint)
		if pointRequest.Error != nil {
			log.Println(pointRequest.Error)
		}

		var paid float64

		if get_payments_request.Error != nil {
			fmt.Println(get_payments_request.Error)
		}

		for _, payment_element := range payments {
			if payment_element.Status == true {
				paid += payment_element.Currency
			}
		}

		not_paid = order.Sum - paid

		var data models.AgentDebtItem
		data.Sum = order.Sum
		data.Paid = paid
		data.NotPaid = not_paid
		data.Status = order.IsClosed
		data.CreatedAt = order.CreatedAt.Format("01.02.2006")

		if language == "ru" {
			data.Name = tradePoint.NameRu
			data.City = tradePoint.City.NameRu
			data.District = tradePoint.District.NameRu
		} else {
			data.Name = tradePoint.NameTm
			data.City = tradePoint.City.NameTm
			data.District = tradePoint.District.NameTm
		}

		list = append(list, data)
		not_paid_cash = not_paid_cash + (not_paid - order_transfer - order_terminal - order_cash)
		sum = sum + not_paid
	}

	if list == nil {
		list = []models.AgentDebtItem{}
	}
	if language == "ru" {
		result.Name = agent.NameRu
	} else {
		result.Name = agent.NameTm
	}

	var self_debt models.Debt
	get_agent_self_debt := r.db.Model(&models.Debt{}).Where("trade_agent_id = ?", AgentId).Find(&self_debt)
	if get_agent_self_debt.Error != nil {
		self_debt = models.Debt{}
	}

	result.Number = agent.Number
	result.Sum = sum
	result.Cash = cash + not_paid_cash
	result.Terminal = terminal
	result.Transfer = transfer
	result.Self = self_debt.Total
	result.List = list

	return result, nil
}

func (r *TradeAgentPostgres) CheckDrivers(Id int, old, new []int) {
	var array []CheckingObject
	var agent models.TradeAgent

	agentRequest := r.db.Model(&models.TradeAgent{}).Where("id = ?", Id).Find(&agent)

	if agentRequest.Error != nil {
		log.Println(agentRequest.Error)
	}

	for _, value := range old {
		var data CheckingObject
		data.Id = value
		data.Count = 0

		array = append(array, data)
	}

	for _, value := range new {
		found := false
		for i := range array {
			if array[i].Id == value {
				// If found, increment count by 1
				array[i].Count += 1
				found = true
				break
			}
		}
		if !found {
			// If not found, add a new element
			array = append(array, CheckingObject{
				Id:    value,
				Count: 2, // Start with a count of 2
			})
		}
	}

	deletedRu := "("
	deletedTm := "("
	addRu := "("
	addTm := "("
	//fmt.Println(array)

	for _, value := range array {
		var driver models.Driver

		request := r.db.Model(&models.Driver{}).Where("id = ?", value.Id).Find(&driver)
		if request.Error != nil {
			continue
		}
		if value.Count == 0 {
			deletedRu += driver.NameRu + ";"
			deletedTm += driver.NameTm + ";"
		}

		if value.Count == 2 {
			addRu += driver.NameRu + ";"
			addTm += driver.NameTm + ";"
		}
	}

	deletedRu += ")"
	deletedTm += ")"
	addRu += ")"
	addTm += ")"

	r.Save(fmt.Sprintf("Водитель %s, для торгового представилтеля:%s был удален и добавлен водитель %s.", deletedRu, agent.NameRu, addRu),
		fmt.Sprintf("%s - atly söwda wekiliň sürüjisi %s aýrylyp, %s-goşuldy", agent.NameTm, deletedTm, addTm),
		"ChangeDriver")

}

func (r *TradeAgentPostgres) CheckEkspeditor(Id int, old, new []int) {
	var array []CheckingObject
	var agent models.TradeAgent

	agentRequest := r.db.Model(&models.TradeAgent{}).Where("id = ?", Id).Find(&agent)

	if agentRequest.Error != nil {
		log.Println(agentRequest.Error)
	}

	for _, value := range old {
		var data CheckingObject
		data.Id = value
		data.Count = 0

		array = append(array, data)
	}

	for _, value := range new {
		found := false
		for i := range array {
			if array[i].Id == value {
				// If found, increment count by 1
				array[i].Count += 1
				found = true
				break
			}
		}
		if !found {
			// If not found, add a new element
			array = append(array, CheckingObject{
				Id:    value,
				Count: 2, // Start with a count of 2
			})
		}
	}

	deletedRu := "("
	deletedTm := "("
	addRu := "("
	addTm := "("
	//fmt.Println(array)

	for _, value := range array {
		var driver models.Driver

		request := r.db.Model(&models.Ekspeditor{}).Where("id = ?", value.Id).Find(&driver)
		if request.Error != nil {
			continue
		}
		if value.Count == 0 {
			deletedRu += driver.NameRu + ";"
			deletedTm += driver.NameTm + ";"
		}

		if value.Count == 2 {
			addRu += driver.NameRu + ";"
			addTm += driver.NameTm + ";"
		}
	}

	deletedRu += ")"
	deletedTm += ")"
	addRu += ")"
	addTm += ")"

	r.Save(fmt.Sprintf("Экспедитор %s, для торгового представилтеля:%s был удален и добавлен %s.", deletedRu, agent.NameRu, addRu),
		fmt.Sprintf("%s - atly söwda wekiliň ekspeditory %s aýrylyp, %s-goşuldy", agent.NameTm, deletedTm, addTm),
		"ChangeEkspeditor")

}

func (r *TradeAgentPostgres) GetMerchendisingInformation(day, month, year, agentid int, language string) (models.PointsMerchendising_List, error) {
	result := models.PointsMerchendising_List{}

	var tradePoints []models.TradePoint
	var result_list []models.PointsMerchendising_Item
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}

	//today := time.Now().Truncate(24 * time.Hour)
	request := r.db.Model(models.TradePoint{}).Preload("City").Preload("District").Where("status = ?", true).
		Where("trade_agent_id = ?", agentid).Where("day_id", dayOfWeek).Find(&tradePoints)

	if request.Error != nil {
		return result, request.Error
	}

	var db_furnitures []models.Furniture

	get_db_furnitures := r.db.Model(&models.Furniture{}).Where("id > 0").Find(&db_furnitures)
	if get_db_furnitures.Error != nil {
		fmt.Println(get_db_furnitures.Error)
	}

	//fmt.Println("furnitures", db_furnitures)

	result_furniture_list := []models.FurnitureInformation{}
	trade_point_count := len(tradePoints)

	for _, item := range db_furnitures {
		var new_furniture_list_element models.FurnitureInformation

		new_furniture_list_element.Count = 0
		new_furniture_list_element.Percentage = 0
		if language == "ru" {
			new_furniture_list_element.Name = item.NameRu
		} else {
			new_furniture_list_element.Name = item.NameTm
		}
		new_furniture_list_element.Item_percentage_value = 100 / float64(trade_point_count)

		result_furniture_list = append(result_furniture_list, new_furniture_list_element)
	}

	var routes []models.RouteOrder

	all_planogramma_true := 0
	all_planogramma_false := 0

	for _, tradePoint := range tradePoints {
		var data models.PointsMerchendising_Item
		var element_marketing_true_value, element_marketing_false_value int64

		data.Id = tradePoint.Id
		data.Code = tradePoint.Code

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

		var furnitures []models.TradePointsFurniture

		furniture_list := []models.FurnitureInformation{}

		get_furniture_request := r.db.Model(&models.TradePointsFurniture{}).Preload("Furniture").Where("trade_point_id = ?", tradePoint.Id).Find(&furnitures)
		if get_furniture_request.Error != nil {
			fmt.Println(get_furniture_request.Error)
		}

		get_marketing_true_value := r.db.Model(&models.Marketing{}).Where("trade_point_id = ? AND exposure = ?", tradePoint.Id, true).
			Where("DATE(created_at) = ?", date.Format("2006-01-02")).Count(&element_marketing_true_value)

		if get_marketing_true_value.Error != nil {
			element_marketing_true_value = 0
		}

		get_marketing_false_value := r.db.Model(&models.Marketing{}).Where("trade_point_id = ? AND exposure = ?", tradePoint.Id, false).
			Where("DATE(created_at) = ?", date.Format("2006-01-02")).Count(&element_marketing_false_value)

		if get_marketing_false_value.Error != nil {
			element_marketing_false_value = 0
		}

		if element_marketing_false_value == 0 && element_marketing_true_value == 0 {
			data.Planogramma = 0
			data.Ishave = false
		} else {
			data.Ishave = true
			marketing_sum := float64(element_marketing_true_value + element_marketing_false_value)
			data.Planogramma = float64(100/marketing_sum) * float64(element_marketing_true_value)
			all_planogramma_true += int(element_marketing_true_value)
			all_planogramma_false += int(element_marketing_false_value)
		}

		for _, item := range furnitures {
			var furniture_item models.FurnitureInformation
			if language == "ru" {
				furniture_item.Name = item.Furniture.NameRu
			} else {
				furniture_item.Name = item.Furniture.NameTm
			}

			furniture_item.Count = item.Count
			furniture_list = append(furniture_list, furniture_item)

			for index, item2 := range result_furniture_list {
				if furniture_item.Name == item2.Name {
					if furniture_item.Count > 0 {
						result_furniture_list[index].Count += furniture_item.Count
						result_furniture_list[index].Percentage += result_furniture_list[index].Item_percentage_value
					}
				}
			}
		}

		data.Furniture = furniture_list

		result_list = append(result_list, data)
	}

	var new_result []models.PointsMerchendising_Item

	for _, route := range routes {
		for _, data := range result_list {
			if route.TradePointId == data.Id {
				new_result = append(new_result, data)
			}
		}
	}

	for _, data := range result_list {
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

	marketing_percentage := 100 / float64(all_planogramma_true+all_planogramma_false)
	marketing_percentage = marketing_percentage * float64(all_planogramma_true)
	marketing_percentage = roundToTwoDecimals(marketing_percentage)

	marketing_information := &models.FurnitureInformation{
		Name:       "Plannogramma",
		Count:      0,
		Percentage: marketing_percentage,
	}

	result_furniture_list = append(result_furniture_list, *marketing_information)
	result.Furniture = result_furniture_list
	result.List = result_list

	return result, nil
}

func (r *TradeAgentPostgres) CreateCSVFileForFactory(day, month, year, agentid int, language string) ([]byte, error) {
	var orders []models.Order
	var list_products []models.ReadProductForTradePoint
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	getOrders := r.db.Model(&models.Order{}).
		Where("trade_agent_id = ? AND created_at = ?", agentid, date).
		Find(&orders)

	if getOrders.Error != nil {
		return nil, getOrders.Error
	}

	for _, order := range orders {
		var orderList []models.OrderList

		getOrderList := r.db.Model(&models.OrderList{}).Where("order_id = ?", order.Id).Find(&orderList)

		if getOrderList.Error != nil {
			return nil, getOrderList.Error
		}

		for _, order_list_element := range orderList {
			ishave := false
			for p_idx, product_element := range list_products {
				if order_list_element.ProductId == product_element.Id {
					list_products[p_idx].Count += order_list_element.Count
					ishave = true
				}
			}
			if ishave == false {
				var product models.Product

				getProduct := r.db.Model(models.Product{}).Preload("Name").Preload("ProductType").Order("id").Where("id = ?", order_list_element.ProductId).First(&product)

				if getProduct.Error != nil {
					return nil, getProduct.Error
				}

				data := models.ReadProductForTradePoint{
					Id:     product.Id,
					Name:   product.Name.Name,
					Volume: product.ProductType.Volume,
					Count:  order_list_element.Count,
				}

				if language == "ru" {
					data.Taste = product.TasteRu
					data.ProductType = product.ProductType.NameRu
				} else {
					data.Taste = product.TasteTm
					data.ProductType = product.ProductType.NameTm
				}

				list_products = append(list_products, data)
			}
		}
	}

	buf := &bytes.Buffer{}
	csvWriter := csv.NewWriter(buf)

	if language == "ru" {
		if err := csvWriter.Write([]string{"Бренд", "Вкус", "Тип", "Объем", "Количество"}); err != nil {
			return nil, err
		}
	} else {
		if err := csvWriter.Write([]string{"Brand", "Tagam", "Type", "Göwrüm", "Sany"}); err != nil {
			return nil, err
		}
	}

	for _, product := range list_products {
		row := []string{
			product.Name,
			product.Taste,
			product.ProductType,
			fmt.Sprintf("%.2f", product.Volume),
			strconv.Itoa(product.Count),
		}

		if err := csvWriter.Write(row); err != nil {
			return nil, err
		}
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *TradeAgentPostgres) GetTradeAgentRevenue(month, year, agent_id int, language string) (float64, error) {
	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	var payments []models.Payment

	get_payments_reqeust := r.db.Model(&models.Payment{}).Where("trade_agent_id = ? AND created_at > ?", agent_id, date).Find(&payments)

	if get_payments_reqeust.Error != nil {
		fmt.Println(get_payments_reqeust.Error)
		return 0, nil
	}

	var sum float64

	for _, payment := range payments {
		sum += payment.Currency
	}

	return sum, nil
}
