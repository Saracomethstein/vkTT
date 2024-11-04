package pool

import (
	"sync"
	"fmt"
	"vkTT/internal/worker"
)

// WorkerPool представляет собой пул воркеров.
type WorkerPool struct {
	workers []*worker.Worker
	work    chan string
	wg      sync.WaitGroup
}

// NewWorkerPool создает новый пул воркеров.
func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		work: make(chan string),
	}
}

// Start запускает указанный пул воркеров.
func (p *WorkerPool) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		w := worker.NewWorker(i+1, &p.wg)
		p.workers = append(p.workers, w)
		p.wg.Add(1)
		go w.Start()
	}
}

// Submit добавляет работу в пул для обработки.
func (p *WorkerPool) Submit(data string) {
	p.work <- data
}

// Run запускает цикл для обработки входящих данных.
func (p *WorkerPool) Run() {
	go func() {
		for data := range p.work {
			for _, w := range p.workers {
				select {
				case w.Work <- data:
					break
				default:
					continue
				}
			}
		}
	}()
}

// Stop останавливает всех воркеров.
func (p *WorkerPool) Stop() {
	for _, w := range p.workers {
		w.Stop <- true
	}
	p.wg.Wait()
}

// AddWorker добавляет нового воркера в пул.
func (p *WorkerPool) AddWorker() {
	w := worker.NewWorker(len(p.workers)+1, &p.wg)
	p.workers = append(p.workers, w)
	p.wg.Add(1)
	go w.Start()
}

// RemoveWorker удаляет воркера из пула.
func (p *WorkerPool) RemoveWorker() {
	if len(p.workers) == 0 {
		fmt.Println("No workers to remove.")
		return
	}
	lastWorker := p.workers[len(p.workers)-1]
	lastWorker.Stop <- true
	p.workers = p.workers[:len(p.workers)-1]
}
