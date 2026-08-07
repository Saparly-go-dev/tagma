package repository

import (
	"github.com/Saparly-go-dev/tagma"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type CityPostgres struct {
	db *gorm.DB
}

// NewCityPostgres creates a new instance of CityPostgres with GORM DB
func NewCityPostgres(db *gorm.DB) *CityPostgres {
	return &CityPostgres{db: db}
}

// Create inserts a new city into the database and returns its ID
func (r *CityPostgres) Create(city tagma.CreateCity) (int, error) {

	var lastdata tagma.City
	var codeindex int

	lastresult := r.db.Model(tagma.City{}).Last(&lastdata).Order("id desc")

	if lastresult.Error != nil {
		codeindex = 1
	} else {
		codeindex = lastdata.Id + 1
	}

	var code = strconv.Itoa(codeindex)

	if len(code) < 2 {
		code = "0" + code
	}

	newCity := tagma.City{
		NameRu:    city.NameRu,
		NameTm:    city.NameTm,
		Code:      code,
		Status:    city.Status,
		CreatedAt: time.Now(),
	}

	result := r.db.Create(&newCity)
	if result.Error != nil {
		return 0, result.Error
	}

	return newCity.Id, nil
}

// GetAll retrieves all cities from the database
func (r *CityPostgres) GetAll(pageSize, pageNumber int, name string) (*tagma.CityPage, error) {
	var cities []tagma.City // Remove unnecessary pointer
	var resultdata []tagma.CityRead
	offset := (pageNumber - 1) * pageSize

	result := r.db.Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Offset(offset).Limit(pageSize).Order("id").Find(&cities)
	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate total page count
	var count int64
	result2 := r.db.Model(&tagma.City{}).Where("name_ru ILIKE ?", "%"+name+"%").Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)
	if result2.Error != nil {
		return nil, result2.Error
	}
	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, city := range cities {
		var data tagma.CityRead

		data.NameRu = city.NameRu
		data.NameTm = city.NameTm
		data.Status = city.Status
		data.Id = city.Id
		data.Code = city.Code
		data.UpdatedAt = city.CreatedAt.Format("02.01.2006")

		resultdata = append(resultdata, data)
	}

	if resultdata == nil {
		resultdata = []tagma.CityRead{}
	}

	var page = tagma.CityPage{
		Items:      resultdata,
		PageCount:  int(pageCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	//Creating page with calculate variables

	return &page, nil
}

func (r *CityPostgres) ChangeStatus(cityId int) error {
	var city tagma.City
	result := r.db.First(&city, cityId)
	if result.Error != nil {
		return result.Error
	}

	if city.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"created_at": time.Now(),
		}
		result := r.db.Model(&tagma.City{}).Where("id = ?", cityId).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"created_at": time.Now(),
		}
		result := r.db.Model(&tagma.City{}).Where("id = ?", cityId).Updates(updateData)
		return result.Error
	}

}

// GetById retrieves a city by its ID
func (r *CityPostgres) GetById(cityId int) (tagma.City, error) {
	var city tagma.City
	result := r.db.First(&city, cityId)
	return city, result.Error
}

// Delete removes a city from the database by its ID
func (r *CityPostgres) Delete(cityId int) error {
	result := r.db.Delete(&tagma.City{}, cityId)
	return result.Error
}

// Update modifies a city's information based on its ID
func (r *CityPostgres) Update(cityId int, city tagma.CreateCity) error {
	updateData := map[string]interface{}{
		"name_ru": city.NameRu,
		"name_tm": city.NameTm,
		"status":  city.Status,
	}
	result := r.db.Model(&tagma.City{}).Where("id = ?", cityId).Updates(updateData)
	return result.Error
}
