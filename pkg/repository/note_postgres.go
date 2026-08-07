package repository

import (
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type NotePostgres struct {
	db *gorm.DB
}

var loc *time.Location

func init() {
	var err error
	loc, err = time.LoadLocation("Asia/Ashgabat")
	if err != nil {
		fmt.Println("Wagt zolagyny alyp bilmedi:", err)
		loc = time.UTC
	}
}

func NewNotePostgres(db *gorm.DB) *NotePostgres {
	return &NotePostgres{db: db}
}

func (r *NotePostgres) Create(note models.CreateNote) error {
	fmt.Println(loc)
	fmt.Println(time.Now().In(loc))
	fmt.Println(time.Now())
	newNote := models.Note{
		Description: note.Description,
		CreatedAt:   time.Now().In(loc),
	}

	request := r.db.Model(&models.Note{}).Create(&newNote)
	return request.Error
}

func (r *NotePostgres) GetAll(pageSize, pageNumber int) (*models.NotePage, error) {
	var notes []models.Note
	var list []models.ReadNote

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.Note{}).Order("id").Where("id > ?", 0)

	var count int64
	request.Count(&count)
	request = request.Offset(offset).Limit(pageSize).Find(&notes)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, note := range notes {
		var data models.ReadNote
		data.Id = note.Id
		data.Description = note.Description
		data.CreateAt = note.CreatedAt.Format("02.01.2006")

		list = append(list, data)
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if list == nil {
		list = []models.ReadNote{}
	}

	// Create the page object
	page := models.NotePage{
		Items:      list,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
	}

	return &page, nil
}

func (r *NotePostgres) Update(Id int, note models.CreateNote) error {
	UpdateData := map[string]interface{}{
		"description": note.Description,
	}

	request := r.db.Model(&models.Note{}).Where("id = ?", Id).Updates(UpdateData)
	return request.Error
}

func (r *NotePostgres) Delete(Id int) error {
	request := r.db.Delete(&models.Note{}, Id)
	return request.Error
}

func (r *NotePostgres) CreatePointNote(note models.CreatePointNote) error {
	newPointNote := models.PointNote{
		Description:  note.Description,
		TradePointId: note.TradePointId,
	}

	request := r.db.Model(&models.PointNote{}).Create(&newPointNote)
	return request.Error
}

func (r *NotePostgres) GetAllPointNote(pageSize, pageNumber, tradePointId int, code, language string) (*models.PointNotePage, error) {
	var notes []models.PointNote
	var items []models.ReadPointNote

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.PointNote{}).Joins("JOIN trade_points ON trade_points.id = point_notes.trade_point_id").
		Preload("TradePoint").Order("id desc")

	if tradePointId > 0 {
		request = request.Where("trade_point_id = ?", tradePointId)
	}

	if len(code) > 0 {
		request = request.Where("trade_points.code ILIKE ?", "%"+code+"%")
	}

	var count int64
	request.Count(&count)
	request = request.Offset(offset).Limit(pageSize).Find(&notes)
	if request.Error != nil {
		return nil, request.Error
	}

	for _, item := range notes {
		var data models.ReadPointNote
		var district models.District

		districtRequest := r.db.Model(&models.District{}).Preload("City").Where("id = ?", item.TradePoint.CityId).Find(&district)

		if districtRequest.Error != nil {

		}

		data.Id = item.Id
		data.CityId = district.City.Id
		data.DisctrictId = district.Id
		data.Description = item.Description
		data.TradePointCode = item.TradePoint.Code
		data.CreatedAt = item.CreatedAt.Format("02.01.2006")
		data.TradePointId = item.TradePointId

		if language == "ru" {
			data.TradePoint = item.TradePoint.NameRu
			data.City = district.City.NameRu
			data.Disctrict = district.NameRu
			data.Location = item.TradePoint.LocationRu
		} else {
			data.TradePoint = item.TradePoint.NameTm
			data.City = district.City.NameTm
			data.Disctrict = district.NameTm
			data.Location = item.TradePoint.LocationTm
		}

		items = append(items, data)
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if items == nil {
		items = []models.ReadPointNote{}
	}

	// Create the page object
	page := models.PointNotePage{
		Items:      items,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
		Total:      int(count),
	}
	return &page, nil
}

func (r *NotePostgres) UpdatePointNote(Id int, pointNote models.CreatePointNote) error {
	updateData := map[string]interface{}{
		"description":    pointNote.Description,
		"trade_point_id": pointNote.TradePointId,
	}

	request := r.db.Model(&models.PointNote{}).Where("id = ?", Id).Updates(updateData)
	return request.Error
}

func (r *NotePostgres) DeletePointNote(Id int) error {
	request := r.db.Delete(&models.PointNote{}, Id)
	return request.Error
}
