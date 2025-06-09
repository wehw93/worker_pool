
package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wehw93/worker_pool/internal/task"
	"github.com/wehw93/worker_pool/pkg/logger"
)

type Pool struct {
	workers         map[int]*Worker
	taskChan        chan task.Task
	workerID        int
	mu              sync.RWMutex
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	logger          logger.Logger
	processor       task.Processor
	shutdownTimeout time.Duration
	shutdownOnce    sync.Once
}

func NewPool(initialWorkerCount int, bufferSize int, shutdownTimeout time.Duration, logger logger.Logger, processor task.Processor) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		workers:         make(map[int]*Worker),
		taskChan:        make(chan task.Task, bufferSize),
		ctx:             ctx,
		cancel:          cancel,
		logger:          logger,
		processor:       processor,
		shutdownTimeout: shutdownTimeout,
	}

	
	for i := 0; i < initialWorkerCount; i++ {
		pool.addWorker()
	}

	return pool
}




func (p *Pool) AddWorker() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.addWorker()
}


func (p *Pool) addWorker() int {
	p.workerID++
	worker := NewWorker(p.workerID, p.taskChan, p.logger, p.processor)
	p.workers[p.workerID] = worker
	
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		worker.Start(p.ctx)
	}()

	p.logger.Printf("Added worker %d", p.workerID)
	return p.workerID
}


func (p *Pool) RemoveWorker(workerID int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	worker, exists := p.workers[workerID]
	if !exists {
		return fmt.Errorf("worker %d not found", workerID)
	}

	worker.Stop()
	delete(p.workers, workerID)

	p.logger.Printf("Removed worker %d", workerID)
	return nil
}

func (p *Pool) SubmitTask(t task.Task) error {
	select {
	case p.taskChan <- t:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool is shutting down")
	default:
		return fmt.Errorf("task queue is full")
	}
}


func (p *Pool) GetWorkerCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.workers)
}


func (p *Pool) Shutdown() error {
	var err error
	p.shutdownOnce.Do(func() {
		p.logger.Info("Shutting down worker pool...")

		
		p.cancel()

		
		close(p.taskChan)

		
		p.mu.Lock()
		for _, worker := range p.workers {
			worker.Stop()
		}
		p.mu.Unlock()

		
		done := make(chan struct{})
		go func() {
			p.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			p.logger.Info("All workers shut down gracefully")
		case <-time.After(p.shutdownTimeout):
			err = fmt.Errorf("shutdown timeout reached after %v", p.shutdownTimeout)
			p.logger.Error("Shutdown timeout reached")
		}
	})
	return err
}