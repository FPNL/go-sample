package batch

import "fmt"

func Batch[T int | string](batchSize int, data []T, fn func([]T, int, int) error) error {
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}

		batch := data[i:end] // 這就是你的批次
		if err := fn(batch, i, end); err != nil {
			return fmt.Errorf(
				"batch from index(%d~%d) fail, handle these(%v), remain these(%v). err: %w",
				i, end,
				batch, data[i:],
				err)
		}
	}

	return nil
}
