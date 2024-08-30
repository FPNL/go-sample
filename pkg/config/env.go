package config

import (
	"os"

	"github.com/go-kratos/kratos/v2/config"
)

// kratos 有內建獲取 env 的方式, 並加上 format 參數，讓使用者可以指定格式
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
