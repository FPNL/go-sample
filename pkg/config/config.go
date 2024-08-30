package config

import (
	"errors"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var ErrEnvNotFound = errors.New("environment variable not found")

type Ops func(config.Config)

func WithWatcher(key string, observer func(key string, value config.Value), errHandler func(err error)) Ops {
	return func(c config.Config) {
		errHandler(c.Watch(key, observer))
	}
}

func LoadFromFile(path string, bc any, ops ...Ops) error {
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)
	if ops == nil {
		defer c.Close()
	}

	if err := c.Load(); err != nil {
		return err
	}

	if err := c.Scan(bc); err != nil {
		return err
	}

	for _, op := range ops {
		op(c)
	}

	return nil
}

func LoadFromEnv(key string, format string, bc any, ops ...Ops) error {
	c := config.New(
		config.WithSource(
			NewSource(key, format),
		),
	)
	if ops == nil {
		defer c.Close()
	}

	if err := c.Load(); err != nil {
		return err
	}

	if len(os.Getenv(key)) > 0 {
		if err := c.Scan(bc); err != nil {
			return err
		}

		for _, op := range ops {
			op(c)
		}

		return nil
	}

	return ErrEnvNotFound
}
