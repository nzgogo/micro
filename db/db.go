package db

import (
	"time"

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
	opts Options
}

var (
	DefaultCharset         = "utf8"
	DefaultParseTime       = true
	DefaultLoc             = "Local"
	DefaultMaxIdleConns    = 4
	DefaultMaxOpenConns    = 16
	DefaultMaxConnLifetime = time.Hour * 2
)

func (d db) Connect() error {
	dsn := d.opts.Username + ":" + d.opts.Password + "@/" + d.opts.DBName
	dsn += "?charset=" + d.opts.Charset

	conn, err := gorm.Open("mysql", dsn)

	if err != nil {
		return err
	}

	d.conn = conn
	return nil
}

func NewDB(u, p, name string, opts ...Option) *db {
	options := Options{
		Username:        u,
		Password:        p,
		DBName:          name,
		Charset:         DefaultCharset,
		ParseTime:       DefaultParseTime,
		Loc:             DefaultLoc,
		MaxIdleConns:    DefaultMaxIdleConns,
		MaxOpenConns:    DefaultMaxOpenConns,
		MaxConnLifetime: DefaultMaxConnLifetime,
	}

	for _, o := range opts {
		o(&options)
	}

	return &db{
		opts: options,
	}
}


