package server

type Server interface {
	Options() Options
	Init(...Option) error
}

type Option func(*Options)
