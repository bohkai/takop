package main

import (
	"errors"
	"strings"
)

func Parse(message string)([]string, error){
	if message[0:1] == "!" {
		return strings.Split(message[1:], " "), nil
	} else {
		return nil, errors.New("no es un command")
	}
}