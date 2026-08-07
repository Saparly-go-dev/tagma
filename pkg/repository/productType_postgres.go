package repository

import (
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type ProductTypePostgres struct {
	db *gorm.DB
}

func NewProductTypePostgres(db *gorm.DB) *ProductTypePostgres {
	return &ProductTypePostgres{db: db}
}

// Create inserts a new city into the database and returns its ID
func (r *ProductTypePostgres) Create(producttype models.CreateProductType) (int, error) {

	var lastProduct models.ProductType
	var codeindex int

	lastresult := r.db.Model(models.ProductType{}).Last(&lastProduct).Order("id desc")

	if lastresult.Error != nil {
		codeindex = 1
	} else {
		codeindex = lastProduct.Id + 1
	}

	newProductType := models.ProductType{
		NameRu:    producttype.NameRu,
		NameTm:    producttype.NameTm,
		Code:      strconv.Itoa(codeindex),
		Volume:    producttype.Volume,
		Count:     producttype.Count,
		Status:    producttype.Status,
		UpdatedAt: time.Now(),
	}

	result := r.db.Create(&newProductType)
	if result.Error != nil {
		return 0, result.Error
	}

	return newProductType.Id, nil
}

// GetAll retrieves all cities from the database
func (r *ProductTypePostgres) GetAll(pageSize, pageNumber int, name string) (*models.ProductTypePage, error) {
	var productTypes []models.ProductType // Remove unnecessary pointer
	var pageData []models.ReadProductType

	offset := (pageNumber - 1) * pageSize

	result := r.db.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Offset(offset).Limit(pageSize).Order("id desc").Find(&productTypes)
	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate total page count
	var count int64
	result2 := r.db.Model(&models.ProductType{}).Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)
	if result2.Error != nil {
		return nil, result2.Error
	}
	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, productType := range productTypes {
		var data models.ReadProductType

		data.Id = productType.Id
		data.Code = productType.Code
		data.NameRu = productType.NameRu
		data.NameTm = productType.NameTm
		data.Count = productType.Count
		data.Volume = productType.Volume
		data.Status = productType.Status
		data.UpdatedAt = productType.UpdatedAt.Format("02.01.2006")

		pageData = append(pageData, data)
	}

	if pageData == nil {
		pageData = []models.ReadProductType{}
	}

	var page = models.ProductTypePage{
		Items:      pageData,
		PageCount:  int(pageCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	//Creating page with calculate variables

	return &page, nil
}

func (r *ProductTypePostgres) ChangeStatus(Id int) error {
	var data models.ProductType
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.ProductType{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.ProductType{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

// GetById retrieves a city by its ID
func (r *ProductTypePostgres) GetById(ProductTypeId int) (models.ProductType, error) {
	var productType models.ProductType
	result := r.db.First(&productType, ProductTypeId)
	return productType, result.Error
}

// Delete removes a city from the database by its ID
func (r *ProductTypePostgres) Delete(productTypeId int) error {
	result := r.db.Delete(&models.ProductType{}, productTypeId)
	return result.Error
}

// Update modifies a city's information based on its ID
func (r *ProductTypePostgres) Update(ProductTypeId int, ProductType models.CreateProductType) error {
	updateData := map[string]interface{}{
		"name_ru": ProductType.NameRu,
		"name_tm": ProductType.NameTm,
		"volume":  ProductType.Volume,
		"count":   ProductType.Count,
		"status":  ProductType.Status,
	}
	result := r.db.Model(&models.ProductType{}).Where("id = ?", ProductTypeId).Updates(updateData)
	return result.Error
}

func (r *ProductTypePostgres) GetAllVolumes() (*[]models.Volume, error) {
	var volumes []models.Volume

	result2 := r.db.Where("id > ?", 0).Order("id").Find(&volumes)

	if result2.Error != nil {
		return nil, result2.Error
	}

	return &volumes, nil
}
