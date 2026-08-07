package repository

import (
	"errors"
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ManualPostgres struct {
	db *gorm.DB
}

func NewManualPostgres(db *gorm.DB) *ManualPostgres {
	return &ManualPostgres{db: db}
}

func (r *ManualPostgres) Create(data models.CreateManual) error {
	newManual := models.Manual{
		TitleRu:   data.TitleRu,
		ContentRu: data.ContentRu,
		TitleTm:   data.TitleTm,
		ContentTm: data.ContentTm,
		CreatedAt: time.Now(),
	}

	saveRequest := r.db.Model(&newManual).Create(&newManual)
	if saveRequest.Error != nil {
		return saveRequest.Error
	}

	return nil
}

func (r *ManualPostgres) Get(language string) (*models.ReadManual, error) {
	var manual models.Manual

	getRequest := r.db.Model(&manual).Where("id > 0").First(&manual)
	if getRequest.Error != nil {
		return &models.ReadManual{}, nil
	}

	result := models.ReadManual{}

	if manual.Id == 0 {
		return nil, nil
	}

	result.Id = manual.Id
	result.TitleRu = manual.TitleRu
	result.ContentRu = manual.ContentRu
	result.TitleTm = manual.TitleTm
	result.ContentTm = manual.ContentTm

	result.CreatedAt = manual.CreatedAt.Format("02.01.2006")

	var images []models.ManualImage

	if manual.Id > 0 {
		imageRequest := r.db.Model(&models.ManualImage{}).Where("manual_id = ?", manual.Id).Find(&images)
		if imageRequest.Error != nil {
			images = []models.ManualImage{}
		}

		for idx, item := range images {
			images[idx].Name = "images/manual/" + item.Name
		}

		result.Images = images
	}

	return &result, nil
}

func (r *ManualPostgres) GetForMobile(language string) (*models.ReadManualForMobile, error) {
	var manual models.Manual

	getRequest := r.db.Model(&manual).Where("id > 0").First(&manual)
	if getRequest.Error != nil {
		return &models.ReadManualForMobile{}, nil
	}

	result := models.ReadManualForMobile{}

	if manual.Id == 0 {
		return nil, nil
	}

	result.Id = manual.Id
	if language == "ru" {
		result.Title = manual.TitleRu
		result.Content = manual.ContentRu
	} else {
		result.Title = manual.TitleTm
		result.Content = manual.ContentTm
	}

	var images []models.ManualImage
	var listOfString []string

	if manual.Id > 0 {
		imageRequest := r.db.Model(&models.ManualImage{}).Where("manual_id = ?", manual.Id).Find(&images)
		if imageRequest.Error != nil {
			listOfString = []string{}
		}

		for _, item := range images {
			image_path := "images/manual/" + item.Name
			listOfString = append(listOfString, image_path)
		}

	}
	result.Images = listOfString

	return &result, nil
}

func (r *ManualPostgres) UpdateManual(Id int, data models.CreateManual) error {
	updateData := map[string]interface{}{
		"title_ru":   data.TitleRu,
		"content_ru": data.ContentRu,
		"title_tm":   data.TitleTm,
		"content_tm": data.ContentTm,
		"created_at": time.Now(),
	}

	saveRequest := r.db.Model(&models.Manual{}).Where("id = ?", Id).Updates(updateData)
	if saveRequest.Error != nil {
		return saveRequest.Error
	}

	return nil
}

func (r *ManualPostgres) DeleteImage(Id int) error {
	var manualImage models.ManualImage

	request := r.db.Model(&manualImage).Where("id = ?", Id).Find(&manualImage)
	if request.Error != nil {
		return request.Error
	}

	outputDir := "images/manual/" + manualImage.Name

	err := deleteFile(outputDir)

	if err != nil {
		return err
	}

	deleteRequest := r.db.Model(&manualImage).Where("id = ?", Id).Delete(&manualImage)
	if deleteRequest.Error != nil {
		return deleteRequest.Error
	}

	return nil
}

func (r *ManualPostgres) SaveImage(Id int, file *multipart.FileHeader) error {

	fileExt := filepath.Ext(file.Filename)
	allowedExtensions := []string{".jpg", ".jpeg", ".png"}
	if !contains(allowedExtensions, fileExt) {
		return errors.New("Invalid file type. Please upload an image file.")
	}

	var imageName, _ = Manual_SaveImages(Id, *file)

	imageData := models.ManualImage{
		Name:     imageName,
		ManualId: Id,
	}

	saveRequest := r.db.Model(&imageData).Create(&imageData)
	if saveRequest.Error != nil {
		return saveRequest.Error
	}

	return nil
}

func Manual_SaveImages(Id int, image1 multipart.FileHeader) (string, error) {
	//outputDir := filepath.Join("images/manual", strconv.Itoa(Id))
	outputDir := "images/manual"
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return "*", fmt.Errorf("failed to create output directory: %w", err)
	}
	//fmt.Printf("Saving images to %s\n", outputDir)

	var firstImageName = RandomString(8) + ".webp"
	if err := manual_processAndSaveImage(image1, filepath.Join(outputDir, firstImageName), "png"); err != nil {
		log.Printf("Failed to process image1: %v", err)
	}

	//fmt.Println(" 12341234")

	return firstImageName, nil
}

func manual_processAndSaveImage(fileHeader multipart.FileHeader, outputPath string, format string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := manual_decodeImage(fileHeader)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	resizedImg := resize.Resize(450, 300, img, resize.Lanczos3)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	switch format {
	case "png":
		err = png.Encode(outFile, resizedImg)
	case "jpeg":
		err = jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 80})
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
	return err
}

func manual_decodeImage(fileHeader multipart.FileHeader) (image.Image, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var img image.Image
	var format string
	switch ext {
	case ".jpeg", ".jpg":
		img, err = jpeg.Decode(file)
		format = "jpeg"
	case ".png":
		img, err = png.Decode(file)
		format = "png"
	default:
		return nil, "", fmt.Errorf("unsupported image format: %s", ext)
	}

	return img, format, err
}
