package config

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/wehw93/worker_pool/pkg/logger"
)

const (
	DefaultConfigPath      = "config/config.yaml"
	DefaultWorkerCount     = 3
	DefaultBufferSize      = 100
	DefaultShutdownTimeout = 5 * time.Second
)

var (
	cfg  *Config
	once sync.Once
)

type Config struct {
	WorkerPool struct {
		InitialWorkerCount int           `yaml:"initial_worker_count" env:"WORKER_COUNT" env-default:"3"`
		BufferSize         int           `yaml:"buffer_size" env:"BUFFER_SIZE" env-default:"100"`
		ShutdownTimeout    time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"5s"`
	} `yaml:"worker_pool"`
	
	Logger struct {
		Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
		Prefix string `yaml:"prefix" env:"LOG_PREFIX" env-default:"[WORKER-POOL] "`
	} `yaml:"logger"`
}

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}
		if err := cleanenv.ReadConfig(DefaultConfigPath, cfg); err != nil {
			if os.IsNotExist(err) {
				cfg.WorkerPool.InitialWorkerCount = DefaultWorkerCount
				cfg.WorkerPool.BufferSize = DefaultBufferSize
				cfg.WorkerPool.ShutdownTimeout = DefaultShutdownTimeout
				cfg.Logger.Level = "info"
				cfg.Logger.Prefix = "[WORKER-POOL] "
			} else {
				panic(err)
			}
		}
	})
	return cfg
}

func (c *Config) GetLogger() logger.Logger {
	level := parseLogLevel(c.Logger.Level)
	return logger.NewStandardLogger(logger.Config{
		Level:  level,
		Prefix: c.Logger.Prefix,
	})
}

func parseLogLevel(levelStr string) logger.LogLevel {
	switch strings.ToLower(levelStr) {
	case "debug":
		return logger.LevelDebug
	case "info":
		return logger.LevelInfo
	case "error":
		return logger.LevelError
	default:
		return logger.LevelInfo
	}
}