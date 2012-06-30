package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"io"
)

// help
func handle_h() {
	usage()
	os.Exit(0)
}

// import plaintext
func handle_i(filename string) {
	password, keys := get_private_data(true)

	// file format is pairs of lines, first the entry name, then the
	// password. the following line is blank, then EOF or another keyset.
	file, err := os.Open(filename)
	if err != nil {
		log_dieln(err)
	}
	reader := bufio.NewReader(file)
	for {
		var entry, key string
		line, _, err := reader.ReadLine() //entry
		if err == io.EOF {
			break
		} else if err != nil {
			log_dieln(err)
		}
		entry = string(line)

		line, _, err = reader.ReadLine() //key
		if err == io.EOF {
			log_die("missing key for entry: %v\n", entry)
		} else if err != nil {
			log_dieln(err)
		}
		key = string(line)

		line, _, err = reader.ReadLine() //newline
		if err != io.EOF {
			log_dieln("missing trailing newline")
		} else if err != nil {
			log_dieln(err)
		}
		if string(line) != "" {
			log_dieln("invalid input\n")
		}

		if keys[entry] != nil {
			log_error("entry %v already exists, ignoring new entry", entry)
		} else {
			fmt.Printf("adding entry %v\n", entry)
			keys[entry] = &s_key{Key: key}
		}
	}
	encrypt(password, keys)
}

// export plaintext
func handle_e(filename string) {
	_, keys := get_private_data(false)

	file, err := os.Create(filename)
	if err != nil {
		log_dieln(err)
	}

	for entry, key := range keys {
		fmt.Fprintf(file, "%v\n%v\n\n", entry, key.Key)
	}
}

// new
func handle_n(entry string) {
	password, keys := get_private_data(true)

	// check for collisions
	if keys[entry] != nil {
		log_die("entry %s already exists\n", entry)
	}

	log_info("adding entry for %s\n", entry)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter key (or leave blank to have one created): ")
	line, _, err := reader.ReadLine()
	if err != nil {
		log_dieln(err)
	}
	if string(line) == "" {
		// auto create a key
		fmt.Println("generating key...")
		line = generate_key()
		log_info("generated key:", string(line))
	}
	keys[entry] = &s_key{Key: string(line)}
	encrypt(password, keys)
	fmt.Printf("added %v.\n", entry)
}

// modify
func handle_m(entry string) {
	password, keys := get_private_data(false)

	// check for entry
	if keys[entry] == nil {
		log_die("entry %s does not exist\n", entry)
	}

	log_info("modifying entry for %s\n", entry)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter key (or leave blank to have one created): ")
	line, _, err := reader.ReadLine()
	if err != nil {
		log_dieln(err)
	}
	if string(line) == "" {
		// auto create a key
		fmt.Println("generating key...")
		line = generate_key()
		log_info("generated key:", string(line))
	}
	keys[entry].Key = string(line)
	encrypt(password, keys)
	fmt.Printf("modified %v.\n", entry)
}

// delete
func handle_d(entry string) {
	password, keys := get_private_data(false)

	// check for entry
	if keys[entry] == nil {
		log_die("entry %s does not exist\n", entry)
	}

	log_info("deleting entry for %s\n", entry)

	delete(keys, entry)
	encrypt(password, keys)
	fmt.Printf("deleted %v.\n", entry)
}

// change key
func handle_p() {
	_, keys := get_private_data(false)
	password := generate_password()
	encrypt(password, keys)
	fmt.Println("password updated successfully.")
}

// list entries
func handle_l() {
	_, keys := get_private_data(false)

	var entries []string
	for entry, _ := range keys {
		entries = append(entries, entry)
	}
	sort.Strings(entries)
	for _, entry := range entries {
		fmt.Println(entry)
	}
}

// read an entry
func handle_read(e string) {
	_, keys := get_private_data(false)
	entry, err := apropos(e, keys)
	if err != nil {
		log_dieln(err)
	}

	key := keys[entry]
	if key == nil {
		log_die("no such entry %v\n", entry)
	}
	fmt.Println(key.Key)
}
