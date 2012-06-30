package main

import (
	"bufio"
	"fmt"
	"os"
)

func generate_password() string {
	reader := bufio.NewReader(os.Stdin)
	echo(false)
	fmt.Print("Enter new password: ")
	line, _, err := reader.ReadLine()
	pass_one := string(line)
	echo(true)
	fmt.Println()
	if err != nil {
		log_dieln(err)
	}
	log_info("got password:", pass_one)
	echo(false)
	fmt.Print("Enter password again: ")
	line, _, err = reader.ReadLine()
	pass_two := string(line)
	echo(true)
	fmt.Println()
	if err != nil {
		log_dieln(err)
	}
	log_info("got password:", pass_two)
	if pass_one != pass_two {
		log_dieln("passwords to not match")
	}
	return string(pass_one)
}

func get_password() string {
	echo(false)
	fmt.Print("Enter password: ")
	reader := bufio.NewReader(os.Stdin)
	password, _, err := reader.ReadLine()
	echo(true)
	fmt.Println()
	if err != nil {
		log_dieln(err)
	}
	return string(password)
}
