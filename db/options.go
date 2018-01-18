package db

import (
	"time"
)

type Options struct {
	Dialects		string //suppport mypath and yet more coming...
	Username        string
	Password        string
	Protocol		string
	Address			string
	DBName          string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	MaxConnLifetime time.Duration
}

type Option func(*Options)

func Dialects(dialects string) Option {
	return func(o *Options){
		o.Dialects = dialects
	}
}

func Protocol(p string) Option {
	return func(o *Options) {
		o.Protocol = p
	}
}

func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
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
