package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// String 返回日志级别的字符串表示
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger 日志记录器
type Logger struct {
	level      LogLevel
	logger     *log.Logger
	file       *os.File
	mu         sync.Mutex
	minLevel   LogLevel
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// InitLogger 初始化日志系统
func InitLogger(logPath string, minLevel LogLevel) error {
	var err error
	once.Do(func() {
		// 创建日志目录
		logDir := filepath.Dir(logPath)
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		// 打开日志文件
		var file *os.File
		file, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return
		}

		// 创建多写入器，同时输出到文件和控制台
		multiWriter := io.MultiWriter(file, os.Stdout)

		defaultLogger = &Logger{
			logger:   log.New(multiWriter, "", 0),
			file:     file,
			minLevel: minLevel,
		}
	})

	return err
}

// GetLogger 获取默认日志记录器
func GetLogger() *Logger {
	if defaultLogger == nil {
		// 如果未初始化，创建一个只输出到控制台的日志记录器
		defaultLogger = &Logger{
			logger:   log.New(os.Stdout, "", 0),
			minLevel: INFO,
		}
	}
	return defaultLogger
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log 内部日志记录方法
func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, v...)
	logLine := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), message)

	l.logger.Println(logLine)
}

// Debug 记录调试级别日志
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

// Info 记录信息级别日志
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

// Warn 记录警告级别日志
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

// Error 记录错误级别日志
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

// 全局便捷函数

// Debug 记录调试级别日志
func Debug(format string, v ...interface{}) {
	GetLogger().Debug(format, v...)
}

// Info 记录信息级别日志
func Info(format string, v ...interface{}) {
	GetLogger().Info(format, v...)
}

// Warn 记录警告级别日志
func Warn(format string, v ...interface{}) {
	GetLogger().Warn(format, v...)
}

// Error 记录错误级别日志
func Error(format string, v ...interface{}) {
	GetLogger().Error(format, v...)
}
