package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var LogTimeFormat string = "2006-01-02T15:04:05.000000"
var LogFormat4Negroni string = "{{.StartTime}} {{.Method}} {{.Path}} {{.Status}} {{.Duration}}"

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("[SAM-API] " + time.Now().Format(LogTimeFormat) + " " + string(bytes))
}

//
// Use custom log fomat
//
func LogInit(silent bool) {
	if !silent {
		log.SetFlags(0)
		log.SetOutput(new(logWriter))
	} else {
		log.SetOutput(ioutil.Discard)
	}
}

//
// Lof each request contents in a chain of handlers
//
func WithLog(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Printf("Log request: %s", RequestInfo(r))
	next(w, r)
}
