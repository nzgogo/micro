package gogo

type Service interface {
	Options() Options
	Init(...Options) error
	Run() error
	Close()
}

type service struct {
	opts    Options
	Name    string
	Version string
	ID      string
}

func (s *service) Options() Options {
	return s.opts
}

func (s *service) Init(opts ...Options) error {
	return nil
}

func (s *service) Run() error {
	return nil
}

func (s *service) Close() {
}

func NewService() *service {
	return &service{}
}
