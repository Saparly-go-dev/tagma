package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type PromoPostgres struct {
	db *gorm.DB
}

func NewPromoPostgres(db *gorm.DB) *PromoPostgres {
	return &PromoPostgres{db: db}
}

func (r *PromoPostgres) Create(promo models.UpdatePromo) error {
	var count int64
	ishave_request := r.db.Model(&models.Promo{}).Where("name_ru = ? OR name_tm = ?", promo.NameRu, promo.NameTm).Count(&count)

	if ishave_request.Error != nil || count > 0 {
		return errors.New("Rugsat berilmeyar")
	}

	newPromo := models.Promo{
		NameRu:    promo.NameRu,
		NameTm:    promo.NameTm,
		Start:     promo.Start,
		End:       promo.End,
		NeedCount: promo.NeedCount,
		FreeCount: promo.FreeCount,
		CreatedAt: time.Now(),
		Status:    false,
		UpdatedAt: time.Now(),
	}

	create_request := r.db.Model(&models.Promo{}).Create(&newPromo)
	if create_request.Error != nil {
		return errors.New("Ýalnyşlyk ýüze çykdy")
	}

	for _, id := range promo.NeedIds {
		newPromoProduct := models.PromoProduct{
			PromoId:   newPromo.Id,
			ProductId: id,
			IsNeeded:  true,
		}

		save_request := r.db.Model(&models.PromoProduct{}).Create(&newPromoProduct)
		if save_request.Error != nil {
			return errors.New("Ýalnyşlyk ýüze çykdy")
		}
	}

	for _, id := range promo.FreeIds {
		newPromoProduct := models.PromoProduct{
			PromoId:   newPromo.Id,
			ProductId: id,
			IsNeeded:  false,
		}

		save_request := r.db.Model(&models.PromoProduct{}).Create(&newPromoProduct)
		if save_request.Error != nil {
			return errors.New("Ýalnyşlyk ýüze çykdy")
		}
	}

	for _, id := range promo.ChannelTypeIds {
		newPromoChannelType := models.PromoChannelType{
			PromoId:       newPromo.Id,
			ChannelTypeId: id,
		}

		save_request := r.db.Model(&models.PromoChannelType{}).Create(&newPromoChannelType)
		if save_request.Error != nil {
		}
	}

	return nil
}

func (r *PromoPostgres) Update(promo_id int, promo models.UpdatePromo) error {
	var db_promo models.Promo

	get_request := r.db.Model(&models.Promo{}).Where("id = ?", promo_id).First(&db_promo)
	if get_request.Error != nil {
		return errors.New("Not Found")
	}

	updateData := map[string]interface{}{
		"name_ru":    promo.NameRu,
		"name_tm":    promo.NameTm,
		"start":      promo.Start,
		"end":        promo.End,
		"need_count": promo.NeedCount,
		"free_count": promo.FreeCount,
		"updated_at": time.Now(),
	}

	update_request := r.db.Model(&models.Promo{}).Where("id = ?", promo_id).Updates(updateData)
	if update_request.Error != nil {
		return errors.New("Ýalnyşlyk ýüze çykdy")
	}

	delete_promo_products := r.db.Model(&models.PromoProduct{}).Where("promo_id = ?", promo_id).Delete(&models.PromoProduct{})

	if delete_promo_products.Error != nil {
		return errors.New("Ýalnyşlyk ýüze çykdy")
	}

	delete_promo_channelTypes := r.db.Model(&models.PromoChannelType{}).Where("promo_id = ?", promo_id).Delete(&models.PromoChannelType{})

	if delete_promo_channelTypes.Error != nil {
		return errors.New("Ýalnyşlyk ýüze çykdy")
	}

	for _, id := range promo.NeedIds {
		newPromoProduct := models.PromoProduct{
			PromoId:   promo_id,
			ProductId: id,
			IsNeeded:  true,
		}

		save_request := r.db.Model(&models.PromoProduct{}).Create(&newPromoProduct)
		if save_request.Error != nil {
			return errors.New("Ýalnyşlyk ýüze çykdy")
		}
	}

	for _, id := range promo.FreeIds {
		newPromoProduct := models.PromoProduct{
			PromoId:   promo_id,
			ProductId: id,
			IsNeeded:  false,
		}

		save_request := r.db.Model(&models.PromoProduct{}).Create(&newPromoProduct)
		if save_request.Error != nil {
			return errors.New("Ýalnyşlyk ýüze çykdy")
		}
	}

	for _, id := range promo.ChannelTypeIds {
		newPromoChannelType := models.PromoChannelType{
			PromoId:       promo_id,
			ChannelTypeId: id,
		}

		save_request := r.db.Model(&models.PromoChannelType{}).Create(&newPromoChannelType)
		if save_request.Error != nil {
		}
	}

	return nil
}

