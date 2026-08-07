package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/handler"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
	"github.com/Saparly-go-dev/tagma/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//	@title			Tagma API
//	@version		1.0
//	@Description	API Server for Tagma

// host 0.0.0.0:8000
//	@BasePath	/

// @securityDefinitions.apikey ApiKeyAuth
// @in                      header
// @name                 Authorization

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Initialize the configuration
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %v", err.Error())
	}

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		logrus.Warnf("error loading .env variables: %s", err.Error())
	}
	if len(os.Getenv("JWT_SIGNING_KEY")) < 32 {
		logrus.Fatal("JWT_SIGNING_KEY must contain at least 32 characters")
	}

	// Initialize GORM database connection
	dsn := "host=" + configValue("DB_HOST", "db.host") +
		" user=" + configValue("DB_USERNAME", "db.username") +
		" dbname=" + configValue("DB_NAME", "db.dbname") +
		" port=" + configValue("DB_PORT", "db.port") +
		" sslmode=" + configValue("DB_SSLMODE", "db.sslmode")
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		dsn += " password=" + password
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Run auto migration for models if necessary
	err = db.AutoMigrate(&tagma.User{}, &tagma.UserTradeAgent{}, models.AgentsDrivers{}, &models.AgentsEkspeditors{}, &tagma.City{},
		&tagma.Trade_category{}, &tagma.Test{}, &models.ChannelType{}, &models.Discount{}, &models.DProduct{},
		&models.ChannelStructure{}, &models.ChannelSize{}, &models.ChannelManagement{}, &models.TradeChannel{}, &models.Product{},
		&models.ProductType{}, &models.Volume{}, &models.Driver{}, &models.Ekspeditor{}, &models.TradeAgent{}, &models.Post{},
		&models.Kind{}, &models.Upr{}, &models.TradePoint{}, &models.Contact{}, &models.ContactUprs{}, &models.Day{},
		&models.Route{}, &models.Image{}, &models.Log{}, &models.Order{}, &models.OrderList{}, &models.Gift{},
		&models.GiftType{}, &models.District{}, &models.Brand{}, &models.Marketing{}, &models.Note{}, &models.PointNote{},
		&models.EkspData{}, &models.UserEkspeditor{}, &models.DiscountChannelType{}, &models.PaymentType{}, &models.Payment{},
		&models.Manual{}, &models.ManualImage{}, &models.Dealer{}, &models.DealerOrder{}, &models.DealerOrderList{}, &models.RouteOrder{}, &models.Debt{},
		&models.RefreshToken{}, &models.Promotion{}, &models.PromotionProducts{}, &models.ExtraTradePointDay{}, models.Plan{},
		&models.Furniture{}, &models.TradePointsFurniture{}, &models.Merchandiser{}, &models.UserMerchandiser{}, &models.Promo{}, &models.PromoProduct{},
		&models.PromoChannelType{}, &models.DiscountTradeAgent{}, &models.DiscountTradePoint{})

	if err != nil {
		logrus.Fatalf("error during migration: %s", err.Error())
	}

	// Initialize repository, services, and handlers
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(tagma.Server)

	go func() {
		if err := srv.Run(configValue("PORT", "port"), handlers.InitRoutes()); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("The app Started")

	location, err := time.LoadLocation(configValue("TZ", "timezone"))
	if err != nil {
		logrus.Fatalf("invalid timezone: %s", err.Error())
	}
	c := cron.New(
		cron.WithLocation(location),
		cron.WithChain(
			cron.Recover(cron.DefaultLogger),
			cron.DelayIfStillRunning(cron.DefaultLogger),
		),
	)

	// Job 1: Run every day at 08:00 Tokyo time
	//_, _ = c.AddFunc("0 8 * * *", func() {
	//	log.Println("Running morning report job")
	//})

	_, _ = c.AddFunc("@every 86400s", func() {
		repos.TimeChecker.Orders()
		log.Println("Heartbeat check at", time.Now().Format("15:04:05"))
	})

	//_, _ = c.AddFunc("@every 10s", func() {
	//	repos.TimeChecker.Orders()
	//	log.Println("Heartbeat check at", time.Now().Format("15:04:05"))
	//})

	c.Start()
	defer c.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Printf("The app Shutting Down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("error occurred on server shutting down : %s", err.Error())
	}

	// Closing the DB connection, GORM doesn't have a direct Close() method so we use a manual DB close
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Errorf("error occurred when getting db: %s", err.Error())
	} else {
		if err := sqlDB.Close(); err != nil {
			logrus.Errorf("error occurred on db connection close: %s", err.Error())
		}
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func configValue(envKey, configKey string) string {
	if value := os.Getenv(envKey); value != "" {
		return value
	}
	return viper.GetString(configKey)
}

//func paginate(value interface{}, pagination *pkg.Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
//	var totalRows int64
//	db.Model(value).Count(&totalRows)
//
//	pagination.TotalRows = totalRows
//	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
//	pagination.TotalPages = totalPages
//
//	return func(db *gorm.DB) *gorm.DB {
//		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
//	}
//}
