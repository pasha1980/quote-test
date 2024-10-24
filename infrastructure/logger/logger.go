package logger

import (
	"io"
	"log"
	"quote-app/config"
)

func Init(wr io.Writer) {
	log.SetOutput(wr)
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	log.SetPrefix(config.Get().AppName + " ")
}
