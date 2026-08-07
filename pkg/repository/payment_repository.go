package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PaymentPostgres struct {
	db *gorm.DB
}

func NewPaymentPostgres(db *gorm.DB) *PaymentPostgres {
	return &PaymentPostgres{db: db}
}

func (r *PaymentPostgres) GetTradePointsForPayment(day, month, year, agentId int, language string) (*[]models.TradePointForPayment, error) {
	var tradePoints []models.TradePoint
	var result []models.TradePointForPayment
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 {
		dayOfWeek = 7
	}
	//today := time.Now().Truncate(24 * time.Hour)
	request := r.db.Model(models.TradePoint{}).Preload("City").
		Preload("District").
		Where("trade_agent_id = ?", agentId).
		Where("day_id", dayOfWeek).
		Where("status = ?", true).
		Find(&tradePoints)

	if request.Error != nil {
		return nil, request.Error
	}

	orderRepo := OrderPostgres{db: r.db}
	tradePoints = orderRepo.ADDExtraTradePoints(agentId, dayOfWeek, tradePoints)
	tradePoints = orderRepo.ADDExtraOrderPoints(day, month, year, agentId, tradePoints)

	var routes []models.RouteOrder

	route_order_get_request := r.db.Model(&models.RouteOrder{}).Order("order_number").
		Where("trade_agent_id = ?", agentId).Where("day_id = ?", dayOfWeek).Find(&routes)

	if route_order_get_request.Error != nil {
		fmt.Println(route_order_get_request.Error)
	}

	for _, tradePoint := range tradePoints {
		var data models.TradePointForPayment

		var contact contact_object
		contact = orderRepo.GetTredePointContact(tradePoint.Id, language)

		var order models.Order
		request3 := r.db.Model(&models.Order{}).Where("trade_point_id = ?", tradePoint.Id).
			Where("DATE(created_at) = ?", date.Format("2006-01-02")).Find(&order)

		if request3.Error != nil {
			data.Status = false
		}

		var paid float64
		var payments []models.Payment

		get_payments_request := r.db.Model(&models.Payment{}).Where("order_id = ?", order.Id).Find(&payments)

		if get_payments_request.Error != nil {
			fmt.Println(get_payments_request.Error)
		}

		for _, payment_element := range payments {
			paid += payment_element.Currency
		}

		not_paid := order.Sum - paid
		paid = Round2(paid)
		not_paid = Round2(not_paid)

		data.Id = tradePoint.Id
		data.Code = tradePoint.Code
		data.OrderId = order.Id
		data.Sum = order.Sum
		data.Paid = paid
		data.NotPaid = not_paid
		data.Contact = contact.Name
		data.Number = contact.Number

		if language == "ru" {
			data.City = tradePoint.City.NameRu + ":" + tradePoint.District.NameTm
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

	var new_result []models.TradePointForPayment

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
		new_result = []models.TradePointForPayment{}
	}

	return &new_result, nil
}

func (r *PaymentPostgres) savePayment(agentId int, data models.CreatePayment) error {
	var order models.Order

	orderRequest := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&models.Order{}).Where("id = ?", data.OrderId).Find(&order)
	if orderRequest.Error != nil {
		return orderRequest.Error
	}

	if order.Id == 0 {
		return errors.New("Sargyt tapylmady")

	}
	if data.PaymentTypeId == 4 {
		var promotion models.Promotion
		var promotion_list []models.PromotionProducts
		var order_list []models.OrderList
		var product models.Product

		get_promotion_request := r.db.Model(&models.Promotion{}).
			Where(`"start" <= ?`, order.CreatedAt).
			Where(`"end" >= ?`, order.CreatedAt).Find(&promotion)

		if get_promotion_request.Error != nil || promotion.Id == 0 {
			return get_promotion_request.Error
		}

		get_promotion_list_request := r.db.Model(&models.PromotionProducts{}).
			Where("promotion_id = ?", promotion.Id).Find(&promotion_list)

		if get_promotion_list_request.Error != nil {
			return get_promotion_list_request.Error
		}

		get_order_list_request := r.db.Model(&models.OrderList{}).
			Where("order_id = ?", order.Id).Find(&order_list)

		if get_order_list_request.Error != nil {
			return get_order_list_request.Error
		}

		mark := 0
		for _, promotion_product := range promotion_list {
			ishave := false
			for _, order_product := range order_list {
				if promotion_product.PromotionId == order_product.ProductId {
					ishave = true

					if product.Id == 0 {
						get_promotion_product := r.db.Model(&models.Product{}).Where("id = ?", promotion_product.ProductId).Find(&product)
						if get_promotion_product.Error != nil {
							return get_promotion_product.Error
						}

					}
				}
			}

			if ishave == true {
				mark = mark + 1
			}
		}

		if mark == 0 {
			return nil
		} else {
			promotion_price := product.Price / float64(promotion.Count)

			newPayment := models.Payment{
				Currency:      promotion_price * data.Currency,
				OrderId:       order.Id,
				CreatedAt:     time.Now(),
				TradeAgentId:  agentId,
				Status:        false,
				IsAgent:       data.Is_Agent,
				PaymentTypeId: data.PaymentTypeId,
				Count:         promotion.Count,
			}

			saveRequest := r.db.Model(&models.Payment{}).Create(&newPayment)
			if saveRequest.Error != nil {
				fmt.Println(saveRequest.Error)
				return saveRequest.Error
			}
			if err := r.checkOrderForPaymentConfirm(newPayment.OrderId); err != nil {
				return err
			}
		}

	} else {
		var payments []models.Payment

		get_payments := r.db.Model(&models.Payment{}).Where("order_id = ?", order.Id).Find(&payments)
		if get_payments.Error != nil {
			return get_payments.Error
		}

		var payment_sum float64
		for _, payment_element := range payments {
			payment_sum = payment_sum + payment_element.Currency
		}

		if (payment_sum + data.Currency) > order.SaleSum {
			return errors.New("Rugsat berilmeýär")
		}
		newPayment := models.Payment{
			Currency:      data.Currency,
			OrderId:       order.Id,
			CreatedAt:     time.Now(),
			TradeAgentId:  agentId,
			Status:        false,
			IsAgent:       data.Is_Agent,
			PaymentTypeId: data.PaymentTypeId,
			Count:         0,
		}

		saveRequest := r.db.Model(&models.Payment{}).Create(&newPayment)
		if saveRequest.Error != nil {
			fmt.Println(saveRequest.Error)
			return saveRequest.Error
		}
		if err := r.checkOrderForPaymentConfirm(newPayment.OrderId); err != nil {
			return err
		}
	}

	return nil
}

