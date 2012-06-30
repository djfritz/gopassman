package main

import (
	"os"
	"path"
)

func file_open() *os.File {
	fi, err := os.Stat(*f_k)
	if err == nil {
		if fi.Mode().Perm() != 0600 {
			log_die("invalid permissions on %v, must be set to 0600\n", 
				*f_k)
		}
	}

	file, err := os.Open(*f_k)

	if os.IsNotExist(err) {
		log_errorln("keyfile does not exist")
		return nil
	} else if err != nil {
		log_dieln(err)
	}
	return file
}

func file_create() *os.File {
	p := path.Dir(*f_k)
	perm := os.FileMode(0700)

	err := os.MkdirAll(p, perm)
	if err != nil {
		log_dieln(err)
	}

	fileperm := os.FileMode(0600)
	file, err := os.OpenFile(*f_k, os.O_TRUNC|os.O_CREATE|os.O_RDWR, fileperm)
	if err != nil {
		log_dieln(err)
	}
	return file
}
