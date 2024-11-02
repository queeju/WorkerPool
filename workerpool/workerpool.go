package workerpool

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errCannotRemoveWorker = errors.New("no more workers left to remove")
	errPoolStopped        = errors.New("worker pool stopped")
)

// worker представляет воркера, который обрабатывает задания из канала jobChan
type worker struct {
	id       int
	jobChan  chan string
	stopChan chan struct{}
	wg       *sync.WaitGroup
}

// WorkerPool представляет пул воркеров, которые обрабатывают входящие задания из канала jobChan
type WorkerPool struct {
	jobChan chan string
	workers []*worker
	mu      sync.Mutex
	wg      *sync.WaitGroup
	stopped atomic.Bool
}

// NewWorkerPool создает новый пул воркеров с заданным количеством воркеров
func NewWorkerPool(numWorkers uint) *WorkerPool {
	jobCh := make(chan string)

	workers := make([]*worker, 0, numWorkers)

	p := &WorkerPool{
		jobChan: jobCh,
		mu:      sync.Mutex{},
		workers: workers,
		wg:      &sync.WaitGroup{},
	}

	for i := 0; i < int(numWorkers); i++ {
		_ = p.AddWorker()
	}

	return p
}

// AddWorker добавляет нового воркера в пул и запускает его
// Возвращает ошибку, если пул уже был остановлен
func (p *WorkerPool) AddWorker() error {
	if p.stopped.Load() {
		return errPoolStopped
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	id := len(p.workers) + 1
	w := newWorker(id, p.jobChan, p.wg)
	p.workers = append(p.workers, w)
	go w.work()
	log.Printf("New worker started: id %d\n", id)

	return nil
}

// RemoveWorker удаляет последнего воркера из пула и завершает его работу
// Возвращает ошибку, если нет воркеров для удаления или пул уже был остановлен
func (p *WorkerPool) RemoveWorker() error {
	if p.stopped.Load() {
		return errPoolStopped
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	last := len(p.workers) - 1
	if last < 0 {
		return errCannotRemoveWorker
	}
	w := p.workers[last]
	w.stop()
	p.workers = p.workers[:last]
	log.Printf("Worker removed: id %d\n", w.id)

	return nil
}

// AddJob добавляет задание в канал для обработки воркерами
// Возвращает ошибку, если пул уже был остановлен
func (p *WorkerPool) AddJob(job string) error {
	if p.stopped.Load() {
		return errPoolStopped
	}

	p.wg.Add(1)
	p.jobChan <- job

	return nil
}

// Stop завершает работу пула после того, как все воркеры закончат работу
// Возвращает ошибку, если пул уже был остановлен
func (p *WorkerPool) Stop() error {
	if p.stopped.Load() {
		return errPoolStopped
	}

	close(p.jobChan)

	p.wg.Wait()
	for _, w := range p.workers {
		w.stop()
	}

	p.stopped.Store(true)

	return nil
}

func newWorker(id int, jobChan chan string, wg *sync.WaitGroup) *worker {
	stopChan := make(chan struct{})
	return &worker{
		id:       id,
		jobChan:  jobChan,
		stopChan: stopChan,
		wg:       wg,
	}
}

func (w *worker) work() {
	for {
		select {
		case job, ok := <-w.jobChan:
			if !ok {
				fmt.Printf("Worker %d: job channel closed\n", w.id)
				return
			}
			defer w.wg.Done()
			time.Sleep(500 * time.Millisecond) // имитация сложной работы
			fmt.Printf("Worker %d processing job: %s\n", w.id, job)
		case <-w.stopChan:
			fmt.Printf("Worker %d stopping\n", w.id)
			return
		}
	}
}

func (w *worker) stop() {
	close(w.stopChan)
}