func (r *PaymentPostgres) checkOrderForPaymentConfirm(orderId int) error {
	var order models.Order

	get_order_request := r.db.Model(&models.Order{}).Where("id = ?", orderId).Find(&order)
	if get_order_request.Error != nil {
		return get_order_request.Error
	}

	var payments []models.Payment
	get_payments_request := r.db.Model(&models.Payment{}).Where("order_id = ?", orderId).Find(&payments)
	if get_payments_request.Error != nil {
		return get_payments_request.Error
	}

	var sum float64

	for _, payment_element := range payments {
		sum = sum + payment_element.Currency
	}

	if sum == order.SaleSum && order.Status == true {
		updateData := map[string]interface{}{
			"is_closed": true,
		}

		update_request := r.db.Model(&models.Order{}).Where("id = ?", orderId).Updates(updateData)
		if update_request.Error != nil {
			return update_request.Error
		}
	}
	return nil
}

func (r *PaymentPostgres) GetPaymentTypes(language string) ([]models.SelectObject, error) {
	var result []models.SelectObject

	var types []models.PaymentType

	reqeust := r.db.Model(&models.PaymentType{}).Where("id > ?", 0).Find(&types)
	if reqeust.Error != nil {
		fmt.Println(reqeust.Error)
	}

	for _, item := range types {
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

	return result, nil
}

func (r *PaymentPostgres) GetPaymentForOrder(orderId int, language string) ([]models.ReadPayment, error) {
	var result []models.ReadPayment
	var list []models.Payment

	request := r.db.Model(&models.Payment{}).Preload("PaymentType").Preload("TradeAgent").Where("order_id = ?", orderId).Find(&list)

	if request.Error != nil {
		fmt.Println(request.Error)
	}

	for _, item := range list {
		var data models.ReadPayment
		data.Id = item.Id
		data.OrderId = item.OrderId
		data.CreatedAt = item.CreatedAt.Format("02.01.2006 15:04")
		data.PaymentTypeId = item.PaymentTypeId
		data.Currency = item.Currency
		data.IsAgent = item.IsAgent

		if language == "ru" {
			data.PaymentType = item.PaymentType.NameRu
			data.TradeAgent = item.TradeAgent.NameRu
		} else {
			data.PaymentType = item.PaymentType.NameTm
			data.TradeAgent = item.TradeAgent.NameTm
		}

		result = append(result, data)
	}

	if result == nil {
		result = []models.ReadPayment{}
	}

	return result, nil
}

func (r *PaymentPostgres) GetAgentPayments(agent_id, payment_type_id int, status bool, point_code, language string) ([]models.ReadPayment, error) {
	var result []models.ReadPayment
	var payments []models.Payment

	get_payments_request := r.db.Model(&models.Payment{}).Preload("PaymentType").Preload("TradeAgent").
		Where("trade_agent_id = ?", agent_id).
		Where("payment_type_id = ?", payment_type_id).
		Where("status = ?", status).Find(&payments)

	if get_payments_request.Error != nil {
	}

	for _, item := range payments {

		var order models.Order

		get_order_request := r.db.Model(&models.Order{}).Preload("TradePoint").Where("id = ?", item.OrderId).Find(&order)

		if get_order_request.Error != nil {

		}

		if len(point_code) > 0 {
			if order.TradePoint.Code != point_code {
				continue
			}
		}

		var data models.ReadPayment
		data.Id = item.Id
		data.Code = order.TradePoint.Code
		data.OrderId = item.OrderId
		data.CreatedAt = item.CreatedAt.Format("02.01.2006 15:04")
		data.PaymentTypeId = item.PaymentTypeId
		data.Currency = item.Currency
		data.IsAgent = item.IsAgent
		if data.PaymentTypeId == 4 {
			data.Currency = float64(item.Count)
		}
		if language == "ru" {
			data.TradePoint = order.TradePoint.NameRu + " " + order.TradePoint.OrientirRu + " " + order.TradePoint.LocationRu
			data.PaymentType = item.PaymentType.NameRu
			data.TradeAgent = item.TradeAgent.NameRu
		} else {
			data.TradePoint = order.TradePoint.NameTm + " " + order.TradePoint.OrientirTm + " " + order.TradePoint.LocationTm
			data.PaymentType = item.PaymentType.NameTm
			data.TradeAgent = item.TradeAgent.NameTm
		}

		result = append(result, data)
	}
	if result == nil {
		result = []models.ReadPayment{}
	}

	return result, nil
}

func (r *PaymentPostgres) ConfirmPayments(agent_id, price1, price2 int, ids []int) error {
	var sum float64
	currency := float64(price1) + float64(price2)/100

	for _, id := range ids {
		var payment models.Payment

		get_payments_request := r.db.Model(&models.Payment{}).Where("trade_agent_id = ?", agent_id).
			Where("id = ?", id).Find(&payment)

		if get_payments_request.Error != nil {
		}

		updateData := map[string]interface{}{
			"status": true,
		}

		update_payment_request := r.db.Model(&models.Payment{}).Where("id = ?", id).Updates(updateData)
		if update_payment_request.Error != nil {
		}

		sum += payment.Currency
	}

	if sum > currency {
		var debt models.Debt
		get_debt_request := r.db.Model(&models.Debt{}).Where("trade_agent_id = ?", agent_id).Find(&debt)
		if get_debt_request.Error != nil {
		}

		if debt.Id == 0 {
			newdebt := models.Debt{
				TradeAgentId: agent_id,
				Total:        Round2(sum - currency),
			}

			save_newdebt_request := r.db.Model(&models.Debt{}).Create(&newdebt)
			if save_newdebt_request.Error != nil {
			}
		} else {
			updateData := map[string]interface{}{
				"total": Round2(debt.Total + (sum - currency)),
			}

			update_request := r.db.Model(&models.Debt{}).Where("id = ?", debt.Id).Updates(updateData)
			if update_request.Error != nil {
			}
		}
	}

	return nil
}

func (r *PaymentPostgres) UpdatePayment(id, price1, price2 int) error {
	var payment models.Payment

	get_payment_request := r.db.Model(&models.Payment{}).Where("id = ?", id).Find(&payment)
	if get_payment_request.Error != nil {
		return get_payment_request.Error
	}

	if payment.Status == true {
		return errors.New("Rugsat berilmeýär")
	}

	currency := float64(price1) + float64(price2)/100
	currency = Round2(currency)

	updateData := map[string]interface{}{
		"currency": currency,
	}

	update_payment_request := r.db.Model(&models.Payment{}).Where("id = ?", id).Updates(updateData)
	if update_payment_request.Error != nil {
		return update_payment_request.Error
	}

	return nil
}

func (r *PaymentPostgres) UpdateDebt(agent_id, price1, price2 int) error {
	var debt models.Debt
	currency := float64(price1) + float64(price2)/100

	get_debt_request := r.db.Model(&models.Debt{}).Where("trade_agent_id = ?", agent_id).Find(&debt)
	if get_debt_request.Error != nil {
	}

	if debt.Total-currency < 0 {
		return errors.New("Rugsat berilmeýär")
	}

	updateData := map[string]interface{}{
		"total": Round2(debt.Total - currency),
	}

	update_debt_request := r.db.Model(&models.Debt{}).Where("id = ?", debt.Id).Updates(updateData)
	if update_debt_request.Error != nil {

	}
	return nil
}

func (r *PaymentPostgres) AddSelfDebt(agent_id int, currency float64) error {
	var count int64

	is_have_agent := r.db.Model(&models.TradeAgent{}).Where("id = ?", agent_id).Count(&count)

	if is_have_agent.Error != nil || count == 0 {
		return errors.New("Not Found")
	}

	var debt models.Debt

	get_debt_request := r.db.Model(&models.Debt{}).Where("trade_agent_id = ?", agent_id).Find(&debt)

	if get_debt_request.Error != nil {
		return errors.New("Not Found")
	}

	if debt.Id == 0 {
		newdebt := models.Debt{
			TradeAgentId: agent_id,
			Total:        currency,
		}

		save_newdebt_request := r.db.Model(&models.Debt{}).Create(&newdebt)
		if save_newdebt_request.Error != nil {
			return save_newdebt_request.Error
		}
	} else {
		updateData := map[string]interface{}{
			"total": debt.Total + currency,
		}

		updateRequest := r.db.Model(&models.Debt{}).Where("id = ?", debt.Id).Updates(updateData)
		if updateRequest.Error != nil {
			return updateRequest.Error
		}
	}

	return nil
}

func (r *PaymentPostgres) GetTradePointAllDebtSum(tradePointId int) float64 {
	var sum float64
	var orders []models.Order

	get_orders_request := r.db.Model(&models.Order{}).Where("is_closed = ? AND trade_point_id = ?", false, tradePointId).Find(&orders)

	if get_orders_request.Error != nil {
		fmt.Println(get_orders_request.Error)
		return 0
	}

	for _, order_item := range orders {
		var payments_sum float64
		var payments []models.Payment

		get_payments := r.db.Model(&models.Payment{}).Where("order_id = ?", order_item.Id).Find(&payments)

		if get_payments.Error != nil {
			fmt.Println(get_payments.Error)
		}

		for _, payment := range payments {
			payments_sum += payment.Currency
		}

		if payments_sum == order_item.Sum || payments_sum > order_item.Sum {
			continue
		}

		sum += order_item.Sum - payments_sum
	}

	return sum
}

func (r *PaymentPostgres) postPayment(currency float64, agent_id, tradepointId, PaymentTypeId int, is_agent bool) error {
	trade_point_debt := r.GetTradePointAllDebtSum(tradepointId)

	if currency > trade_point_debt {
		return errors.New("Rugsat berilmeýär")
	}

	var orders []models.Order

	get_orders_request := r.db.Model(&models.Order{}).Order("created_at").
		Where("is_closed = ? AND trade_point_id = ?", false, tradepointId).Find(&orders)

	if get_orders_request.Error != nil {
		return get_orders_request.Error
	}

	for _, order_item := range orders {
		var payments_sum float64
		var payments []models.Payment

		get_payments := r.db.Model(&models.Payment{}).Where("order_id = ?", order_item.Id).Find(&payments)

		if get_payments.Error != nil {
			return get_payments.Error
		}

		for _, payment := range payments {
			payments_sum += payment.Currency
		}

		if payments_sum == order_item.SaleSum || payments_sum > order_item.SaleSum {
			continue
		}

		payment_currency := order_item.SaleSum - payments_sum

		if payment_currency > currency {
			payment_currency = currency
		}

		payment_models := models.CreatePayment{
			PaymentTypeId: PaymentTypeId,
			Currency:      payment_currency,
			OrderId:       order_item.Id,
			Is_Agent:      is_agent,
		}

		if err := r.savePayment(agent_id, payment_models); err != nil {
			return err
		}

		currency = currency - payment_currency

		if currency == 0 {
			break
		}
	}

	return nil
}

func (r *PaymentPostgres) GetTradePointPaymentsHistory(agent_id, tradePointId int, language string) ([]models.ReadPayment, error) {
	var payments []models.Payment
	date_check := time.Now().Add(-24 * 30 * time.Hour)

	err := r.db.Model(&models.Payment{}).
		Where("DATE(payments.created_at) > ?", date_check.Format("2006-01-02")).
		Joins("JOIN orders ON orders.id = payments.order_id").
		Preload("PaymentType").
		Preload("TradeAgent").
		Preload("Order.TradePoint").
		Where("payments.trade_agent_id = ?", agent_id).
		Where("orders.trade_point_id = ?", tradePointId).
		Find(&payments).Error

	if err != nil {
		return []models.ReadPayment{}, err
	}

	var result []models.ReadPayment
	for _, item := range payments {
		data := models.ReadPayment{
			Id:            item.Id,
			OrderId:       item.OrderId,
			CreatedAt:     item.CreatedAt.Format("02.01.2006 15:04"),
			PaymentTypeId: item.PaymentTypeId,
			Currency:      item.Currency,
			Status:        item.Status,
			IsAgent:       item.IsAgent,
		}

		if data.PaymentTypeId == 4 {
			data.Currency = float64(item.Count)
		}

		if language == "ru" {
			data.TradePoint = item.Order.TradePoint.NameRu + " " + item.Order.TradePoint.OrientirRu + " " + item.Order.TradePoint.LocationRu
			data.PaymentType = item.PaymentType.NameRu
			data.TradeAgent = item.TradeAgent.NameRu
		} else {
			data.TradePoint = item.Order.TradePoint.NameTm + " " + item.Order.TradePoint.OrientirTm + " " + item.Order.TradePoint.LocationTm
			data.PaymentType = item.PaymentType.NameTm
			data.TradeAgent = item.TradeAgent.NameTm
		}

		result = append(result, data)
	}

	if result == nil {
		return []models.ReadPayment{}, nil
	}

	return result, nil
}

func (r *PaymentPostgres) GetDailyPaymentsInformations(agent_id, day, month, year int, language string) (models.ObjectWithIdNameAndCurrency, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	var order_value, self_debt, trade_point_all_debt float64

	result := models.ObjectWithIdNameAndCurrency{}
	var payments_list []models.Payment

	var trade_agent_self_debt models.Debt
	var trade_points []models.TradePoint

	get_agent_debt_request := r.db.Model(&models.Debt{}).Where("trade_agent_id = ?", agent_id).Find(&trade_agent_self_debt)

	if get_agent_debt_request.Error != nil {
		self_debt = 0
	} else {
		self_debt = trade_agent_self_debt.Total
	}

	var orders []models.Order

	get_orders_request := r.db.Model(models.Order{}).Where("DATE(created_at) = ? AND trade_agent_id = ?", date.Format("2006-01-02"), agent_id).Find(&orders)

	if get_orders_request.Error != nil {
		order_value = 0
	} else {
		for _, order_item := range orders {
			order_value += order_item.Sum
		}
	}

	get_trade_points_request := r.db.Model(models.TradePoint{}).Where("trade_agent_id = ?", agent_id).Find(&trade_points)

	if get_trade_points_request.Error != nil {
		trade_point_all_debt = 0
	} else {
		for _, point := range trade_points {
			point_debt := r.GetTradePointAllDebtSum(point.Id)

			if point_debt > 0 {
				trade_point_all_debt += point_debt
			}
		}
	}

	get_payments_request := r.db.Model(&models.Payment{}).Where("DATE(created_at) = ?", date.Format("2006-01-02"))

	if agent_id > 0 {
		get_payments_request = get_payments_request.Where("trade_agent_id = ?", agent_id)
	}

	get_payments_request.Find(&payments_list)

	for _, payment := range payments_list {
		if payment.PaymentTypeId == 1 {
			if payment.IsAgent == true {
				result.AgentCash += payment.Currency
			} else {
				result.EkspeditorCash += payment.Currency
			}
		}

		if payment.PaymentTypeId == 2 {
			result.AgentTransfer += payment.Currency
		}

		if payment.PaymentTypeId == 3 {
			result.AgentTerminal += payment.Currency
		}

		if payment.PaymentTypeId == 4 {
			result.EkspeditorLed += payment.Currency
		}
	}

	result.CashSum = result.AgentCash + result.EkspeditorCash
	result.SelfDebt = self_debt
	result.TradePointAllDebt = trade_point_all_debt
	result.OrdersValue = order_value

	return result, nil
}

func (r *PaymentPostgres) GetTradePointsDebtList(agentid int, language string) ([]models.TradePointMobile, error) {
	var tradePoints []models.TradePoint
	var result []models.TradePointMobile

	for dayOfWeek := 1; dayOfWeek <= 6; dayOfWeek++ {
		var first_list, route_list []models.TradePointMobile

		request := r.db.Model(models.TradePoint{}).Preload("City").Preload("District").Where("status = ?", true).
			Where("trade_agent_id = ?", agentid).Where("day_id", dayOfWeek).Find(&tradePoints)

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
			var note models.PointNote
			var current_trade_point_debt float64
			current_trade_point_debt = r.GetTradePointAllDebtSum(tradePoint.Id)

			if current_trade_point_debt == 0 {
				continue
			}

			data.Id = tradePoint.Id
			data.Code = tradePoint.Code
			data.Contact = contact.Name
			data.Number = contact.Number
			data.Status = false
			data.Note = ""
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

			first_list = append(first_list, data)
		}

		for _, route := range routes {
			for _, data := range first_list {
				if route.TradePointId == data.Id {
					route_list = append(route_list, data)
				}
			}
		}

		for _, data := range first_list {
			ishave := false
			for _, route := range routes {
				if route.TradePointId == data.Id {
					ishave = true
				}
			}

			if ishave == false {
				route_list = append(route_list, data)
			}
		}

		result = append(result, route_list...)
	}

	if result == nil {
		return []models.TradePointMobile{}, nil
	}

	return result, nil
}
