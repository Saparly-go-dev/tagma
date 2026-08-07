package repository

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type ContactPostgres struct {
	db *gorm.DB
}

func NewContactPostgres(db *gorm.DB) *ContactPostgres {
	return &ContactPostgres{db: db}
}

func (r *ContactPostgres) Create(contact models.CreateContact) (int, error) {
	newContact := models.Contact{
		NameRu:       contact.NameRu,
		NameTm:       contact.NameTm,
		Number:       contact.Number,
		KindId:       contact.KindId,
		PostId:       contact.PostId,
		TradePointId: contact.TradePointId,
	}

	result := r.db.Create(&newContact)

	if result.Error != nil {
		return 0, result.Error
	}

	for _, element := range contact.UprsIds {
		newdata := models.ContactUprs{
			ContactId: newContact.Id,
			UprId:     element,
		}

		result := r.db.Model(&models.ContactUprs{}).Create(newdata)

		if result.Error != nil {
			fmt.Println(result.Error)
			return 0, result.Error
		}
	}

	return newContact.Id, nil
}

func (r *ContactPostgres) GetAll(PointId int, language string) (*[]models.ReadContact, error) {
	var resultData []models.ReadContact
	var contacts []models.Contact
	result := r.db.Model(&models.Contact{}).Preload("Kind").Preload("Post").
		Preload("Uprs").Where("trade_point_id = ?", PointId).Find(&contacts)

	if result.Error != nil {
		return &resultData, result.Error
	}

	if len(contacts) == 0 {
		return &resultData, nil
	}

	for _, contact := range contacts {
		var data models.ReadContact

		data.Id = contact.Id

		data.Number = contact.Number
		data.KindId = contact.KindId
		data.PostId = contact.PostId
		if language == "ru" {
			data.Name = contact.NameRu
			data.Kind = contact.Kind.NameRu
			data.Post = contact.Post.NameRu

			for _, element := range contact.Uprs {
				data.Uprs = append(data.Uprs, element.NameRu)
			}
		} else {
			data.Name = contact.NameTm
			data.Kind = contact.Kind.NameTm
			data.Post = contact.Post.NameTm

			for _, element := range contact.Uprs {
				data.Uprs = append(data.Uprs, element.NameTm)
			}
		}

		resultData = append(resultData, data)
	}

	return &resultData, nil
}

func (r *ContactPostgres) GetById(Id int) (*models.Contact, error) {
	var contact models.Contact

	result := r.db.Model(&models.Contact{}).Preload("Kind").Preload("Post").
		Preload("Uprs").Where("id = ?", Id).Find(&contact)

	if result.Error != nil {
		return &contact, result.Error
	}

	return &contact, nil
}

func (r *ContactPostgres) Delete(Id int) error {
	deleteResult := r.db.Model(&models.ContactUprs{}).Where("contact_id = ?", Id).Delete(&models.ContactUprs{})

	if deleteResult.Error != nil {

		return deleteResult.Error
	}

	result := r.db.Model(&models.Contact{}).Delete(&models.Contact{}, Id)

	return result.Error
}

func (r *ContactPostgres) Update(Id int, contact models.CreateContact) error {

	deleteResult := r.db.Model(&models.ContactUprs{}).Where("contact_id = ?", Id).Delete(&models.ContactUprs{})

	if deleteResult.Error != nil {

		return deleteResult.Error
	}

	updateData := map[string]interface{}{
		"name_ru":        contact.NameRu,
		"name_tm":        contact.NameTm,
		"number":         contact.Number,
		"kind_id":        contact.KindId,
		"post_id":        contact.PostId,
		"trade_point_id": contact.TradePointId,
	}
	result := r.db.Model(&models.Contact{}).Where("id = ?", Id).Updates(updateData)

	for _, element := range contact.UprsIds {
		newdata := models.ContactUprs{
			ContactId: Id,
			UprId:     element,
		}

		result := r.db.Model(&models.ContactUprs{}).Create(newdata)

		if result.Error != nil {
			//fmt.Println(result.Error)
		}
	}

	return result.Error
}

func (r *ContactPostgres) GetKinds(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.Kind

	result := r.db.Model(&models.Kind{}).Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range listOfData {
		var item models.SelectObject

		item.Id = element.Id
		if language == "ru" {
			item.Name = element.NameRu
		} else {
			item.Name = element.NameTm
		}

		responsedata = append(responsedata, item)
	}
	return &responsedata, nil
}

func (r *ContactPostgres) GetPosts(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.Post

	result := r.db.Model(&models.Post{}).Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range listOfData {
		var item models.SelectObject

		item.Id = element.Id
		if language == "ru" {
			item.Name = element.NameRu
		} else {
			item.Name = element.NameTm
		}

		responsedata = append(responsedata, item)
	}
	return &responsedata, nil
}

func (r *ContactPostgres) GetUprs(name, language string) (*[]models.SelectObject, error) {
	var responsedata []models.SelectObject
	var listOfData []models.Upr

	result := r.db.Model(&models.Upr{}).Where("name_tm LIKE ?", "%"+name+"%").Or("name_ru LIKE ?", "%"+name+"%").Find(&listOfData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, element := range listOfData {
		var item models.SelectObject

		item.Id = element.Id
		if language == "ru" {
			item.Name = element.NameRu
		} else {
			item.Name = element.NameTm
		}

		responsedata = append(responsedata, item)
	}
	return &responsedata, nil
}
