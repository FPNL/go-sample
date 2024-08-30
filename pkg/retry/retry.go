package retry

import (
	"math"
	"math/rand"
	"time"
)

type retry struct {
	maxTimes      int
	execute       func(count int) error
	deathHandler  func(err error)
	delayStrategy func(i int) time.Duration
}

type RetryOptions func(retry)

// NewRetry
// maxTimes: 重试次数
// execute: 執行函数
// deathHandler: 重试次数用完后的处理函数
// 預設蟲是策略是退避指數
func NewRetry(maxTimes int, execute func(count int) error, deathHandler func(err error), options ...RetryOptions) retry {
	r := retry{
		maxTimes:     maxTimes,
		execute:      execute,
		deathHandler: deathHandler,
		delayStrategy: func(i int) time.Duration {
			return func(n int, maximumBackoff float64) time.Duration {
				basicNumber := 1 << n
				randomNumber := float64(basicNumber) + rand.Float64()

				waitTime := math.Min(randomNumber, maximumBackoff)

				return time.Duration(waitTime) * time.Second
			}(i, 60)
		},
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func (receiver retry) Do() {
	var err error

	for i := 0; i < receiver.maxTimes; i++ {
		if err = receiver.execute(i); err == nil {
			return
		}
		if receiver.delayStrategy != nil {
			time.Sleep(receiver.delayStrategy(i))
		}
	}

	receiver.deathHandler(err)
}

// WithDelayStrategy 重試策略，輸入參數代表當前重試至第幾次，輸出表示要等待多就才執行下一次重試
func WithDelayStrategy(fn func(int) time.Duration) RetryOptions {
	return func(r retry) {
		r.delayStrategy = fn
	}
}
