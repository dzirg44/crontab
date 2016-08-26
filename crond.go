package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
)

var port *string = flag.String("port", "127.0.0.1:4444", "web port")
var logs *string = flag.String("logs", "/var/log/croncli/", "log path")
var conf *string = flag.String("conf", "crontab.conf", "crontab config")
var stopCh chan bool = make(chan bool)
var startCh chan bool = make(chan bool)

const (
	RUN_LOG_POSTFIX = `run.log`
	SVR_LOG         = `sys.log`
	DATEFORMAT      = `20060102`
	TIMEFORMAT      = `2006-01-02 15:04:05`
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	initLog()

	loaded, loadErr := loadConf()
	if !loaded {
		fmt.Printf("Err %s exit.\n", loadErr)
		os.Exit(1)
	}

	go jobHandle()

	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	http.HandleFunc("/del", del)
	http.HandleFunc("/log", loger)
	http.HandleFunc("/load", load)
	http.HandleFunc("/stop", stop)
	http.HandleFunc("/start", start)
	http.HandleFunc("/status", status)

	startErr := http.ListenAndServe(*port, nil)
	if startErr != nil {
		fmt.Println("Start server failed.", startErr)
		os.Exit(1)
	}
}
