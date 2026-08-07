package repository

import (
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type DealerPostgres struct {
	db *gorm.DB
}

func NewDealerPostgres(db *gorm.DB) *DealerPostgres {
	return &DealerPostgres{db: db}
}

func (r *DealerPostgres) Create(dealer models.CreateDealer) error {
	newDealer := models.Dealer{
		NameRu:      dealer.NameRu,
		NameTm:      dealer.NameTm,
		AddressRu:   dealer.AddressRu,
		AddressTm:   dealer.AddressTm,
		Company:     dealer.Company,
		PhoneNumber: dealer.PhoneNumber,
		Email:       dealer.Email,
		CityId:      dealer.CityId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      dealer.Status,
		Account:     0,
	}

	result := r.db.Model(&models.Dealer{}).Create(&newDealer)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *DealerPostgres) GetPage(pageSize, pageNumber, cityId int, name, language string) (models.DealerPage, error) {
	var list []models.Dealer
	var items []models.ReadDealer
	offset := (pageNumber - 1) * pageSize

	result := r.db.Model(&models.Dealer{}).Preload("City").Where("id > 0")

	if cityId > 0 {
		result = result.Where("city_id = ?", cityId)
	}

	if len(name) > 0 {
		result = result.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_rm ILIKE ?", "%"+name+"%")
	}

	var count int64

	result = result.Count(&count)

	result = result.Order("id desc").Offset(offset).Limit(pageSize).Find(&list)

	if result.Error != nil {
		return models.DealerPage{}, result.Error
	}

	for _, element := range list {
		var data models.ReadDealer
		var dealer_orders_uppaid []models.DealerOrder
		get_dealer_not_paid_order_request := r.db.Model(&models.DealerOrder{}).Where("dealer_id = ?", element.Id).
			Where("is_closed = ?", false).Where("status = ?", true).Find(&dealer_orders_uppaid)

		if get_dealer_not_paid_order_request.Error != nil {
			fmt.Println(get_dealer_not_paid_order_request.Error)
		}

		var debt float64

		for _, item := range dealer_orders_uppaid {
			debt += item.Sum
		}
		data.Id = element.Id
		data.NameRu = element.NameRu
		data.NameTm = element.NameTm
		data.AddressRu = element.AddressRu
		data.AddressTm = element.AddressTm
		data.Company = element.Company
		data.PhoneNumber = element.PhoneNumber
		data.Email = element.Email
		data.CityId = element.CityId
		data.Status = element.Status
		data.Account = element.Account
		data.Debt = debt
		data.CreatedAt = element.CreatedAt.Format("02.01.2006")
		data.UpdatedAt = element.UpdatedAt.Format("02.01.2006 15:04")
		if language == "ru" {
			data.City = element.City.NameRu
		} else {
			data.City = element.City.NameTm
		}

		items = append(items, data)
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if items == nil {
		items = []models.ReadDealer{}
	}

	page := models.DealerPage{
		Items:      items,
		PageCount:  pageCount,
		PageSize:   pageSize,
		PageNumber: pageNumber,
		Total:      int(count),
	}

	return page, nil
}

func (r *DealerPostgres) Update(Id int, dealer models.CreateDealer) error {
	var data models.Dealer

	result := r.db.Model(&models.Dealer{}).Where("id = ?", Id).Updates(&dealer)

	if result.Error != nil || data.Id == 0 {
		return nil
	}

	updateDate := map[string]interface{}{
		"name_ru":      dealer.NameRu,
		"name_tm":      dealer.NameTm,
		"address_ru":   dealer.AddressRu,
		"address_tm":   dealer.AddressTm,
		"company":      dealer.Company,
		"phone_number": dealer.PhoneNumber,
		"email":        dealer.Email,
		"status":       dealer.Status,
	}

	updateRequest := r.db.Model(&models.Dealer{}).Where("id = ?", Id).Updates(updateDate)
	if updateRequest.Error != nil {
		return updateRequest.Error
	}

	return nil
}

func (r *DealerPostgres) ChangeStatus(Id int) error {
	var data models.Dealer
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"created_at": time.Now(),
		}
		result := r.db.Model(&models.Dealer{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"created_at": time.Now(),
		}
		result := r.db.Model(&models.Dealer{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *DealerPostgres) GetDealerForSelect(name, language string) ([]models.SelectObject, error) {
	var result []models.SelectObject
	var dealers_list []models.Dealer

	request := r.db.Model(&models.Dealer{}).Preload("City").Where("status = ?", true)

	if len(name) > 0 {
		request = request.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_rm ILIKE ?", "%"+name+"%")
	}

	request.Order("id desc").Find(&dealers_list)

	for _, element := range dealers_list {
		var data models.SelectObject

		data.Id = element.Id
		if language == "ru" {
			data.Name = element.NameRu
		} else {
			data.Name = element.NameTm
		}

		result = append(result, data)
	}

	if result == nil {
		result = []models.SelectObject{}
	}

	return result, nil
}

func (r *DealerPostgres) TopUpDealerAccount(dealerId, price1, price2 int) error {
	var dealer models.Dealer

	get_dealer_request := r.db.Model(&models.Dealer{}).Where("id = ?", dealerId).First(&dealer)
	if get_dealer_request.Error != nil {
		return get_dealer_request.Error
	}
	sum := float64(price2)/100 + float64(price1)
	sum += dealer.Account

	update_dealer := map[string]interface{}{
		"account":    Round2(sum),
		"updated_at": time.Now(),
	}

	dealer_update_request := r.db.Model(&models.Dealer{}).Where("id = ?", dealerId).Updates(update_dealer)
	if dealer_update_request.Error != nil {
		return dealer_update_request.Error
	}

	var dealer_orders []models.DealerOrder
	get_dealer_order_request := r.db.Model(&models.DealerOrder{}).Order("id").Where("dealer_id = ? AND is_closed = ? AND status = ?", dealerId, false, true).Find(&dealer_orders)
	fmt.Println(dealer_orders)
	if get_dealer_order_request.Error != nil {
		return get_dealer_order_request.Error
	}

	//Dileriň töllenmedik sargytlaryna töleg edip çykýarys
	account_balance := sum
	for _, item := range dealer_orders {
		if item.Sum <= account_balance {
			update_dealer_order_data := map[string]interface{}{
				"is_closed": true,
				"paid_date": time.Now(),
			}

			update_dealer_order_request := r.db.Model(&models.DealerOrder{}).Where("id = ?", item.Id).Updates(update_dealer_order_data)
			if update_dealer_order_request.Error != nil {
				return update_dealer_order_request.Error
			}
			account_balance = account_balance - item.Sum

			update_dealer_data := map[string]interface{}{
				"account":    Round2(account_balance),
				"updated_at": time.Now(),
			}

			update_dealer_request := r.db.Model(&models.Dealer{}).Where("id = ?", dealerId).Updates(update_dealer_data)
			if update_dealer_request.Error != nil {
				fmt.Println("yalnysh123")
				return update_dealer_request.Error
			}

		} else {
			break
		}
	}
	return nil
}
