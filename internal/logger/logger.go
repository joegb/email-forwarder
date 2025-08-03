package logger

import (
	"log"
	"os"
	"io"
	"sync"
)

var (
	once sync.Once
	file *os.File
)

// 初始化日志
func Init() {
	once.Do(func() {
		// 创建日志目录
		if err := os.MkdirAll("logs", 0755); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}
		
		// 打开日志文件
		var err error
		file, err = os.OpenFile("logs/mail_processing.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Failed to open mail log file: %v", err)
		}
		
		log.SetOutput(io.MultiWriter(os.Stdout, file))
	})
}

func Info(format string, v ...interface{}) {
	Init()
	log.Printf("[INFO] "+format, v...)
}

func Error(format string, v ...interface{}) {
	Init()
	log.Printf("[ERROR] "+format, v...)
}

func Warn(format string, v ...interface{}) {
	Init()
	log.Printf("[WARN] "+format, v...)
}