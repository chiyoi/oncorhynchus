package logs

import (
	"log"
	"os"
)

var (
	infoLogger    = log.New(os.Stderr, "(info) ", log.Lmsgprefix|log.LstdFlags|log.LUTC)
	warningLogger = log.New(os.Stderr, "(warning) ", log.Lmsgprefix|log.Lshortfile|log.LstdFlags|log.LUTC)
	errorLogger   = log.New(os.Stderr, "(error) ", log.Lmsgprefix|log.Llongfile|log.LstdFlags|log.LUTC)
	fatalLogger   = log.New(os.Stderr, "(fatal) ", log.Lmsgprefix|log.Llongfile|log.LstdFlags|log.LUTC)
	debugLogger   = log.New(os.Stderr, "(debug) ", log.Lmsgprefix|log.Llongfile|log.LstdFlags|log.LUTC)
)

var (
	Info    = infoLogger.Println
	Warning = warningLogger.Println
	Error   = errorLogger.Println
	Panic   = fatalLogger.Panicln
	Debug   = debugLogger.Println
)
