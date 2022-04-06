package main

import (
	"errors"
	"strings"
)

func Parse(message string)([]string, error){
	if len(message) < 1 {
		return nil, errors.New("message is empty")
	}

	if message[0:1] == "?" {
		return strings.Split(message[1:], " "), nil
	} else {
		return nil, errors.New("no es un command")
	}
}