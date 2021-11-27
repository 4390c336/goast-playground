package main

import (
	"os"
)

func another_main() []byte {
	data, err := os.ReadFile("/tmp/dat")
	if err != nil {
		panic("can't open file")
	}
	return data
}

/*
some blah blah here
*/
func AnotherFunction_fun() bool {
	//this comment will break everything
	return false
}
