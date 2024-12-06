package wg_log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	ERROR int = iota
	WARRING
	PANIC
	FATAL
)

var logFile *os.File

func InitLogger(errLogPath string) {
	if errLogPath == "" {
		errLogPath = "./err.log"
	}
	file, err := os.OpenFile(errLogPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	// 设置log的输出目标为文件
	log.SetOutput(file)
	// 保持句柄,用于优雅关闭
	logFile = file
}

func CloseLogFile() {
	if logFile != nil {
		err := logFile.Close()
		if err != nil {
			Log(ERROR, err)
			return
		}
	}
}

func Log(level int, v ...any) {
	switch level {
	case ERROR:
		Error(v)
	case WARRING:
		Warring(v)
	case PANIC:
		Panic(v)
	case FATAL:
		Fatal(v)
	default:
		Error(v)
	}
}

func Error(v ...any) {
	log.Println(fmt.Sprintf("- %s [ERROR] ", time.Now().Format(time.DateTime)), v)
}

func Warring(v ...any) {
	log.Println(fmt.Sprintf("- %s [WARRING] ", time.Now().Format(time.DateTime)), v)
}

// Panic recover捕获到的时候使用,会自动打印调用栈
func Panic(v ...any) {
	log.Println(fmt.Sprintf("- %s [PANIC] ", time.Now().Format(time.DateTime)), v, GetStack())
}

// Fatal 和log.Fatal一样会自动退出程序,同时也会打印调用栈
func Fatal(v ...any) {
	log.Println(fmt.Sprintf("- %s [FATAL] ", time.Now().Format(time.DateTime)), v, GetStack())
	os.Exit(1)
}

// WarringIf 逻辑错误时可以使用此函数来简单warring
func WarringIf(doWarring bool, v ...any) {
	if doWarring {
		Warring(v)
	}
}

// FatalIfErr 严重错误时可以使用
func FatalIfErr(err error) {
	if err != nil {
		Fatal(err)
	}
}

func GetStack() string {
	// 初始分配一个较小的缓冲区
	stackBuf := make([]byte, 1024)
	// 动态调整缓冲区大小，直到足够大以容纳完整的调用栈
	for {
		stackSize := runtime.Stack(stackBuf, false)
		if stackSize < len(stackBuf) {
			// 如果返回的栈信息小于缓冲区的大小，说明我们已经获取了完整的栈信息
			return string(stackBuf[:stackSize])
		}
		// 如果栈信息被截断，扩展缓冲区大小并继续
		stackBuf = make([]byte, len(stackBuf)*2)
	}
}
