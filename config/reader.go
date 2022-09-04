package config

type Reader interface {
	Merge(...*KeyValue) error

	Resolver() error
}
