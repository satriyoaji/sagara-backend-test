package database

import (
	"github.com/satriyoaji/sagara-backend-test/helper"
	"github.com/satriyoaji/sagara-backend-test/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect()  {

	connection, err := gorm.Open(mysql.Open(helper.GoDotEnvVariable("DB_USERNAME") + ":" + helper.GoDotEnvVariable("DB_PASSWORD") + "@tcp("+helper.GoDotEnvVariable("DB_HOST")+":" + helper.GoDotEnvVariable("DB_PORT") + ")/" + helper.GoDotEnvVariable("DB_NAME") + "?parseTime=True&loc=Local"), &gorm.Config{})
	helper.OutputPanicError(err)

	DB = connection

	connection.AutoMigrate(models.User{})
	connection.AutoMigrate(models.Product{})
}