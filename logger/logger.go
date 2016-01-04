package logger

import (
	"log"
	"github.com/fatih/color"
	"os"
)

func Fatal(values ...interface{}) {
	color.Set(color.FgHiRed)
	log.Println(values...)
	color.Unset()
	os.Exit(1)
}
func Fatalf(format string, values ...interface{}) {
	color.Set(color.FgHiRed)
	log.Printf(format, values...)
	color.Unset()
	os.Exit(1)
}

func Warning(values ...interface{}) {
	color.Set(color.FgYellow)
	log.Println(values...)
	color.Unset()
}

func Success(values ...interface{}) {
	color.Set(color.FgGreen)
	log.Println(values...)
	color.Unset()
}
func Successf(format string, values ...interface{}) {
	color.Set(color.FgGreen)
	log.Printf(format, values...)
	color.Unset()
}

func Info(values ...interface{}) {
	color.Set(color.FgCyan)
	log.Println(values...)
	color.Unset()
}