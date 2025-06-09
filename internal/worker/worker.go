
package worker

import (
    "context"
    "sync"
    "github.com/wehw93/worker_pool/internal/task"
    "github.com/wehw93/worker_pool/pkg/logger"
)

type Worker struct {
    id       int
    taskChan <-chan task.Task
    quit     chan struct{}
    wg       *sync.WaitGroup
    logger   logger.Logger
    processor task.Processor
    stopped  bool
    mu       sync.RWMutex
}

func NewWorker(id int, taskChan <-chan task.Task, logger logger.Logger, processor task.Processor) *Worker {
    return &Worker{
        id:       id,
        taskChan: taskChan,
        quit:     make(chan struct{}),
        wg:       &sync.WaitGroup{},
        logger:   logger,
        processor: processor,
    }
}

func (w *Worker) Start(ctx context.Context) {
    w.mu.Lock()
    if w.stopped {
        w.mu.Unlock()
        return
    }
    w.mu.Unlock()

    w.wg.Add(1)
    go w.run(ctx)
}

func (w *Worker) Stop() {
    w.mu.Lock()
    defer w.mu.Unlock()
    
    if w.stopped {
        return
    }
    
    w.stopped = true
    close(w.quit)
    w.wg.Wait()
}

func (w *Worker) run(ctx context.Context) {
    defer w.wg.Done()
    w.logger.Printf("Worker %d started", w.id)

    for {
        select {
        case task, ok := <-w.taskChan:
            if !ok {
                w.logger.Printf("Worker %d: task channel closed", w.id)
                return
            }
            if w.processor != nil {
                w.processor.Process(w.id, task)
            } else {
                w.logger.Printf("Worker %d: processor is nil, skipping task: %s", w.id, task.Data)
            }
            
        case <-w.quit:
            w.logger.Printf("Worker %d: received quit signal", w.id)
            return
            
        case <-ctx.Done():
            w.logger.Printf("Worker %d: context cancelled", w.id)
            return
        }
    }
}