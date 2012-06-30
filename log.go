package main

import (
	"log"
	"os"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(0)
}

func log_prologue() {
	_, file, line, _ := runtime.Caller(2)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	log.SetFlags(0)
	log.SetPrefix(short + ":" + strconv.Itoa(line) + ": ")
}

func log_epilogue() {
	log.SetPrefix("")
}

func log_info(format string, arg ...interface{}) {
	if *f_v {
		log_prologue()
		log.Printf(format, arg...)
		log_epilogue()
	}
}

func log_infoln(arg ...interface{}) {
	if *f_v {
		log_prologue()
		log.Println(arg...)
		log_epilogue()
	}
}

func log_errorln(arg ...interface{}) {
	if *f_v {
		log_prologue()
	}
	log.SetOutput(os.Stderr)
	log.Println(arg...)
	log.SetOutput(os.Stdout)
	log_epilogue()
}

func log_error(format string, arg ...interface{}) {
	if *f_v {
		log_prologue()
	}
	log.SetOutput(os.Stderr)
	log.Printf(format, arg...)
	log.SetOutput(os.Stdout)
	log_epilogue()
}

func log_dieln(arg ...interface{}) {
	if *f_v {
		log_prologue()
	}
	log.SetOutput(os.Stderr)
	log.Println(arg...)
	log.SetOutput(os.Stdout)
	os.Exit(1)
}

func log_die(format string, arg ...interface{}) {
	if *f_v {
		log_prologue()
	}
	log.SetOutput(os.Stderr)
	log.Printf(format, arg...)
	log.SetOutput(os.Stdout)
	os.Exit(1)
}