func (r *PromoPostgres) GetPage(pageSize, pageNumber int, name, language string) (models.PromoPage, error) {
	var db_promos []models.Promo
	items := []models.ReadPromo{}
	var count int64
	result := models.PromoPage{}

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.Promo{}).Preload("ChannelTypes").Order("created_at desc").
		Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%")

	request = request.Count(&count)
	request = request.Offset(offset).Limit(pageSize).Find(&db_promos)
	//fmt.Println(db_promos)

	for _, element := range db_promos {
		var read_promo models.ReadPromo
		var promoProducts []models.PromoProduct
		var channels []models.SelectObject
		var needIds, freeIds []int

		get_promoProducts := r.db.Model(&models.PromoProduct{}).Where("promo_id = ?", element.Id).Find(&promoProducts)
		if get_promoProducts.Error != nil {
			fmt.Println("promoProducts not Found")
		}

		for _, channel := range element.ChannelTypes {
			newdata := models.SelectObject{
				Id:   channel.Id,
				Name: channel.NameTm,
			}

			if language == "ru " {
				newdata.Name = channel.NameRu
			}

			channels = append(channels, newdata)
		}

		for _, promo_product := range promoProducts {
			if promo_product.IsNeeded {
				needIds = append(needIds, promo_product.ProductId)
			} else {
				freeIds = append(freeIds, promo_product.ProductId)
			}
		}

		read_promo.Id = element.Id
		read_promo.NameRu = element.NameRu
		read_promo.NameTm = element.NameTm
		read_promo.Start = element.Start
		read_promo.End = element.End
		read_promo.NeedCount = element.NeedCount
		read_promo.FreeCount = element.FreeCount
		read_promo.Status = element.Status
		read_promo.CreatedAt = element.CreatedAt
		read_promo.UpdatedAt = element.UpdatedAt
		read_promo.ChannelTypes = channels
		read_promo.NeedIds = needIds
		read_promo.FreeIds = freeIds

		items = append(items, read_promo)
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

func (r *PromoPostgres) ChangeStatus(id int) error {
	var data models.Promo
	request := r.db.Model(&models.Promo{}).Where("id = ?", id).First(&data)
	if request.Error != nil || data.Id == 0 {
		return errors.New("Not Found")
	}
	fmt.Println(data.Status)
	status := !data.Status

	fmt.Println(status)

	updateData := map[string]interface{}{
		"status":     !data.Status,
		"updated_at": time.Now(),
	}

	result := r.db.Model(&models.Promo{}).Where("id = ?", id).Updates(updateData)
	return result.Error
}

func (r *PromoPostgres) Delete(id int) error {
	var data models.Promo
	request := r.db.Model(&models.Promo{}).Where("id = ?", id).First(&data)
	if request.Error != nil || data.Id == 0 {
		return errors.New("Not Found")
	}

	delete_promo_products := r.db.Model(&models.PromoProduct{}).Where("promo_id = ?", id).Delete(&models.PromoProduct{})

	if delete_promo_products.Error != nil {
		return errors.New("Ýalnyşlyk ýüze çykdy")
	}

	delete_promo_channelTypes := r.db.Model(&models.PromoChannelType{}).Where("promo_id = ?", id).Delete(&models.PromoChannelType{})

	if delete_promo_channelTypes.Error != nil {
		return errors.New("Ýalnyşlyk ýüze çykdy")
	}

	delete_request := r.db.Model(&models.Promo{}).Delete(&data)

	return delete_request.Error
}
