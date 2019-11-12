//+build !test

package main

import (
	"log"
	"os"
)

const name = "logger"
const version = "1.1.0"

func main() {
	var svc Service

	err := svc.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	if err != nil {
		os.Exit(-1)
	}
}
