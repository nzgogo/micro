package selector

type Options struct {
	Strategy Strategy
}

// Option used to initialise the selector
type Option func(*Options)

// SetStrategy sets the default strategy for the selector
func SetStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}
