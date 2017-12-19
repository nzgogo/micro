package gogo

type Service interface {
	Options() Options
	Init(...Options) error
	Run() error
	Close() error
}

type service struct {
	opts    Options
	Name    string
	Version string
	ID      string
}
