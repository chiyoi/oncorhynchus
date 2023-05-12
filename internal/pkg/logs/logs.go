package logs

import (
	"log"
	"os"
)

var (
	infoLogger    = log.New(os.Stderr, "(info) ", log.Lmsgprefix|log.LstdFlags|log.LUTC)
	warningLogger = log.New(os.Stderr, "(warning) ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	errorLogger   = log.New(os.Stderr, "(error) ", log.Lmsgprefix|log.Llongfile|log.LstdFlags|log.LUTC)
	panicLogger   = log.New(os.Stderr, "(fatal) ", log.Lmsgprefix|log.Llongfile|log.LstdFlags|log.LUTC)
	debugLogger   = log.New(os.Stderr, "(debug) ", log.Lmsgprefix|log.Llongfile|log.LstdFlags|log.LUTC)
)

func Info(v ...any)    { infoLogger.Println(v...) }
func Warning(v ...any) { warningLogger.Println(v...) }
func Error(v ...any)   { errorLogger.Println(v...) }
func Panic(v ...any)   { panicLogger.Panicln(v...) }
func Debug(v ...any)   { debugLogger.Println(v...) }
