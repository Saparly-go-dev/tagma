package repository

import (
	"errors"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type PromotionPostgres struct {
	db *gorm.DB
}

func NewPromotionPostgres(db *gorm.DB) *PromotionPostgres {
	return &PromotionPostgres{db: db}
}

func (r *PromotionPostgres) Create(promotion models.CreatePromotion) error {
	newPromotion := models.Promotion{
		NameRu:    promotion.NameRu,
		NameTm:    promotion.NameTm,
		Start:     promotion.Start,
		End:       promotion.End,
		CreatedAt: time.Now(),
		Count:     promotion.Count,
		Status:    true,
	}

	create_promotion_request := r.db.Model(&models.Promotion{}).Create(&newPromotion)
	if create_promotion_request.Error != nil {
		return create_promotion_request.Error
	}
	for _, item := range promotion.List {
		newPromotion_Product := models.PromotionProducts{
			ProductId:   item,
			PromotionId: newPromotion.Id,
		}
		save_promotion_product_request := r.db.Create(&newPromotion_Product)

		if save_promotion_product_request.Error != nil {
		}
	}
	return nil
}

func (r *PromotionPostgres) ChangeStatus(Id int) error {
	var promotion models.Promotion

	get_promotion_request := r.db.First(&promotion, Id)
	if get_promotion_request.Error != nil {
		return get_promotion_request.Error
	}

	if promotion.Status == true {
		updateData := map[string]interface{}{
			"status": false,
		}
		result := r.db.Model(&models.Promotion{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status": true,
		}
		result := r.db.Model(&models.Promotion{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *PromotionPostgres) GetAll(pageSize, pageNumber int, name, language string) (*models.PromotionPage, error) {
	var promotions []models.Promotion
	var promotion_list []models.ReadPromotion

	offset := (pageNumber - 1) * pageSize

	get_promotions_request := r.db.Model(&models.Promotion{}).Order("created_at desc").
		Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%").
		Order("id desc").
		Find(&promotions)

	var count int64
	get_promotions_request = get_promotions_request.Count(&count)

	get_promotions_request = get_promotions_request.Offset(offset).Limit(pageSize)

	if get_promotions_request.Error != nil {
	}

	for _, item := range promotions {
		var new_item models.ReadPromotion
		var promotion_elements []models.PromotionProducts
		var list []int

		get_promotion_elements := r.db.Model(&models.PromotionProducts{}).Where("promotion_id = ?", item.Id).Order("product_id").
			Find(&promotion_elements)

		if get_promotion_elements.Error != nil {
		}

		for _, element := range promotion_elements {
			list = append(list, element.ProductId)
		}

		new_item.Id = item.Id
		new_item.NameRu = item.NameRu
		new_item.NameTm = item.NameTm
		new_item.Start = item.Start.Format("02.01.2006")
		new_item.End = item.End.Format("02.01.2006")
		new_item.List = list
		new_item.Count = item.Count
		new_item.Status = item.Status
		new_item.CreatedAt = item.CreatedAt.Format("02.01.2006")

		promotion_list = append(promotion_list, new_item)
	}

	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if promotion_list == nil {
		promotion_list = []models.ReadPromotion{}
	}

	page := models.PromotionPage{
		Items:      promotion_list,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
	}

	return &page, nil
}

func (r *PromotionPostgres) Update(id int, data models.CreatePromotion) error {
	var promotion models.Promotion

	get_promotion_promotion := r.db.Model(&models.Promotion{}).Where("id = ?", id).First(&promotion)

	if get_promotion_promotion.Error != nil {
		return get_promotion_promotion.Error
	}

	delete_promotion_products_request := r.db.Model(&models.PromotionProducts{}).Where("promotion_id = ?", id).Delete(&models.PromotionProducts{})

	if delete_promotion_products_request.Error != nil {
		return delete_promotion_products_request.Error
	}

	updateData := map[string]interface{}{
		"name_ru": data.NameRu,
		"name_tm": data.NameTm,
		"start":   data.Start,
		"end":     data.End,
		"count":   data.Count,
	}

	update_promotion_request := r.db.Model(&models.Promotion{}).Where("id = ?", id).Updates(updateData)

	if update_promotion_request.Error != nil {
		return update_promotion_request.Error
	}

	for _, item := range data.List {
		newPromotion_Product := models.PromotionProducts{
			ProductId:   item,
			PromotionId: promotion.Id,
		}

		save_promotion_product_request := r.db.Create(&newPromotion_Product)

		if save_promotion_product_request.Error != nil {
		}
	}

	return nil
}

func (r *PromotionPostgres) UpdateNameAndTime(id int, data models.UpdateDiscountInfo) error {
	var count int64

	check_is_name_unique := r.db.Model(&models.Promotion{}).Where("name_ru = ? OR name_tm = ? ", data.NameRu, data.NameTm).Where("id != ?", id).Count(&count)
	if count > 0 || check_is_name_unique.Error != nil {
		return errors.New("Rugsat berilmeýär")
	}

	check_is_discount_have := r.db.Model(&models.Promotion{}).Where("id = ?", id).Count(&count)

	if check_is_discount_have.Error != nil || count == 0 {
		return errors.New("Rugsat berilmeýär")
	}

	update_request := r.db.Model(&models.Promotion{}).Where("id = ?", id).Updates(data)
	if update_request.Error != nil {
		return errors.New("Ýalňyşlyk ýüze çykdy")
	}

	return nil
}

func (r *PromotionPostgres) Delete(id int) error {

	delete_promotion_products_request := r.db.Model(&models.PromotionProducts{}).Where("promotion_id = ?", id).Delete(&models.PromotionProducts{})

	if delete_promotion_products_request.Error != nil {
		return delete_promotion_products_request.Error
	}

	delete_promotion_request := r.db.Model(&models.Promotion{}).Where("id = ?", id).Delete(&models.Promotion{})

	return delete_promotion_request.Error

}
