package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id   int
	work chan string
	stop chan bool
	wg   *sync.WaitGroup
}

func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:   id,
		work: make(chan string),
		stop: make(chan bool),
		wg:   wg,
	}
}

func (w *Worker) Start() {
	defer w.wg.Done()
	for {
		select {
		case data := <-w.work:
			fmt.Printf("Worker %d processing data: %s\n", w.id, data)
		case <-w.stop:
			fmt.Printf("Worker %d stopping.\n", w.id)
			return
		}
	}
}

type WorkerPool struct {
	workers []*Worker
	work    chan string
	wg      sync.WaitGroup
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		work: make(chan string),
	}
}

func (p *WorkerPool) Start(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, &p.wg)
		p.workers = append(p.workers, worker)
		p.wg.Add(1)
		go worker.Start()
	}
}

func (p *WorkerPool) Submit(data string) {
	p.work <- data
}

func (p *WorkerPool) Run() {
	go func() {
		for data := range p.work {
			for _, worker := range p.workers {
				select {
				case worker.work <- data:
					break
				default:
					continue
				}
			}
		}
	}()
}

func (p *WorkerPool) Stop() {
	for _, worker := range p.workers {
		worker.stop <- true
	}
	p.wg.Wait()
}

func (p *WorkerPool) AddWorker() {
	worker := NewWorker(len(p.workers)+1, &p.wg)
	p.workers = append(p.workers, worker)
	p.wg.Add(1)
	go worker.Start()
}

func (p *WorkerPool) RemoveWorker() {
	if len(p.workers) == 0 {
		fmt.Println("No workers to remove.")
		return
	}
	lastWorker := p.workers[len(p.workers)-1]
	lastWorker.stop <- true
	p.workers = p.workers[:len(p.workers)-1]
}

func main() {
	pool := NewWorkerPool()
	pool.Start(3)
	pool.Run()

	for i := 0; i < 10; i++ {
		pool.Submit(fmt.Sprintf("Data %d", i+1))
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)
	pool.Stop()
}
