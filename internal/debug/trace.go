package debug

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var Debug = false

func init() {
	// og.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetFlags(0)
	log.SetPrefix("zipline: ")
	log.SetOutput(os.Stderr)
}

func Trace(msg string, p ...interface{}) {
	if !Debug {
		return
	}

	pc, _, line, _ := runtime.Caller(1)

	caller := fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
	if len(p) == 0 {
		log.Println(caller + " - " + msg)
	} else {
		log.Printf(caller+" - "+msg, p...)
	}
}
