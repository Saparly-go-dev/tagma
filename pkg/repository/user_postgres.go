package repository

import (
	"errors"
	"fmt"
	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"time"
)

type UserPostgres struct {
	db *gorm.DB
}

func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) create(user tagma.CreateUser) error {
	var count int64
	var ekspeditor models.Ekspeditor
	var agent models.TradeAgent
	var merchandiser models.Merchandiser
	var nameRu, nameTm string

	// Egen ulanyjynyň tipi ekspeditor ýa-da agent bolsa onda onuň adyny bazadan gözläp alýarys
	if user.EkspeditorId > 0 {
		ekspeditorRequest := r.db.Model(&ekspeditor).Where("id = ?", user.EkspeditorId).Find(&ekspeditor)

		if ekspeditorRequest.Error != nil {
			return ekspeditorRequest.Error
		}
		nameRu = ekspeditor.NameRu
		nameTm = ekspeditor.NameTm
	}

	if user.AgentId > 0 {
		agentRequest := r.db.Model(&models.TradeAgent{}).Where("id = ?", user.AgentId).Find(&agent)

		if agentRequest.Error != nil {
			return agentRequest.Error
		}

		nameRu = agent.NameRu
		nameTm = agent.NameTm
	}

	if user.MerchandiserId > 0 {
		merchandiserRequest := r.db.Model(&models.Merchandiser{}).Where("id = ?", user.MerchandiserId).Find(&merchandiser)
		if merchandiserRequest.Error != nil {
			return merchandiserRequest.Error
		}

		nameRu = merchandiser.NameRu
		nameTm = merchandiser.NameTm
	}

	if user.EkspeditorId == 0 && user.AgentId == 0 && user.MerchandiserId == 0 {
		nameRu = user.NameRu
		nameTm = user.NameTm
	}

	// Ulanyjynyň login-i ýete-täkligini barlaýarys
	validation := r.db.Model(&tagma.User{}).Where("username = ?", user.Username).Count(&count)

	if validation.Error != nil {
		return validation.Error
	}

	if count > 0 {
		return errors.New("user already exists")
	}

	newUser := tagma.User{
		NameRu:    nameRu,
		NameTm:    nameTm,
		Username:  user.Username,
		Password:  user.Password,
		Status:    true,
		CreatedAt: time.Now(),
		Type:      user.Type,
	}

	saveRequest := r.db.Model(&tagma.User{}).Create(&newUser)
	if saveRequest.Error != nil {
		return saveRequest.Error
	}

	if user.Type == "ekspeditor" {

		userEkspeditor := models.UserEkspeditor{
			UserId:       newUser.Id,
			EkspeditorId: user.EkspeditorId,
		}

		request := r.db.Model(&models.UserEkspeditor{}).Create(&userEkspeditor)
		if request.Error != nil {
			return request.Error
		}
	}

	if user.Type == "agent" {
		userAgent := tagma.UserTradeAgent{
			UserId:       newUser.Id,
			TradeAgentId: user.AgentId,
		}

		request := r.db.Model(&tagma.UserTradeAgent{}).Create(&userAgent)
		if request.Error != nil {
			return request.Error
		}
	}

	if user.Type == "merchandiser" {
		userMerchandiser := models.UserMerchandiser{
			UserId:         newUser.Id,
			MerchandiserId: user.MerchandiserId,
		}

		request := r.db.Model(&models.UserMerchandiser{}).Create(&userMerchandiser)
		if request.Error != nil {
			return request.Error
		}

	}

	return nil
}

