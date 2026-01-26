package consumer

type WorkerPool struct {
	sem chan struct{}
}

func NewWorkerPool(size int) *WorkerPool {
	return &WorkerPool{sem: make(chan struct{}, size)}
}

func (p *WorkerPool) Submit(fn func()) {
	p.sem <- struct{}{}
	go func() {
		defer func() { <-p.sem }()
		fn()
	}()
}
