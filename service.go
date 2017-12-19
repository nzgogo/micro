package gogo

type Service interface {
	Options() Options
	Init(...Options) error
	Run() error
	Close()
}

type service struct {
	opts    Options
	name    string
	version string
	id      string
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

func NewService(n string, v string, i string, opts ...Option) *service {
	return &service{
		opts:    newOptions(opts...),
		name:    n,
		version: v,
		id:      i,
	}
}
