package out

import (
	"fmt"
	"os"
	"time"

	ct "github.com/daviddengcn/go-colortext"
)

const (
	timeFormat = "02-01-2006 15:04:05.000"
)

var (
	debug bool
)

var (
	ErrorHandler func(string)
)

func SetDebug(d bool) {
	debug = d
}

func Fatal(msg ...interface{}) {
	ct.Foreground(ct.Red, true)
	fmt.Print("[F]: " + time.Now().Format(timeFormat) + " -> ")
	fmt.Println(msg...)
	os.Exit(1)
}

func Err(handled bool, msg ...interface{}) {
	ct.Foreground(ct.Red, true)
	fmt.Print("[E]: " + time.Now().Format(timeFormat) + " -> ")
	fmt.Println(msg...)
	ct.ResetColor()

	if handled {
		if ErrorHandler != nil {
			ErrorHandler(fmt.Sprint(msg...))
		}
	}
}

func Info(msg ...interface{}) {
	ct.Foreground(ct.Green, true)
	fmt.Print(msg...)
	ct.ResetColor()
}

func Infoln(msg ...interface{}) {
	ct.Foreground(ct.Green, true)
	fmt.Println(msg...)
	ct.ResetColor()
}

func Infof(format string, a ...interface{}) {
	ct.Foreground(ct.Green, true)
	fmt.Printf(format, a...)
	ct.ResetColor()
}

func Debug(a ...interface{}) {
	if !debug {
		return
	}
	fmt.Print(a...)
}

func Debugln(a ...interface{}) {
	if !debug {
		return
	}
	fmt.Println(a...)
}

func Debugf(format string, a ...interface{}) {
	if !debug {
		return
	}
	fmt.Printf(format, a...)
}
