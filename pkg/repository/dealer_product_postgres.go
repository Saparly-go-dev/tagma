package repository

import (
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type DealerProductPostgres struct {
	db *gorm.DB
}

func NewDealerProductPostgres(db *gorm.DB) *DealerProductPostgres {
	return &DealerProductPostgres{db: db}
}

func (r *DealerProductPostgres) Create(product models.CreateDealerProduct) (int, error) {

	var lastDealerProduct models.DealerProduct
	var codeindex int
	var productPrice float64

	lastresult := r.db.Model(models.DealerProduct{}).Last(&lastDealerProduct).Order("id desc")

	if lastresult.Error != nil {
		codeindex = 1
	} else {
		codeindex = lastDealerProduct.Id + 1
	}

	productPrice = float64(product.Price1) + float64(product.Price2)/100

	newDealerProduct := models.DealerProduct{
		BrandId:       product.BrandId,
		TasteRu:       product.TasteRu,
		TasteTm:       product.TasteTm,
		Price:         productPrice,
		Code:          strconv.Itoa(codeindex),
		Status:        product.Status,
		ProductTypeId: product.ProductTypeId,
		UpdatedAt:     time.Now(),
	}

	result := r.db.Create(&newDealerProduct)
	if result.Error != nil {
		return 0, result.Error
	}

	return newDealerProduct.Id, nil
}

// GetAll retrieves all cities from the database
func (r *DealerProductPostgres) GetAll(pageSize, pageNumber int, name, language string) (*models.DealerProductPage, error) {
	var products []models.DealerProduct
	var readDealerProducts []models.ReadDealerProduct

	// Calculate offset for pagination
	offset := (pageNumber - 1) * pageSize

	// Use correct join clause for product and product type, assuming a foreign key like `product_type_id` existsam
	result := r.db.Model(&models.DealerProduct{}).Preload("DealerProductType").Preload("Name").Where("taste_ru ILIKE ? OR taste_tm ILIKE ?", "%"+name+"%", "%"+name+"%").Order("id desc").
		Offset(offset).Limit(pageSize).Find(&products)

	if result.Error != nil {
		return nil, result.Error
	}

	// Iterate through the retrieved products
	for _, product := range products {
		var newitem models.ReadDealerProduct

		// Fill newitem with product details
		newitem.Id = product.Id
		newitem.Code = product.Code
		newitem.Name = product.Name.Name
		newitem.Price = product.Price
		newitem.Status = product.Status
		newitem.TasteRu = product.TasteRu
		newitem.TasteTm = product.TasteTm
		newitem.BrandId = product.BrandId
		newitem.ProductTypeId = product.ProductTypeId
		newitem.BlockPrice = roundToTwoDecimals((float64(product.ProductType.Count) * product.Price))
		newitem.Volume = product.ProductType.Volume
		newitem.Updatedat = product.UpdatedAt.Format("02.01.2006")

		// Set DealerProductType based on language
		if language == "ru" {
			newitem.ProductType = product.ProductType.NameRu
		} else {
			newitem.ProductType = product.ProductType.NameTm
		}
		// Append to the result list
		readDealerProducts = append(readDealerProducts, newitem)
	}

	// Calculate total count for pagination
	var count int64
	result2 := r.db.Model(&models.DealerProduct{}).
		Where("taste_ru ILIKE ? OR taste_tm ILIKE ?", "%"+name+"%", "%"+name+"%").
		Count(&count)

	if result2.Error != nil {
		return nil, result2.Error
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if readDealerProducts == nil {
		readDealerProducts = []models.ReadDealerProduct{}
	}

	// Create the page object
	page := models.DealerProductPage{
		Items:      readDealerProducts,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
	}

	return &page, nil
}

func (r *DealerProductPostgres) ChangeStatus(Id int) error {
	var data models.DealerProduct
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.DealerProduct{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.DealerProduct{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *DealerProductPostgres) Delete(productId int) error {
	result := r.db.Delete(&models.DealerProduct{}, productId)
	return result.Error
}

// Update modifies a city's information based on its ID
func (r *DealerProductPostgres) Update(DealerProductId int, DealerProduct models.CreateDealerProduct) error {
	var productPrice float64

	productPrice = float64(DealerProduct.Price1) + float64(DealerProduct.Price2)/100

	updateData := map[string]interface{}{
		"brand_id":        DealerProduct.BrandId,
		"status":          DealerProduct.Status,
		"taste_ru":        DealerProduct.TasteRu,
		"taste_tm":        DealerProduct.TasteTm,
		"price":           productPrice,
		"product_type_id": DealerProduct.ProductTypeId,
	}
	result := r.db.Model(&models.DealerProduct{}).Where("id = ?", DealerProductId).Updates(updateData)
	return result.Error
}
