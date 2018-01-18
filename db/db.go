package db

import (
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type DB interface {
	Options() Options
	Connect() error
	Close() error
	DB() *gorm.DB
}

type db struct {
	conn *gorm.DB
	opts Options
}

var (
	DefaultDialect         = "mysql"
	DefaultProtocol        = "tcp"
	DefaultAddress         = "workbench.cugybz6qn13l.ap-southeast-2.rds.amazonaws.com"
	DefaultCharset         = "utf8"
	DefaultParseTime       = true
	DefaultLoc             = "Local"
	DefaultMaxIdleConns    = 4
	DefaultMaxOpenConns    = 16
	DefaultMaxConnLifetime = time.Hour * 2
)

func (d *db) Connect() error {
	// The Data Source Name has a common format, like e.g. PEAR DB uses it, but without type-prefix (optional parts marked by squared brackets):
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := d.opts.Username
	dsn += ":" + d.opts.Password + "@"
	dsn += d.opts.Protocol
	dsn += "(" + d.opts.Address + ")"
	dsn += "/" + d.opts.DBName
	dsn += "?charset=" + d.opts.Charset
	dsn += "&parseTime=" + strconv.FormatBool(d.opts.ParseTime)

	conn, err := gorm.Open(d.opts.Dialects, dsn)

	if err != nil {
		return err
	}

	conn.DB().SetMaxIdleConns(d.opts.MaxIdleConns)
	conn.DB().SetMaxOpenConns(d.opts.MaxOpenConns)
	conn.DB().SetConnMaxLifetime(d.opts.MaxConnLifetime)
	d.conn = conn
	return nil
}

func (d *db) Options() Options {
	return d.opts
}

func (d *db) Close() error {
	return d.conn.Close()
}

func (d *db) DB() *gorm.DB {
	return d.conn
}

func NewDB(u, p, name string, opts ...Option) *db {
	options := Options{
		Dialects:        DefaultDialect,
		Username:        u,
		Password:        p,
		Protocol:        DefaultProtocol,
		Address:         DefaultAddress,
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
