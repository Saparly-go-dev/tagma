package repository

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type LogPostgres struct {
	db *gorm.DB
}

func NewLogPostgres(db *gorm.DB) *LogPostgres {
	return &LogPostgres{db: db}
}

//func (r *LogPostgres) Save(infoRu, infoTm, logType string) (int, error) {
//	data := models.Log{
//		InfoRu:   infoRu,
//		InfoTm:   infoTm,
//		LogType:  logType,
//		User:     "admin",
//		CreateAt: time.Now(),
//	}
//
//	result := r.db.Model(&models.Log{}).Create(&data)
//
//	if result.Error != nil {
//		log.Println(result.Error)
//		return 0, result.Error
//	}
//
//	return data.Id, result.Error
//}

func (r *LogPostgres) GetAll(language string, pageSize, pageNumber int, logType string) (*models.LogPage, error) {
	var logs []models.Log
	var pageData []models.ReadLog

	offset := (pageNumber - 1) * pageSize
	fmt.Println(pageSize)

	result := r.db.Model(&models.Log{}).Where("log_type ILIKE ?", "%"+logType+"%").Order("id desc").Offset(offset).Limit(pageSize).Find(&logs)
	fmt.Println(result.Error)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, log := range logs {
		var data models.ReadLog

		data.Id = log.Id
		data.CreateAt = log.CreateAt.Format("02.01.2006")
		data.User = log.User
		if language == "ru" {
			data.Info = log.InfoRu
			data.LogType = logType_Ru[log.LogType]

		} else {
			data.Info = log.InfoTm
			data.LogType = logType_Tm[log.LogType]
		}

		pageData = append(pageData, data)
	}
	fmt.Println("mapper gecdi geldi")

	var count int64
	result2 := r.db.Model(&models.Log{}).Where("log_type ILIKE ?", "%"+logType+"%").Count(&count)

	if result2.Error != nil {
		return nil, result2.Error
	}

	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	//if pageData == nil {
	//	pageData = []models.ReadLog{}
	//}

	page := models.LogPage{
		Items:      pageData,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	return &page, nil
}

func (r *LogPostgres) GetLogTypes(language string) (*[]string, error) {
	var result []string

	if language == "ru" {
		for _, element := range logType_Ru {
			result = append(result, element)
		}
	} else {
		for _, element := range logType_Tm {
			result = append(result, element)
		}
	}

	return &result, nil
}

var (
	logType_Ru = map[string]string{
		"ChangeEkspeditor": "Смена экспедитора",
		"ChangeDriver":     "Смена водителя",
	}

	logType_Tm = map[string]string{
		"ChangeEkspeditor": "Ekspeditor çalyşmak",
		"ChangeDriver":     "Sürüjiniň çalyşmak",
	}
)
