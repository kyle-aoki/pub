package main

import "fmt"

func mainRecover() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

func mustExec(err error) {
	if err != nil {
		panic(err)
	}
}
