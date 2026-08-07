package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type DiscountPostgres struct {
	db *gorm.DB
}

func NewDiscountPostgres(db *gorm.DB) *DiscountPostgres {
	return &DiscountPostgres{db: db}
}

func (r *DiscountPostgres) Create(discount models.CreateDiscount) error {
	newdiscount := models.Discount{
		NameRu:     discount.NameRu,
		NameTm:     discount.NameTm,
		Start:      discount.Start,
		End:        discount.End,
		SortOrder:  discount.SortOrder,
		Status:     discount.Status,
		IsProcent:  discount.IsProcent,
		IsDirect:   discount.IsDirect,
		IsIncrease: discount.IsIncrease,
		IsGroup:    discount.IsGroup,
		CreatedAt:  time.Now(),
	}

	request := r.db.Model(&models.Discount{}).Create(&newdiscount)
	if request.Error != nil {
		return request.Error
	}

	var secondCoefficient int

	for _, item := range discount.First {
		coefficient := item.Coefficient
		if discount.IsDirect == false {
			coefficient = 0
			secondCoefficient = item.Coefficient
		}
		newDProduct := models.DProduct{
			ProductId:   item.ProductId,
			DiscountId:  newdiscount.Id,
			Count:       item.Count,
			Coefficient: coefficient,
			Type:        true,
		}

		saveRequest := r.db.Model(&models.DProduct{}).Create(&newDProduct)
		if saveRequest.Error != nil {
			fmt.Println("CreateDProduct", saveRequest.Error)
			//return saveRequest.Error
		}
	}

	if discount.IsDirect == false {
		for _, item := range discount.Second {
			newDProduct := models.DProduct{
				ProductId:   item.ProductId,
				DiscountId:  newdiscount.Id,
				Count:       item.Count,
				Coefficient: secondCoefficient,
				Type:        false,
			}

			saveRequest := r.db.Model(&models.DProduct{}).Create(&newDProduct)
			if saveRequest.Error != nil {
				fmt.Println("CreateDProduct", saveRequest.Error)
				//return saveRequest.Error
			}
		}
	}

	for _, element := range discount.ChannelTypeIds {
		newDiscountChannelType := models.DiscountChannelType{
			ChannelTypeId: element,
			DiscountId:    newdiscount.Id,
		}

		saveRequest := r.db.Model(&models.DiscountChannelType{}).Create(&newDiscountChannelType)
		if saveRequest.Error != nil {
			fmt.Println("Skidka we ChannelType birikdirilmedi", element, " ", saveRequest.Error)
		}
	}

	for _, element := range discount.TradePointIds {
		newDiscountTradePoint := models.DiscountTradePoint{
			DiscountId:   newdiscount.Id,
			TradePointId: element,
		}

		saveDiscountTradePoint := r.db.Model(&models.DiscountTradePoint{}).Create(&newDiscountTradePoint)
		if saveDiscountTradePoint.Error != nil {
			fmt.Println("Yalnyshlyk yuze cykdy. Discount bilen TradePoint birikdirilmedi:", element)
		}
	}

	for _, element := range discount.TradeAgents {
		newDiscountTradeAgent := models.DiscountTradeAgent{
			DiscountId:   newdiscount.Id,
			TradeAgentId: element,
		}

		saveDiscountTradeAgent := r.db.Model(&models.DiscountTradeAgent{}).Create(&newDiscountTradeAgent)

		if saveDiscountTradeAgent.Error != nil {
			fmt.Println("Yalnyshlyk yuze cykdy. Discount bilen TradeAgent birikdirilmedi:", element)
		}
	}

	return request.Error
}

