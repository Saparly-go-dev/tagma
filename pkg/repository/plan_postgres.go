package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type PlanPostgres struct {
	db *gorm.DB
}

func NewPlanPostgres(db *gorm.DB) *PlanPostgres {
	return &PlanPostgres{db: db}
}

func (r *PlanPostgres) Create(list []models.CreatePlan, month, year int) error {
	var plan models.Plan

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)

	get_plan_request := r.db.Model(&models.Plan{}).
		Where("id >0").
		Where("created_at >= ?", firstOfMonth).
		Where("created_at <= ?", lastOfMonth).
		First(&plan)

	if get_plan_request.Error != nil {

	}

	if plan.Id > 0 {
		return errors.New("Rugsat berilmeýär")
	} //

	for _, item := range list {
		for idx, element := range item.List {
			newPlan := models.Plan{
				ProductId:    item.ProductId,
				TradeAgentId: idx,
				Count:        element,
				CreatedAt:    firstOfMonth,
				Status:       true,
			}

			save_plan_requst := r.db.Model(&models.Plan{}).Create(&newPlan)
			if save_plan_requst.Error != nil {
				fmt.Println(save_plan_requst.Error)
				return save_plan_requst.Error
			} else {
				//fmt.Print("succesfully saved")
			}
		}
	}

	return nil
}

func (r *PlanPostgres) GetPlans(month, year int, language string) ([]models.ReadPlan, error) {
	var plans []models.Plan
	var result []models.ReadPlan

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	//firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	//lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	//fmt.Println(firstOfMonth)

	get_plans_request := r.db.Model(&models.Plan{}).
		Preload("Product").
		Preload("TradeAgent").
		Where("created_at = ?", firstOfMonth).
		Find(&plans)

	//fmt.Println(plans)

	if get_plans_request.Error != nil {
		return []models.ReadPlan{}, get_plans_request.Error
	}

	var products []models.Product

	get_products_request := r.db.Model(&models.Product{}).Where("id > 0 ").Find(&products)
	if get_products_request.Error != nil {
		return []models.ReadPlan{}, get_products_request.Error
	}

	for _, product := range products {
		var data models.ReadPlan
		data.ProductId = product.Id
		data.List = make(map[int]int)
		for _, plan := range plans {
			if plan.ProductId == product.Id {
				data.List[plan.TradeAgentId] = plan.Count
			}
		}

		if len(data.List) > 0 {
			result = append(result, data)
		}
	}

	if len(result) == 0 {
		result = []models.ReadPlan{}
	}

	return result, nil
}

func (r *PlanPostgres) GetPlanStatisticsForMain(month, year, agent_id int, language string) ([]models.PlanStatisticsMain, error) {
	var plans []models.Plan
	var result []models.PlanStatisticsMain

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)

	get_plans_request := r.db.Model(&models.Plan{}).
		Preload("Product").
		Preload("TradeAgent").
		Where("id >0").
		Where("created_at >= ?", firstOfMonth).
		Where("created_at <= ?", lastOfMonth)

	if agent_id > 0 {
		get_plans_request = get_plans_request.Where("trade_agent_id = ?", agent_id)
	}

	get_plans_request.Find(&plans)

	if get_plans_request.Error != nil {
		return []models.PlanStatisticsMain{}, get_plans_request.Error
	}

	var orders []models.Order

	get_orders_request := r.db.Model(&models.Order{}).
		Preload("List").
		Where("created_at >= ?", firstOfMonth).
		Where("created_at <= ?", lastOfMonth)

	if agent_id != 0 {
		get_orders_request = get_orders_request.Where("trade_agent_id = ?", agent_id)
	}

	get_orders_request.Find(&orders)

	if get_orders_request.Error != nil {
		return []models.PlanStatisticsMain{}, get_orders_request.Error
	}

	var products []models.Product

	get_products_request := r.db.Model(&models.Product{}).
		Where("status = ?", true).Order("id").Find(&products)

	if get_products_request.Error != nil {
		return []models.PlanStatisticsMain{}, get_products_request.Error
	}

	for _, item := range products {
		var plan_count int
		for _, plan_element := range plans {
			if item.Id == plan_element.ProductId {
				plan_count = plan_count + plan_element.Count
			}
		}

		var data models.PlanStatisticsMain
		day_list := GetWorkingDaysInMonth(month, year)
		//var result_data models.PlanStatisticsMain
		for idx, day_element := range day_list {
			var count int
			for _, order := range orders {
				order_date := time.Date(year, time.Month(month), day_element.Day, 0, 0, 0, 0, time.UTC)

				if order.CreatedAt.Format("02.01.2006") == order_date.Format("02.01.2006") {
					for _, order_list_element := range order.List {
						if item.Id == order_list_element.ProductId {
							count = count + order_list_element.Count
						}
					}
				}

			}
			day_list[idx].Count = count
		}
		data.Plan = plan_count
		data.ProductId = item.Id
		data.Data = day_list
		if language == "ru" {
			data.Product = item.TasteRu
		} else {
			data.Product = item.TasteTm
		}
		pass_day_count := GetWorkingPassDaysCount(month, year)
		count_sum := GetCount(day_list)
		if count_sum == 0 || pass_day_count == 0 || plan_count == 0 {
			data.Forecast = 0
			data.ForecastVsPlan = 0
		} else {
			data.Forecast = int(float64(count_sum) / float64(pass_day_count) * float64((len(day_list))))
			data.ForecastVsPlan = Round2(float64(data.Forecast)/float64(plan_count)-1) * 100
		}

		result = append(result, data)
	}
	orderRepo := OrderPostgres{db: r.db}
	product_distinct, err := orderRepo.GetProducts(language)

	if err != nil {
	}

	for _, distinct_element := range *product_distinct {
		for j, result_element := range result {
			if distinct_element.Id == result_element.ProductId {
				result[j].Name = distinct_element.Name
				result[j].Volume = distinct_element.Volume
				result[j].ProductType = distinct_element.ProductTypeName
				result[j].GroupId = distinct_element.GroupId
			}
		}
	}

	if result == nil || len(result) == 0 {
		return []models.PlanStatisticsMain{}, nil
	}

	return result, nil
}

