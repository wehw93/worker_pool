
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wehw93/worker_pool/internal/config"
	"github.com/wehw93/worker_pool/internal/task"
	"github.com/wehw93/worker_pool/internal/worker"
	"github.com/wehw93/worker_pool/pkg/logger"
)

type Processor struct {
	logger logger.Logger
}

func (p *Processor) Process(workerID int, t task.Task) {
	p.logger.Printf("Worker %d processing task: %s", workerID, t.Data)
	time.Sleep(100 * time.Millisecond)
}

func main() {
	
	cfg := config.GetConfig()
	
	
	appLogger := cfg.GetLogger()
	
	
	processor := &Processor{
		logger: appLogger,
	}

	
	pool := worker.NewPool(
		cfg.WorkerPool.InitialWorkerCount,
		cfg.WorkerPool.BufferSize,
		cfg.WorkerPool.ShutdownTimeout,
		appLogger,
		processor,
	)
	
	defer func() {
		if err := pool.Shutdown(); err != nil {
			appLogger.Printf("Error during shutdown: %v", err)
		}
	}()

	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	
	go func() {
		appLogger.Info("Starting task submission...")
		for i := 0; i < 30; i++ {
			task := task.Task{
				Data: fmt.Sprintf("Task %d", i+1),
			}

			if err := pool.SubmitTask(task); err != nil {
				appLogger.Printf("Failed to submit task: %v", err)
				continue
			}

			time.Sleep(50 * time.Millisecond)
		}
		appLogger.Info("Finished submitting tasks")
	}()

	
	go func() {
		time.Sleep(500 * time.Millisecond)
		
		appLogger.Printf("Current worker count: %d", pool.GetWorkerCount())
		
		appLogger.Info("Adding 5 more workers...")
		for i := 0; i < 5; i++ {
			workerID := pool.AddWorker()
			appLogger.Printf("Added worker with ID: %d", workerID)
		}
		
		appLogger.Printf("Worker count after adding: %d", pool.GetWorkerCount())

		time.Sleep(1 * time.Second)
		
		appLogger.Info("Removing 1 worker...")
		if err := pool.RemoveWorker(1); err != nil {
			appLogger.Printf("Failed to remove worker: %v", err)
		} else {
			appLogger.Printf("Worker count after removal: %d", pool.GetWorkerCount())
		}
	}()

	<-sigChan
	appLogger.Info("Received shutdown signal, initiating graceful shutdown...")
}