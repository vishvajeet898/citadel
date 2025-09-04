package psql

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/Orange-Health/citadel/common/constants"
	"github.com/Orange-Health/citadel/conf"
)

type Database struct {
	DB *gorm.DB
}

var dbORM Database

func Initialize() {
	config := conf.GetConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Calcutta connect_timeout=30",
		config.GetString("database.host"),
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.database"),
		config.GetString("database.port"),
		config.GetString("database.sslmode"),
	)

	logMode := logger.Silent
	if gin.IsDebugging() {
		logMode = logger.Info
	}

	psqlDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		log.Fatalf("%s , error: '%+v'", constants.ERROR_WHILE_CONNECTING_TO_DATABASE, err)
	}

	dbORM.DB = psqlDb
	dbORM.SetConnectionPool()
}

func (d *Database) SetConnectionPool() {
	config := conf.GetConfig()
	sqlDB, err := d.DB.DB()
	if err != nil {
		log.Fatalf("%s , error: '%+v'", constants.ERROR_WHILE_CONNECTING_TO_DATABASE, err)
	}

	sqlDB.SetMaxIdleConns(config.GetInt("database.maxIdleConnection"))

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(config.GetInt("database.maxOpenConnection"))

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func GetDbInstance() *gorm.DB {
	return dbORM.DB
}
