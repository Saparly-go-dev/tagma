package repository

import (
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"time"
)

type DriverPostgres struct {
	db *gorm.DB
}

func NewDriverPostgres(db *gorm.DB) *DriverPostgres {
	return &DriverPostgres{db: db}
}

func (r *DriverPostgres) Create(driver models.CreateDriver) (int, error) {
	newDriver := models.Driver{
		NameRu:    driver.NameRu,
		NameTm:    driver.NameTm,
		Status:    driver.Status,
		Number:    driver.Number,
		UpdatedAt: time.Now(),
	}

	result := r.db.Create(&newDriver)
	if err := result.Error; err != nil {
		return 0, err
	}

	return newDriver.Id, nil
}

func (r *DriverPostgres) GetAll(pageNumber, pageSize int, name, language string) (*models.DriverPage, error) {
	var drivers []models.Driver
	var pagedata []models.ReadDriver

	offset := (pageNumber - 1) * pageSize

	result := r.db.Model(&models.Driver{}).Preload("TradeAgent").Where("name_ru ILIKE ?", "%"+name+"%").
		Or("name_tm ILIKE ?", "%"+name+"%").Order("id").Offset(offset).Limit(pageSize).Find(&drivers)

	if result.Error != nil {
		return nil, result.Error
	}

	var count int64
	result2 := r.db.Model(&models.Driver{}).Where("name_ru ILIKE ?", "%"+name+"%").
		Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)

	if result2.Error != nil {
		return nil, result.Error
	}

	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, driver := range drivers {
		var data models.ReadDriver
		var tradeAgent models.TradeAgent
		//fmt.Println(driver)
		if len(driver.TradeAgent) > 0 {
			tradeAgent = driver.TradeAgent[0]
		}

		data.Id = driver.Id
		data.Code = tradeAgent.Code
		data.NameRu = driver.NameRu
		data.NameTm = driver.NameTm
		data.Status = driver.Status
		data.Number = driver.Number
		data.UpdatedAt = driver.UpdatedAt.Format("02.01.2006")
		if language == "ru" {
			data.TradeAgent = tradeAgent.NameRu
		} else {
			data.TradeAgent = tradeAgent.NameTm
		}

		pagedata = append(pagedata, data)
	}

	var page = models.DriverPage{
		Items:      pagedata,
		PageCount:  int(pageCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	//Creating page with calculate variables

	return &page, nil
}
func (r *DriverPostgres) ChangeStatus(Id int) error {
	var driver models.Driver
	result := r.db.First(&driver, Id)

	if err := result.Error; err != nil {
		return err
	}

	// Toggle status
	newStatus := !driver.Status
	updateData := map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	// Use "=" in the Where clause instead of "=="
	result = r.db.Model(&models.Driver{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *DriverPostgres) GetById(Id int) (models.Driver, error) {
	var driver models.Driver
	result := r.db.First(&driver, Id)
	return driver, result.Error
}

func (r *DriverPostgres) Delete(Id int) error {
	result := r.db.Delete(&models.Driver{}, Id)
	return result.Error
}

func (r *DriverPostgres) Update(Id int, driver models.CreateDriver) error {
	updateData := map[string]interface{}{
		"name_ru": driver.NameRu,
		"name_tm": driver.NameTm,
		"status":  driver.Status,
		"number":  driver.Number,
	}

	result := r.db.Model(&models.Driver{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}
