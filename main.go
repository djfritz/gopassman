package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"os"
	"os/signal"
)

// TODO: document the correct way!

// flags:
// -i import unencrypted data
// -e export unencrypted data
// -n new entry, needs some other fields?
// <arg0> the entry
// -l list entries
// -d remove entry
// -m modify entry
// -p change password
// -k use a different keyfile
// usage: gp -i <some plaintext file>
//        gp -e <some plaintext file>
//	  gp -n amazon
//	  gp -l
//	  gp -d amazon
//	  gp -m amazon
// 	  gp -p
//	  gp -k ~/different.gp
// Commands can stack and should be evaluated in the following order:
// -k <>, -i <>, -m <>, -d <>, -n <>, -e <>, -p, -l, <entry>
// -l is exclusive to all others, and will error is muxed with the others

var (
	default_path    string = os.Getenv("HOME") + "/.gp/keys.gp"
	cached_password string = ""

	f_i = flag.String("i", "", "import unencrypted keys")
	f_e = flag.String("e", "", "export unencrypted keys")
	f_n = flag.String("n", "", "new entry")
	f_l = flag.Bool("l", false, "list entries")
	f_d = flag.String("d", "", "delete entry")
	f_m = flag.String("m", "", "modify an entry")
	f_p = flag.Bool("p", false, "change master password")
	f_k = flag.String("k", default_path, "keyfile to use")
	f_h = flag.Bool("h", false, "display detailed help")
	f_v = flag.Bool("v", false, "display extra info")
)

func usage() {
	fmt.Println("usage: gp <options> <entry>")
	flag.PrintDefaults()
	if *f_h {
		fmt.Println()
		fmt.Println("flags can stack and are evaluated in order of:")
		fmt.Println("\tk i m d n e p l")
		fmt.Println("\tyou can import, merge, delete, etc.,",
			"in one command")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("\tlookup an entry named gopher:")
		fmt.Println("\t\tgp gopher")
		fmt.Println("\tembed an entry on the command line:")
		fmt.Println("\t\techo \"password = `gp gopher`\"")
		fmt.Println("\tmerge into your local db, delete entry gopher,",
			"and lookup entry bunny")
		fmt.Println("\t\tgp -k otherkeys.gp -d gopher bunny")
	}
}

func init() {
	flag.Usage = usage
}

type s_key struct {
	Key string
}

type keymap map[string]*s_key

func (k keymap) Encode() []byte {
	buffer := new(bytes.Buffer)
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(k)
	if err != nil {
		log_dieln(err)
	}
	return buffer.Bytes()
}

func (k keymap) Decode(input []byte) {
	buffer := bytes.NewBuffer(input)
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(&k)
	if err != nil {
		log_dieln(err)
	}
}

func main() {
	flag.Parse()
	sane := false

	// trap signals so we can try to get echo turned back on before
	// we bail
	sig := make(chan os.Signal, 1024)
	signal.Notify(sig)
	go func() {
		<-sig
		echo(true)
		os.Exit(0)
	}()

	// parse the flags in a specific order. if any fails, then bail.
	// some flags, like -h, will call exit themselves.
	if *f_h {
		// we won't get past this one
		handle_h()
		sane = true
	}
	if *f_i != "" {
		handle_i(*f_i)
		sane = true
	}
	if *f_m != "" {
		handle_m(*f_m)
		sane = true
	}
	if *f_d != "" {
		handle_d(*f_d)
		sane = true
	}
	if *f_n != "" {
		handle_n(*f_n)
		sane = true
	}
	if *f_e != "" {
		handle_e(*f_e)
		sane = true
	}
	if *f_p {
		handle_p()
		sane = true
	}
	if *f_l {
		handle_l()
		sane = true
	}

	// special case, if no regular flags are given (v,h are special)
	// and no extra args have been given, then throw an error
	if flag.NArg() == 0 && sane == false {
		log_errorln("main: not enough arguments")
		usage()
	} else if flag.NArg() > 1 {
		log_errorln("main: too many arguments")
		usage()
	} else if flag.NArg() == 1 {
		handle_read(flag.Arg(0))
	}
}

func get_private_data(create bool) (string, keymap) {
	var keys keymap = make(keymap)
	file := file_open()
	if file == nil {
		if create {
			fmt.Println("creating new keyfile:", *f_k)
			cached_password = generate_password()
			file = file_create()
		} else {
			os.Exit(1)
		}
	} else {
		if cached_password == "" {
			cached_password = get_password()
		}
		keys = decrypt(file, cached_password)
	}
	file.Close()
	return cached_password, keys
}
