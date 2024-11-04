package worker

import (
	"fmt"
	"sync"
)

// Worker представляет собой воркер, который будет обрабатывать входящие данные.
type Worker struct {
	ID   int
	Work chan string
	Stop chan bool
	wg   *sync.WaitGroup
}

// NewWorker создает нового воркера.
func NewWorker(id int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:   id,
		Work: make(chan string),
		Stop: make(chan bool),
		wg:   wg,
	}
}

// Start запускает воркера для обработки данных из канала.
func (w *Worker) Start() {
	defer w.wg.Done()
	for {
		select {
		case data := <-w.Work:
			fmt.Printf("Worker %d processing data: %s\n", w.ID, data)
		case <-w.Stop:
			fmt.Printf("Worker %d stopping.\n", w.ID)
			return
		}
	}
}

