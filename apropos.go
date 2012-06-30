package main

import (
	"strings"
	"fmt"
	"errors"
)

func apropos(s string, keys map[string]*s_key) (string, error) {
	var entries []string
	for entry, _ := range keys {
		if strings.HasPrefix(entry, s) {
			log_infoln(s, "matches:", entry)
			entries = append(entries, entry)
		}
	}
	if len(entries) == 1 {
		return entries[0], nil
	} else if len(entries) > 1 {
		e := ""
		for _, i := range entries {
			e = e + i + "\n"
		}
		e = fmt.Sprintf("%v matched:\n%v", s, e)
		return "", errors.New(e)
	}
	return "", errors.New("no key found")
}