func (r *PlanPostgres) GetPlansForAgent(month, year, agent_id int, language string) ([]models.PlanStatisticsForAgent, error) {
	var plans []models.Plan
	var result []models.PlanStatisticsForAgent

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)

	get_plans_request := r.db.Model(&models.Plan{}).
		Preload("Product").
		Preload("TradeAgent").
		Where("id >0").
		Where("created_at >= ?", firstOfMonth).
		Where("created_at <= ?", lastOfMonth).
		Where("trade_agent_id = ?", agent_id).
		Find(&plans)

	if get_plans_request.Error != nil {
		return []models.PlanStatisticsForAgent{}, get_plans_request.Error
	}

	var orders []models.Order

	get_orders_request := r.db.Model(&models.Order{}).
		Preload("List").
		Where("created_at >= ?", firstOfMonth).
		Where("created_at <= ?", lastOfMonth).
		Where("trade_agent_id = ?", agent_id).
		Find(&orders)

	if get_orders_request.Error != nil {
		return []models.PlanStatisticsForAgent{}, get_orders_request.Error
	}

	var products []models.Product

	get_products_request := r.db.Model(&models.Product{}).
		Where("status = ?", true).Order("id").Find(&products)

	if get_products_request.Error != nil {
		return []models.PlanStatisticsForAgent{}, get_products_request.Error
	}

	for _, product := range products {
		var count int
		for _, order := range orders {
			for _, item := range order.List {
				if product.Id == item.ProductId {
					count = count + item.Count
				}
			}
		}

		var data models.PlanStatisticsForAgent
		data.ProductId = product.Id
		data.Selled = count

		result = append(result, data)
	}

	for _, element := range plans {
		for idx, result_element := range result {
			if element.ProductId == result_element.ProductId {
				result[idx].Count = element.Count
			}
		}
	}

	return result, nil
}

func (r *PlanPostgres) Update(month, year int, list []models.UpdatedPlan) error {

	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	for _, item := range list {
		for idx, element := range item.List {
			var plan models.Plan

			get_plan_request := r.db.Model(&models.Plan{}).Where("product_id = ?", item.ProductId).
				Where("trade_agent_id = ?", idx).Where("created_at = ?", firstOfMonth).Find(&plan)

			if get_plan_request.Error != nil {
				continue
			}

			if plan.Id == 0 {
				newPlan := models.Plan{
					ProductId:    item.ProductId,
					TradeAgentId: idx,
					Count:        element,
					CreatedAt:    firstOfMonth,
					Status:       true,
				}

				save_request := r.db.Model(&models.Plan{}).Create(&newPlan)
				if save_request.Error != nil {
					continue
				}
			} else {

				updateData := map[string]interface{}{
					"count":  element,
					"status": true,
				}

				update_request := r.db.Model(&models.Plan{}).Where("id = ?", plan.Id).Updates(updateData)
				if update_request.Error != nil {
					continue
				}
			}

		}
	}

	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////
func GetWorkingDaysInMonth(month, year int) []models.PlanData {
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	var result []models.PlanData
	currentDate := firstOfMonth

	for currentDate.Month() == time.Month(month) {
		weekday := currentDate.Weekday()

		if weekday != time.Sunday {
			var data models.PlanData
			data.Day = currentDate.Day()
			data.Count = 0
			result = append(result, data)

		}

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return result
}

func GetWorkingPassDaysCount(month, year int) int {
	firstOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	var result int
	currentDate := firstOfMonth
	now := time.Now().In(time.UTC)

	for currentDate.Month() == time.Month(month) {
		if currentDate.After(now) {
			break
		}

		weekday := currentDate.Weekday()

		if weekday != time.Sunday {
			result++
		}

		currentDate = currentDate.AddDate(0, 0, 1)

	}
	return result
}

func GetCount(list []models.PlanData) int {
	var result int
	for _, item := range list {
		result += item.Count
	}

	return result
}
