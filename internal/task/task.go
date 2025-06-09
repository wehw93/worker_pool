// task.go
package task

type Task struct {
    Data string
}

type Processor interface {
    Process(workerID int, task Task)
}