func (r *DiscountPostgres) GetAll(pageSize, pageNumber int, name, language string) (*models.DiscountPage, error) {
	var discounts []models.Discount
	var items []models.ReadDiscount

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.Discount{}).Preload("ChannelTypes").Preload("TradePoints").
		Preload("TradeAgents").Order("created_at desc").Where("id > ?", 1)

	if len(name) > 0 {
		request = request.Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	var count int64
	request = request.Count(&count)
	request = request.Offset(offset).Limit(pageSize).Find(&discounts)

	for _, discount := range discounts {
		var item models.ReadDiscount
		var dpProducts []models.DProduct
		var first []models.ReadDProduct
		var second []models.ReadDProduct
		var channelTypes, tradePoints, tradeAgents []models.SelectObject

		dpRequest := r.db.Model(&models.DProduct{}).Where("discount_id = ?", discount.Id).Where("type = ?", true).Find(&dpProducts)
		if dpRequest.Error != nil {
			fmt.Println("Dproductlar tapylmady", dpRequest.Error)
		}

		for _, item := range dpProducts {
			var element models.ReadDProduct

			element.ProductId = item.ProductId
			element.Count = item.Count
			element.Coefficient = item.Coefficient

			first = append(first, element)
		}

		if discount.IsDirect == false {
			secondDpRequest := r.db.Model(&models.DProduct{}).Where("discount_id = ?", discount.Id).Where("type = ?", false).Find(&dpProducts)
			if secondDpRequest.Error != nil {
				fmt.Println("second Dproductlar tapylmady", dpRequest.Error)
			}

			for _, item := range dpProducts {
				var element models.ReadDProduct

				element.ProductId = item.ProductId
				element.Count = item.Count
				element.Coefficient = item.Coefficient

				second = append(second, element)
			}
		}

		for _, tip := range discount.ChannelTypes {
			var temporary models.SelectObject
			temporary.Id = tip.Id
			if language == "ru" {
				temporary.Name = tip.NameRu
			} else {
				temporary.Name = tip.NameTm
			}
			channelTypes = append(channelTypes, temporary)
		}

		for _, tip := range discount.TradeAgents {
			var temporary models.SelectObject
			temporary.Id = tip.Id
			if language == "ru" {
				temporary.Name = tip.NameRu
			} else {
				temporary.Name = tip.NameTm
			}
			tradeAgents = append(tradeAgents, temporary)
		}

		for _, tip := range discount.TradePoints {
			var temporary models.SelectObject
			temporary.Id = tip.Id
			temporary.Name = tip.Code
			tradePoints = append(tradePoints, temporary)
		}

		if len(channelTypes) == 0 {
			channelTypes = []models.SelectObject{}
		}

		if len(tradeAgents) == 0 {
			tradeAgents = []models.SelectObject{}
		}

		if len(tradePoints) == 0 {
			tradePoints = []models.SelectObject{}
		}

		item.Id = discount.Id
		item.NameRu = discount.NameRu
		item.NameTm = discount.NameTm
		item.Start = discount.Start.Format("02.01.2006")
		item.End = discount.End.Format("02.01.2006")
		item.SortOrder = discount.SortOrder
		item.IsProscent = discount.IsProcent
		item.IsDirect = discount.IsDirect
		item.IsIncrease = discount.IsIncrease
		item.IsGroup = discount.IsGroup
		item.First = first
		item.Second = second
		item.Status = discount.Status
		item.ChannelTypes = channelTypes
		item.TradePoints = tradePoints
		item.TradeAgent = tradeAgents

		items = append(items, item)
	}

	//Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if items == nil {
		items = []models.ReadDiscount{}
	}

	page := models.DiscountPage{
		Items:      items,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		PageCount:  pageCount,
		Total:      int(count),
	}
	return &page, nil
}