func (r *UserPostgres) GetPage(pageSize, pageNumber int, name, tip, language string) (tagma.UserPage, error) {
	var users []tagma.User
	var result []tagma.ReadUser
	offset := (pageNumber - 1) * pageSize
	var count int64
	request := r.db.Model(&tagma.User{}).Where("id > ?", 0)

	if len(name) > 0 {
		request = request.Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	if len(tip) > 0 {
		request = request.Where("type = ?", tip)
	}

	if err := request.Count(&count).Error; err != nil {
		return tagma.UserPage{}, err
	}
	if err := request.Order("id").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return tagma.UserPage{}, err
	}

	for _, user := range users {
		var item tagma.ReadUser

		item.Id = user.Id
		item.NameRu = user.NameRu
		item.NameTm = user.NameTm
		item.Username = user.Username
		item.Status = user.Status
		item.CreatedAt = user.CreatedAt.Format("02.01.2006")
		item.Type = user.Type

		if user.Type == "ekspeditor" {
			var userEkspeditor models.UserEkspeditor

			request2 := r.db.Model(&models.UserEkspeditor{}).Where("user_id = ?", user.Id).Find(&userEkspeditor)
			if request2.Error != nil {
				fmt.Println("UserEkspeditor tapylmady", user.Id)
			}
			item.EkspeditorId = userEkspeditor.EkspeditorId

			if language == "ru" {
				item.TypeName = "Экспедитор"
			} else {
				item.TypeName = "Ekspeditor"
			}

		}

		if user.Type == "agent" {
			var userAgent tagma.UserTradeAgent

			request2 := r.db.Model(&tagma.UserTradeAgent{}).Where("user_id = ?", user.Id).Find(&userAgent)
			if request2.Error != nil {
				fmt.Println("UserAgent tapylmady", user.Id)
			}
			item.AgentId = userAgent.TradeAgentId
			if language == "ru" {
				item.TypeName = "Торговые представитель"
			} else {
				item.TypeName = "Söwda wekili"
			}
		}

		if user.Type == "merchandiser" {
			var userMerchandiser models.UserMerchandiser

			request2 := r.db.Model(&models.UserMerchandiser{}).Where("user_id = ?", user.Id).Find(&userMerchandiser)
			if request2.Error != nil {
				fmt.Println(request2.Error)
			}
			item.MerchandiserId = userMerchandiser.MerchandiserId
			if language == "ru" {
				item.TypeName = "Мерчендайзер"
			} else {
				item.TypeName = "Merçandiser"
			}
		}

		if user.Type == "viewer" {
			if language == "ru" {
				item.TypeName = "Зритель"
			} else {
				item.TypeName = "Tomaşaçy"
			}
		}

		if user.Type == "admin" {
			if language == "ru" {
				item.TypeName = "Администратор"
			} else {
				item.TypeName = "Administrator"
			}
		}

		result = append(result, item)
	}

	if result == nil {
		result = []tagma.ReadUser{}
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	page := tagma.UserPage{
		Items:      result,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	return page, nil
}

func (r *UserPostgres) ChangePassword(data tagma.Change_Password) error {
	updateData := map[string]interface{}{
		"password": data.Password,
	}

	if data.UserId <= 0 {
		return errors.New("Not Found")
	}

	result := r.db.Model(&tagma.User{}).Where("id = ?", data.UserId).Updates(updateData)
	return result.Error
}

func (r *UserPostgres) ChangeStatus(Id int) error {
	var user tagma.User
	result := r.db.Model(&tagma.User{}).Where("id = ?", Id).First(&user)

	if result.Error != nil {
		return result.Error
	}

	updateData := map[string]interface{}{
		"status":     !user.Status,
		"created_at": time.Now(),
	}

	request := r.db.Model(&tagma.User{}).Where("id = ?", user.Id).Updates(updateData)
	return request.Error
}

func (r *UserPostgres) Delete(Id int) error {
	var user tagma.User

	request := r.db.Model(&tagma.User{}).Where("id = ?", Id).Find(&user)
	if request.Error != nil || user.Id == 0 {
		return errors.New("Ulanyjy tapylmady")
	}

	if user.Type == "ekspeditor" {
		deleteRequest := r.db.Model(&models.UserEkspeditor{}).Where("user_id = ?", user.Id).Delete(&models.UserEkspeditor{})
		if deleteRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy:", deleteRequest.Error)
		}
	}

	if user.Type == "agent" {
		deleteRequest := r.db.Model(&tagma.UserTradeAgent{}).Where("user_id = ?", user.Id).Delete(&tagma.UserTradeAgent{})
		if deleteRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy:", deleteRequest.Error)
		}
	}

	if user.Type == "merchandiser" {
		deleteRequest := r.db.Model(&models.UserMerchandiser{}).Where("user_id = ?", user.Id).Delete(&models.UserMerchandiser{})
		if deleteRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy:", deleteRequest.Error)
		}
	}

	result := r.db.Model(&tagma.User{}).Delete(&tagma.User{}, Id)
	return result.Error
}

