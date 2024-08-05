package workpool

func NewWorkerPool[T any](numWorkers int, jobs <-chan T, process func(T)) {
	var workers = make(chan struct{}, numWorkers)

	go func() {
		for job := range jobs {
			workers <- struct{}{}
			go func(job T) {
				process(job)
				<-workers
			}(job)
		}
	}()
}
