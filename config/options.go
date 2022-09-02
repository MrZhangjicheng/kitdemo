package config

type Option func(*options)

func WithSources(sources ...Source) Option {
	return func(o *options) {
		o.sources = sources
	}

}

type options struct {
	sources []Source
}
