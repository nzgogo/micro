package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB interface {
	Options() Options
	Connect() error
	DB() *gorm.DB
}

type db struct {
	conn *gorm.DB
	opt  Options
}