func (r *DiscountPostgres) ChangeStatus(Id int) error {
	var data models.Discount
	request := r.db.Model(&models.Discount{}).Where("id = ?", Id).First(&data)
	if request.Error != nil {
		return request.Error
	}
	status := !data.Status

	updateData := map[string]interface{}{
		"status": status,
	}

	result := r.db.Model(&models.Discount{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *DiscountPostgres) Delete(Id int) error {
	var count int64

	ishaveOrderList := r.db.Model(&models.OrderList{}).Where("discount_id = ?", Id).Count(&count)

	if ishaveOrderList.Error != nil || count > 0 {
		return errors.New("Rugsat berilmeýär")
	}

	deleteDiscountChannelTypes := r.db.Model(&models.DiscountChannelType{}).Where("discount_id = ?", Id).Delete(&models.DiscountChannelType{})

	if deleteDiscountChannelTypes.Error != nil {
		fmt.Println("DiscountChannelType", deleteDiscountChannelTypes.Error)
	}

	deleteTradePoints := r.db.Model(&models.DiscountTradePoint{}).Where("discount_id = ?", Id).Delete(&models.DiscountTradePoint{})
	if deleteTradePoints.Error != nil {
		fmt.Println("DiscountTradePoint pozup bolmady", deleteTradePoints.Error)
	}

	deleteTradeAgents := r.db.Model(&models.DiscountTradeAgent{}).Where("discount_id = ?", Id).Delete(&models.DiscountTradeAgent{})

	if deleteTradeAgents.Error != nil {
		fmt.Println("DiscountTradeAgent pozup bolmady", deleteTradeAgents.Error)
	}

	var data models.DProduct
	request := r.db.Model(&models.DProduct{}).Where("discount_id = ?", Id).Delete(&data)

	if request.Error != nil {

	}

	deleteRequest := r.db.Delete(&models.Discount{}, Id)
	return deleteRequest.Error

}

func (r *DiscountPostgres) UpdateNameAndDate(id int, data models.UpdateDiscountInfo) error {
	var count int64

	check_is_name_unique := r.db.Model(&models.Discount{}).Where("name_ru = ? OR name_tm = ?", data.NameRu, data.NameTm).Where("id != ?", id).Count(&count)
	if count > 0 || check_is_name_unique.Error != nil {
		return errors.New("Rugsat berilmeýär 1 ")
	}

	check_is_discount_have := r.db.Model(&models.Discount{}).Where("id = ?", id).Count(&count)

	if check_is_discount_have.Error != nil || count == 0 {
		return errors.New("Rugsat berilmeýär")
	}

	updateData := map[string]interface{}{
		"name_ru":    data.NameRu,
		"name_tm":    data.NameTm,
		"start":      data.Start,
		"end":        data.End,
		"sort_order": data.SortOrder,
	}

	update_request := r.db.Model(&models.Discount{}).Where("id = ?", id).Updates(updateData)
	if update_request.Error != nil {
		return errors.New("Ýalňyşlyk ýüze çykdy")
	}

	deleteDiscountChannelTypes := r.db.Model(&models.DiscountChannelType{}).Where("discount_id= ?", id).Delete(&models.DiscountChannelType{})

	if deleteDiscountChannelTypes.Error != nil {
		fmt.Println("DiscountChannelType", deleteDiscountChannelTypes.Error)
	}

	deleteTradePoints := r.db.Model(&models.DiscountTradePoint{}).Where("discount_id = ?", id).Delete(&models.DiscountTradePoint{})
	if deleteTradePoints.Error != nil {
		fmt.Println("DiscountTradePoint pozup bolmady", deleteTradePoints.Error)
	}

	deleteTradeAgents := r.db.Model(&models.DiscountTradeAgent{}).Where("discount_id = ?", id).Delete(&models.DiscountTradeAgent{})

	if deleteTradeAgents.Error != nil {
		fmt.Println("DiscountTradeAgent pozup bolmady", deleteTradeAgents.Error)
	}

	deleteDProducts := r.db.Model(&models.DProduct{}).Where("discount_id = ?", id)
	if deleteDProducts.Error != nil {
		fmt.Println("DProducts pozup bolmady", deleteDProducts.Error)
	}

	for _, element := range data.ChannelTypeIds {
		newDiscountChannelType := models.DiscountChannelType{
			ChannelTypeId: element,
			DiscountId:    id,
		}

		saveRequest := r.db.Model(&models.DiscountChannelType{}).Create(&newDiscountChannelType)
		if saveRequest.Error != nil {
			fmt.Println("Skidka we ChannelType birikdirilmedi", element, " ", saveRequest.Error)
		}
	}

	for _, element := range data.TradePointIds {
		newDiscountTradePoint := models.DiscountTradePoint{
			DiscountId:   id,
			TradePointId: element,
		}

		saveDiscountTradePoint := r.db.Model(&models.DiscountTradePoint{}).Create(&newDiscountTradePoint)
		if saveDiscountTradePoint.Error != nil {
			fmt.Println("Yalnyshlyk yuze cykdy. Discount bilen TradePoint birikdirilmedi:", element)
		}
	}

	for _, element := range data.TradeAgents {
		newDiscountTradeAgent := models.DiscountTradeAgent{
			DiscountId:   id,
			TradeAgentId: element,
		}

		saveDiscountTradeAgent := r.db.Model(&models.DiscountTradeAgent{}).Create(&newDiscountTradeAgent)

		if saveDiscountTradeAgent.Error != nil {
			fmt.Println("Yalnyshlyk yuze cykdy. Discount bilen TradeAgent birikdirilmedi:", element)
		}
	}

	return nil
}

func (r *DiscountPostgres) Update(Id int, discount models.CreateDiscount) error {
	updateData := map[string]interface{}{
		"name_ru":     discount.NameRu,
		"name_tm":     discount.NameTm,
		"start":       discount.Start,
		"end":         discount.End,
		"is_procent":  discount.IsProcent,
		"is_direct":   discount.IsDirect,
		"is_increase": discount.IsIncrease,
		"sort_order":  discount.SortOrder,
	}

	request := r.db.Model(&models.Discount{}).Where("id = ?", Id).Updates(updateData)

	//delete Old ChannelType relations

	deleteDiscountChannelTypes := r.db.Model(&models.DiscountChannelType{}).Where("discount_id= ?", Id).Delete(&models.DiscountChannelType{})

	if deleteDiscountChannelTypes.Error != nil {
		fmt.Println("DiscountChannelType", deleteDiscountChannelTypes.Error)
	}

	deleteTradePoints := r.db.Model(&models.DiscountTradePoint{}).Where("discount_id = ?", Id).Delete(&models.DiscountTradePoint{})
	if deleteTradePoints.Error != nil {
		fmt.Println("DiscountTradePoint pozup bolmady", deleteTradePoints.Error)
	}

	deleteTradeAgents := r.db.Model(&models.DiscountTradeAgent{}).Where("discount_id = ?", Id).Delete(&models.DiscountTradeAgent{})

	if deleteTradeAgents.Error != nil {
		fmt.Println("DiscountTradeAgent pozup bolmady", deleteTradeAgents.Error)
	}

	deleteDProducts := r.db.Model(&models.DProduct{}).Where("discount_id = ?", Id)
	if deleteDProducts.Error != nil {
		fmt.Println("DProducts pozup bolmady", deleteDProducts.Error)
	}

	for _, element := range discount.ChannelTypeIds {
		newDiscountChannelType := models.DiscountChannelType{
			ChannelTypeId: element,
			DiscountId:    Id,
		}

		saveRequest := r.db.Model(&models.DiscountChannelType{}).Create(&newDiscountChannelType)
		if saveRequest.Error != nil {
			fmt.Println("Skidka we ChannelType birikdirilmedi", element, " ", saveRequest.Error)
		}
	}

	for _, element := range discount.TradePointIds {
		newDiscountTradePoint := models.DiscountTradePoint{
			DiscountId:   Id,
			TradePointId: element,
		}

		saveDiscountTradePoint := r.db.Model(&models.DiscountTradePoint{}).Create(&newDiscountTradePoint)
		if saveDiscountTradePoint.Error != nil {
			fmt.Println("Yalnyshlyk yuze cykdy. Discount bilen TradePoint birikdirilmedi:", element)
		}
	}

	for _, element := range discount.TradeAgents {
		newDiscountTradeAgent := models.DiscountTradeAgent{
			DiscountId:   Id,
			TradeAgentId: element,
		}

		saveDiscountTradeAgent := r.db.Model(&models.DiscountTradeAgent{}).Create(&newDiscountTradeAgent)

		if saveDiscountTradeAgent.Error != nil {
			fmt.Println("Yalnyshlyk yuze cykdy. Discount bilen TradeAgent birikdirilmedi:", element)
		}
	}

	for _, item := range discount.First {
		newDProduct := models.DProduct{
			ProductId:   item.ProductId,
			DiscountId:  Id,
			Count:       item.Count,
			Coefficient: item.Coefficient,
			Type:        true,
		}

		saveRequest := r.db.Model(&models.DProduct{}).Create(&newDProduct)
		if saveRequest.Error != nil {
			fmt.Println("CreateDProduct", saveRequest.Error)
			//return saveRequest.Error
		}
	}

	if discount.IsDirect == false {
		for _, item := range discount.Second {
			newDProduct := models.DProduct{
				ProductId:   item.ProductId,
				DiscountId:  Id,
				Count:       item.Count,
				Coefficient: item.Coefficient,
				Type:        false,
			}

			saveRequest := r.db.Model(&models.DProduct{}).Create(&newDProduct)
			if saveRequest.Error != nil {
				fmt.Println("CreateDProduct", saveRequest.Error)
				//return saveRequest.Error
			}
		}
	}

	return request.Error
}

func (r *DiscountPostgres) GetTradePointByTradeAgent(tradeAgentId int) ([]models.SelectObject, error) {
	var result []models.SelectObject
	var tradePoints []models.TradePoint

	get_trade_points_request := r.db.Model(&models.TradePoint{}).Where("trade_agent_id = ?", tradeAgentId).Find(&tradePoints)

	if get_trade_points_request.Error != nil {
		result = []models.SelectObject{}
		return result, get_trade_points_request.Error
	}

	for _, element := range tradePoints {
		var temporary models.SelectObject
		temporary.Id = element.Id
		temporary.Name = element.Code

		result = append(result, temporary)
	}

	if len(result) == 0 {
		result = []models.SelectObject{}
		return result, nil
	}

	return result, nil
}