func (r *UserPostgres) Update(Id int, data tagma.CreateUser) error {
	var user tagma.User
	var ekspeditor models.Ekspeditor
	var agent models.TradeAgent
	var merchandiser models.Merchandiser
	var nameRu, nameTm string

	// Ulanyjynyň öňki relation-karyny barlaýarys we bazadan pozýarys we täze maglumatlary bazadan alýarys
	request := r.db.Model(&tagma.User{}).Where("id = ?", Id).Find(&user)
	if request.Error != nil || user.Id == 0 {
		return errors.New("Ulanyjy tapylmady")
	}

	// Üýtgedilýän loginiň ýeke-täkligini barlaýarys
	if user.Username != data.Username {
		var count int64
		isHaveRequest := r.db.Model(&tagma.User{}).Where("username = ?", data.Username).Count(&count)
		if isHaveRequest.Error != nil {

		}
		if count > 0 {
			return errors.New("Ulanyjy ady ulanylyar")
		}
	}

	if user.Type == "ekspeditor" {
		deleteRequest := r.db.Model(&models.UserEkspeditor{}).Where("user_id = ?", user.Id).Delete(&models.UserEkspeditor{})
		if deleteRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy:", deleteRequest.Error)
		}

		//Bazadan gerekli ekspeditoryň adyny tapýarys we täze nameRu, nameTm ululyklara dakýarys
		ekspeditorRequest := r.db.Model(&models.Ekspeditor{}).Where("id = ?", data.EkspeditorId).Find(&ekspeditor)

		if ekspeditorRequest.Error != nil {
			fmt.Println("ekspeditor tapylmady", ekspeditorRequest.Error)
		}
		nameRu = ekspeditor.NameRu
		nameTm = ekspeditor.NameTm
	}

	if user.Type == "agent" {
		deleteRequest := r.db.Model(&tagma.UserTradeAgent{}).Where("user_id = ?", user.Id).Delete(&tagma.UserTradeAgent{})
		if deleteRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy:", deleteRequest.Error)
		}

		//Bazadan gerekli agentiň adyny tapýarys we täze nameRu, nameTm ululyklara dakýarys
		agentRequest := r.db.Model(&models.TradeAgent{}).Where("id = ?", data.AgentId).Find(&agent)

		if agentRequest.Error != nil {
			fmt.Println("Agent tapylmady", agentRequest.Error)
		}

		nameRu = agent.NameRu
		nameTm = agent.NameTm
	}

	if user.Type == "merchandiser" {
		deleteRequest := r.db.Model(&models.UserMerchandiser{}).Where("user_id = ?", user.Id).Delete(&models.UserMerchandiser{})
		if deleteRequest.Error != nil {
			fmt.Println("yalnyshlyk yuze cykdy:", deleteRequest.Error)
		}

		//Bazadan gerekli agentiň adyny tapýarys we täze nameRu, nameTm ululyklara dakýarys
		merchandiserReqeust := r.db.Model(&models.Merchandiser{}).Where("id = ?", data.MerchandiserId).Find(&merchandiser)
		if merchandiserReqeust.Error != nil {
			fmt.Println("Merchandiser tapylmady", merchandiserReqeust.Error)
		}

		nameRu = merchandiser.NameRu
		nameTm = merchandiser.NameTm
	}

	if data.EkspeditorId == 0 && data.AgentId == 0 && data.MerchandiserId == 0 {
		nameRu = user.NameRu
		nameTm = user.NameTm
	}

	updateData := map[string]interface{}{
		"name_ru":  nameRu,
		"name_tm":  nameTm,
		"username": data.Username,
		"type":     data.Type,
	}

	result := r.db.Model(&tagma.User{}).Where("id = ?", user.Id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if data.Type == "ekspeditor" {

		userEkspeditor := models.UserEkspeditor{
			UserId:       user.Id,
			EkspeditorId: data.EkspeditorId,
		}

		saveRequest := r.db.Model(&models.UserEkspeditor{}).Create(&userEkspeditor)
		//fmt.Println("ekspeditruser yatda saklanyldy")
		if saveRequest.Error != nil {
			return saveRequest.Error
		}
	}

	if data.Type == "agent" {
		userAgent := tagma.UserTradeAgent{
			UserId:       user.Id,
			TradeAgentId: data.AgentId,
		}

		//fmt.Println("agentruser yatda saklanyldy")
		saveRequest := r.db.Model(&tagma.UserTradeAgent{}).Create(&userAgent)
		if saveRequest.Error != nil {
			return saveRequest.Error
		}
	}

	if data.Type == "merchandiser" {
		userMerchandiser := models.UserMerchandiser{
			UserId:         user.Id,
			MerchandiserId: data.MerchandiserId,
		}

		saveRequest := r.db.Model(&models.UserMerchandiser{}).Create(&userMerchandiser)
		if saveRequest.Error != nil {
			return saveRequest.Error
		}
	}

	return nil
}
