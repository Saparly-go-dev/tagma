package repository

import (
	"github.com/Saparly-go-dev/tagma"
	"gorm.io/gorm"
	"time"
)

type TradeCategoryPostgres struct {
	db *gorm.DB
}

// NewTradeCategoryPostgres initializes a new TradeCategoryPostgres instance
func NewTradeCategoryPostgres(db *gorm.DB) *TradeCategoryPostgres {
	return &TradeCategoryPostgres{db: db}
}

// Create creates a new trade category in the database and returns the category's ID
func (r *TradeCategoryPostgres) Create(category tagma.CreateTrade_category) (int, error) {
	newCatgegory := tagma.Trade_category{
		NameRu:    category.Name,
		NameTm:    category.Name,
		Min_Sum:   category.Min_Sum,
		Max_Sum:   category.Max_Sum,
		Status:    category.Status,
		UpdatedAt: time.Now(),
	}

	result := r.db.Create(&newCatgegory)
	if result.Error != nil {
		return 0, result.Error
	}

	return newCatgegory.Id, nil
}

// GetAll retrieves all trade categories from the database
func (r *TradeCategoryPostgres) GetAll(pageNumber, pageSize int, name string) (*tagma.TradeCategoryPage, error) {
	var categories []tagma.Trade_category
	var pagedata []tagma.ReadTrade_category

	offset := (pageNumber - 1) * pageSize

	result := r.db.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Offset(offset).Limit(pageSize).Order("id").Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate total page count
	var count int64
	result2 := r.db.Model(&tagma.Trade_category{}).Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)
	if result2.Error != nil {
		return nil, result2.Error
	}
	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, category := range categories {
		var data tagma.ReadTrade_category

		data.Id = category.Id
		data.NameRu = category.NameRu
		data.NameTm = category.NameTm
		data.Min_Sum = category.Min_Sum
		data.Max_Sum = category.Max_Sum
		data.Status = category.Status
		data.UpdatedAt = category.UpdatedAt.Format("02.01.2006")

		pagedata = append(pagedata, data)
	}

	var page = tagma.TradeCategoryPage{
		Items:      pagedata,
		PageCount:  int(pageCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	//Creating page with calculate variables

	return &page, nil
}

func (r *TradeCategoryPostgres) ChangeStatus(Id int) error {
	var data tagma.Trade_category
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&tagma.Trade_category{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&tagma.Trade_category{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

// GetById retrieves a trade category by its ID
func (r *TradeCategoryPostgres) GetById(Id int) (tagma.Trade_category, error) {
	var category tagma.Trade_category
	err := r.db.First(&category, Id).Error
	return category, err
}

// Delete removes a trade category by its ID
func (r *TradeCategoryPostgres) Delete(Id int) error {
	return r.db.Delete(&tagma.Trade_category{}, Id).Error
}

// Update updates the fields of a specific trade category
func (r *TradeCategoryPostgres) Update(Id int, category tagma.CreateTrade_category) error {
	updateData := map[string]interface{}{
		"name_ru": category.Name,
		"name_tm": category.Name,
		"status":  category.Status,
		"Min_Sum": category.Min_Sum,
		"Max_Sum": category.Max_Sum,
	}

	request := r.db.Model(&tagma.Trade_category{}).Where("id = ?", Id).Updates(updateData)
	return request.Error
}
