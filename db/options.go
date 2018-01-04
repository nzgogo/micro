package db

import "time"

type Options struct {
	Username        string
	Password        string
	DBName          string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxConnLifetime time.Duration
}

type Option func(*Options)

func Username(username string) Option {
	return func(o *Options) {
		o.Username = username
	}
}

func Password(pwd string) Option {
	return func(o *Options) {
		o.Password = pwd
	}
}

func DBName(db string) Option {
	return func(o *Options) {
		o.DBName = db
	}
}

func Charset(charset string) Option {
	return func(o *Options) {
		o.Charset = charset
	}
}

func ParseTime(parseTime bool) Option {
	return func(o *Options) {
		o.ParseTime = parseTime
	}
}

func Loc(loc string) Option {
	return func(o *Options) {
		o.Loc = loc
	}
}

func MaxIdleConns(num int) Option {
	return func(o *Options) {
		o.MaxIdleConns = num
	}
}

func MaxOpenConns(num int) Option {
	return func(o *Options) {
		o.MaxOpenConns = num
	}
}

func MaxConnLifetime(duration time.Duration) Option {
	return func(o *Options) {
		o.MaxConnLifetime = duration
	}
}
