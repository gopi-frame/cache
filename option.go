package cache

type Option[T any] interface {
	Apply(c *Cache[T]) error
}

type OptionFunc[T any] func(c *Cache[T]) error

func (f OptionFunc[T]) Apply(c *Cache[T]) error {
	return f(c)
}

// WithEncoder sets encoder.
func WithEncoder[T any](encoder func(T) ([]byte, error)) OptionFunc[T] {
	return func(c *Cache[T]) error {
		if encoder == nil {
			return nil
		}
		c.encoder = encoder
		return nil
	}
}

// WithDecoder sets decoder.
func WithDecoder[T any](decoder func([]byte) (T, error)) OptionFunc[T] {
	return func(c *Cache[T]) error {
		if decoder == nil {
			return nil
		}
		c.decoder = decoder
		return nil
	}
}

// ForString sets encoder and decoder for string type.
func ForString() OptionFunc[string] {
	return func(c *Cache[string]) error {
		c.encoder = func(value string) ([]byte, error) {
			return []byte(value), nil
		}
		c.decoder = func(bs []byte) (string, error) {
			return string(bs), nil
		}
		return nil
	}
}
