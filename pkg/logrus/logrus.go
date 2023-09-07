package logrus

import (
	"fmt"
	"time"
)

func Info(msg ...interface{}) {
	msg = prepend(msg, fmt.Sprintf("[INFO] %s", time.Now().Format(time.RFC3339)))
	fmt.Println(msg...)
}

func Infof(msg string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[INFO] %s %s\n", time.Now().Format(time.RFC3339), msg), args...)
}

func Debug(msg ...interface{}) {
	msg = prepend(msg, fmt.Sprintf("[DBUG] %s", time.Now().Format(time.RFC3339)))
	fmt.Println(msg...)
}

func Debugf(msg string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[DBUG] %s %s\n", time.Now().Format(time.RFC3339), msg), args...)
}

func Error(msg ...interface{}) {
	msg = prepend(msg, fmt.Sprintf("[EROR] %s", time.Now().Format(time.RFC3339)))
	fmt.Println(msg...)
}

func Errorf(msg string, args ...interface{}) {
	fmt.Printf(fmt.Sprintf("[EROR] %s %s\n", time.Now().Format(time.RFC3339), msg), args...)
}

func prepend(items []interface{}, item interface{}) []interface{} {
	return append([]interface{}{item}, items...)
}
