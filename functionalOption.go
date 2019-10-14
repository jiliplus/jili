package main

import (
	"fmt"
	"time"
)

func main() {
	addr := "1270.0.0.1"
	newTimeout := 10 * time.Second
	// Options must be provided only if needed.
	Connect(addr)
	Connect(addr, WithTimeout(newTimeout))
	Connect(addr, WithCaching(false))
	Connect(
		addr,
		WithCaching(false),
		WithTimeout(newTimeout),
	)
}

// Options is DB's options
type Options struct {
	timeout time.Duration
	caching bool
}

// ApplyOption overrides behavior of Connect.
type ApplyOption func(*Options)

// Option overrides behavior of Connect.
// type Option interface {
// 	apply(*Options)
// }

// type optionFunc func(*Options)

// func (f optionFunc) apply(o *Options) {
// f(o)
// }

// WithTimeout if you need
func WithTimeout(t time.Duration) ApplyOption {
	return func(o *Options) {
		o.timeout = t
		fmt.Printf("With Timeout: %s\n", t)
	}
}

// WithCaching if you need
func WithCaching(cache bool) ApplyOption {
	return func(o *Options) {
		o.caching = cache
		fmt.Printf("With Caching: %t\n", cache)
	}
}

const (
	defaultTimeout = time.Second
	defaultCaching = true
)

// Connect creates a connection.
func Connect(
	addr string,
	opts ...ApplyOption,
) (*Options, error) {
	fmt.Println("---== new Connect ==---")
	options := Options{
		timeout: defaultTimeout,
		caching: defaultCaching,
	}

	for _, apply := range opts {
		apply(&options)
	}

	return &options, nil
}
