package main

import (
	"syscall"
	"unsafe"
)

func echo(e bool) {
	var t syscall.Termios
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0),
		syscall.TCGETS, uintptr(unsafe.Pointer(&t)))
	if e {
		t.Lflag |= syscall.ECHO
	} else {
		t.Lflag &= ^uint32(syscall.ECHO)
	}
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(0),
		syscall.TCSETS, uintptr(unsafe.Pointer(&t)))
}
