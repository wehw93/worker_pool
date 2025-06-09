// pool.go
package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wehw93/worker_pool/internal/config"
	"github.com/wehw93/worker_pool/internal/task"
	"github.com/wehw93/worker_pool/pkg/logger"
)

type Pool struct {
	workers      map[int]*Worker
	taskChan     chan task.Task
	workerID     int
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	logger       logger.Logger
	processor    task.Processor
	config       config.WorkerPoolConfig
	shutdownOnce sync.Once
}

func NewPool(cfg config.WorkerPoolConfig, processor task.Processor) *Pool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		workers:   make(map[int]*Worker),
		taskChan:  make(chan task.Task, cfg.BufferSize),
		ctx:       ctx,
		cancel:    cancel,
		logger:    cfg.Logger,
		processor: processor,
		config:    cfg,
	}

	
	for i := 0; i < cfg.InitialWorkerCount; i++ {
		pool.AddWorker()
	}

	return pool
}

func (p *Pool) AddWorker() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.workerID++
	worker := NewWorker(p.workerID, p.taskChan, p.logger, p.processor)
	p.workers[p.workerID] = worker
	worker.Start(p.ctx)

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

func (p *Pool) Shutdown() error {
	var err error
	p.shutdownOnce.Do(func() {
		p.logger.Info("Shutting down worker pool...")

		
		p.cancel()

		
		close(p.taskChan)

		
		done := make(chan struct{})
		go func() {
			p.wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			p.logger.Info("All workers shut down")
		case <-time.After(p.config.ShutdownTimeout):
			err = fmt.Errorf("shutdown timeout")
			p.logger.Error("Shutdown timeout reached")
		}
	})
	return err
}
