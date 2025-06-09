// main.go
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

type SimpleTaskProcessor struct {
	logger logger.Logger
}

func (p *SimpleTaskProcessor) Process(workerID int, t task.Task) {
	p.logger.Printf("Worker %d processing task: %s", workerID, t.Data)
	time.Sleep(100 * time.Millisecond) 
}

func main() {
	
	cfg := config.NewDefaultConfig()

	
	processor := &SimpleTaskProcessor{
		logger: cfg.Logger,
	}

	
	pool := worker.NewPool(*cfg, processor)
	defer pool.Shutdown()

	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	
	go func() {
		for i := 0; i < 20; i++ {
			task := task.Task{
				Data: fmt.Sprintf("Task %d", i+1),
			}

			if err := pool.SubmitTask(task); err != nil {
				cfg.Logger.Printf("Failed to submit task: %v", err)
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()

	
	go func() {
		time.Sleep(500 * time.Millisecond)
		cfg.Logger.Info("Adding 2 more workers...")
		pool.AddWorker()
		pool.AddWorker()

		time.Sleep(1 * time.Second)
		cfg.Logger.Info("Removing 1 worker...")
		if err := pool.RemoveWorker(1); err != nil {
			cfg.Logger.Printf("Failed to remove worker: %v", err)
		}
	}()

	
	<-sigChan
	cfg.Logger.Info("Received shutdown signal")
}
