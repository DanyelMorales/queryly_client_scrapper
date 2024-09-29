package style

import (
	"fmt"
	"log"
	"strings"
)

func OnReportF(format string, a ...interface{}) {
	LogF(Documentation, format, a...)
}

func OnReport(format string) {
	Log(Documentation, format)
}

func OnNoticeF(format string, a ...interface{}) {
	LogF(Notice, format, a...)
}

func OnNotice(format string) {
	Log(Notice, format)
}

func OnTitle(title string) {
	Log(AddonEnable, title)
}

func OnTitleF(format string, a ...interface{}) {
	LogF(AddonEnable, format, a...)
}

func OnAction(format string) {
	Log(Running, format)
}
func SuccessfulAction(format string) {
	Log(ThumbsUp, format)
}
func FailedAction(format string) {
	Log(Failure, format)
}
func ExitAction(format string) {
	Log(Deleted, format)
}
func ExitActionF(format string, a ...interface{}) {
	LogF(Deleted, format, a...)
}
func OnActionF(format string, a ...interface{}) {
	LogF(Running, format, a...)
}
func SuccessfulActionF(format string, a ...interface{}) {
	LogF(ThumbsUp, format, a...)
}
func FailedActionF(format string, a ...interface{}) {
	LogF(Failure, format, a...)
}
func ActLogCF(print bool, format string, a ...interface{}) {
	if !print {
		return
	}
	LogF(Running, format, a...)
}
func OkLogCF(print bool, format string, a ...interface{}) {
	if !print {
		return
	}
	LogF(ThumbsUp, format, a...)
}
func ErrorLogCF(print bool, format string, a ...interface{}) {
	if !print {
		return
	}
	LogF(Failure, format, a...)
}

func PrintF(st Enum, format string, a ...interface{}) {
	msg := StringF(st, format, a...)
	fmt.Print(msg)
}

func Print(st Enum, format string) {
	msg := String(st, format)
	fmt.Print(msg)
}

func LogF(st Enum, format string, a ...interface{}) {
	msg := StringF(st, format, a...)
	log.Print(msg)
}

func Log(st Enum, format string) {
	msg := String(st, format)
	log.Print(msg)
}

// StringF writes a basic formatted string to stdout
func String(st Enum, format string) string {
	prefix, ok := Config[st]
	if ok {
		var color = string(prefix.Color)
		str := color + prefix.Prefix + format + "\n"
		return strings.ToLower(str)
	}
	return ""
}

func StringF(st Enum, format string, a ...interface{}) string {
	str := String(st, format)
	return fmt.Sprintf(str, a...)
}
