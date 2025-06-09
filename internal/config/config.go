// config.go
package config

import (
    "time"
    "github.com/wehw93/worker_pool/pkg/logger"
)

const (
    DefaultWorkerCount     = 3
    DefaultBufferSize      = 100
    DefaultShutdownTimeout = 5 * time.Second
)

type WorkerPoolConfig struct {
    InitialWorkerCount int
    BufferSize         int
    ShutdownTimeout    time.Duration
    Logger             logger.Logger
}

func NewDefaultConfig() *WorkerPoolConfig {
    return &WorkerPoolConfig{
        InitialWorkerCount: DefaultWorkerCount,
        BufferSize:         DefaultBufferSize,
        ShutdownTimeout:    DefaultShutdownTimeout,
        Logger: logger.NewStandardLogger(logger.Config{
            Level:  logger.LevelInfo,
            Prefix: "[WORKER-POOL] ",
        }),
    }
}