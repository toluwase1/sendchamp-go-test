package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sendchamp-go-test/config"
	"sendchamp-go-test/models"
	"time"
)

type GormDB struct {
	DB *gorm.DB
}

func GetDB(c *config.Config) *GormDB {
	gormDB := &GormDB{}
	gormDB.Init(c)
	return gormDB
}

func (g *GormDB) Init(c *config.Config) {
	g.DB = getMySqlDB(c)

	if err := migrate(g.DB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getMySqlDB(c *config.Config) *gorm.DB {
	log.Printf("Connecting to mysql: %+v", c)
	dsn := "root:toluwase@tcp(127.0.0.1:3306)/sendchamp?charset=utf8mb4&parseTime=True&loc=Local"
	//postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
	//	c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level Info, Silent, Warn, Error
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	gormConfig := &gorm.Config{
		Logger: newLogger,
	}
	if c.Env == "prod" {
		gormConfig = &gorm.Config{}
	}
	postgresDB, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		log.Fatal(err)
	}
	return postgresDB
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.BlackList{}, &models.Task{})
	if err != nil {
		return fmt.Errorf("migrations error: %v", err)
	}

	return nil
}
