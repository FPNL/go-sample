package axiom

import (
	"context"
	"errors"
	"os"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var ErrEnvNotFound = errors.New("environment variable not found")

func LoadFromFile(path string, bc any) error {
	c := config.New(
		config.WithSource(
			file.NewSource(path),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		return err
	}

	return c.Scan(bc)
}

func LoadFromEnv(key string, format string, bc any) error {
	c := config.New(
		config.WithSource(
			NewSource(key, format),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		return err
	}

	if len(os.Getenv(key)) > 0 {
		if err := c.Scan(bc); err != nil {
			return err
		}

		return nil
	}

	return ErrEnvNotFound
}

// kratos 有內建獲取 env 的方式,
// 可是我們公司一般用 env 不是單一個 key 對應單一個 value,
// 而是一個 key 裡頭的 value 就是 json, 包含所有的環境變數.
// 因此這邊加上 format 參數，讓使用者可以指定格式
type env struct {
	key    string
	format string
}

func NewSource(key string, format string) config.Source {
	return &env{key, format}
}

func (e *env) Load() (kv []*config.KeyValue, err error) {
	return e.load(os.Getenv(e.key)), nil
}

func (e *env) load(value string) []*config.KeyValue {
	var kv []*config.KeyValue
	if len(value) != 0 {
		value := config.KeyValue{
			Key:    e.key,
			Value:  []byte(value),
			Format: e.format,
		}
		kv = append(kv, &value)
	}
	return kv
}

func (e *env) Watch() (config.Watcher, error) {
	w, err := NewWatcher()
	if err != nil {
		return nil, err
	}
	return w, nil
}

var _ config.Watcher = (*watcher)(nil)

type watcher struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWatcher() (config.Watcher, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &watcher{ctx: ctx, cancel: cancel}, nil
}

// Next will be blocked until the Stop method is called
func (w *watcher) Next() ([]*config.KeyValue, error) {
	<-w.ctx.Done()
	return nil, w.ctx.Err()
}

func (w *watcher) Stop() error {
	w.cancel()
	return nil
}
