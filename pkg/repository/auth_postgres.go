package repository

import (
	"fmt"
	"time"

	"github.com/Saparly-go-dev/tagma/models"

	"github.com/Saparly-go-dev/tagma"
	"gorm.io/gorm"
)

type AuthPostgres struct {
	db *gorm.DB
}

// NewAuthPostgres creates a new instance of AuthPostgres with GORM DB
func NewAuthPostgres(db *gorm.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateUser creates a new user in the database and returns the user's ID
func (r *AuthPostgres) CreateUser(user tagma.User) (int, error) {
	user.CreatedAt = time.Now() // Set the created time
	result := r.db.Create(&user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.Id, nil
}

// GetUser retrieves a user by username. Password verification belongs to the service layer.
func (r *AuthPostgres) GetUser(username string) (tagma.User, error) {
	var user tagma.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err

}

func (r *AuthPostgres) UpdatePassword(userID int, password string) error {
	return r.db.Model(&tagma.User{}).Where("id = ?", userID).Update("password", password).Error
}

func (r *AuthPostgres) GetUserByID(userID int) (tagma.User, error) {
	var user tagma.User
	err := r.db.First(&user, userID).Error
	return user, err
}

func (r *AuthPostgres) GetAgentId(Id int) (int, error) {
	var userAgent tagma.UserTradeAgent
	fmt.Println(Id)
	err := r.db.Where("user_id = ?", Id).First(&userAgent).Error

	if err != nil {
		return 0, err
	}

	return userAgent.TradeAgentId, nil

}

func (r *AuthPostgres) SaveRefreshToken(userid int, token, refresh_token string) {

}

func (r *AuthPostgres) GetEkspeditorId(Id int) (int, error) {
	var ekspeditor models.UserEkspeditor

	request := r.db.Model(&models.UserEkspeditor{}).Where("user_id = ?", Id).First(&ekspeditor)

	if request.Error != nil {
		return 0, request.Error
	}

	return ekspeditor.EkspeditorId, nil
}

func (r *AuthPostgres) GetTradeAgentIdFromEkspeditorId(Id int) (int, error) {
	var data models.AgentsEkspeditors

	get_agent_ekspeditors_request := r.db.Model(&models.AgentsEkspeditors{}).Where("ekspeditor_id = ?", Id).First(&data)

	if get_agent_ekspeditors_request.Error != nil {
		return 0, get_agent_ekspeditors_request.Error
	}

	return data.TradeAgentId, nil
}

func (r *AuthPostgres) GetEkspeditoryIdFromAgentId(Id int) (int, error) {
	var data models.AgentsEkspeditors
	get_agent_ekspeditors_request := r.db.Model(&models.AgentsEkspeditors{}).Where("trade_agent_id = ?", Id).First(&data)

	if get_agent_ekspeditors_request.Error != nil {
		return 0, get_agent_ekspeditors_request.Error
	}

	return data.EkspeditorId, nil
}
