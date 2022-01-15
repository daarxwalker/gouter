package main

import (
	"fmt"
)

func main() {
	fmt.Println("-- GOUTER SIMPLE REVERSE PROXY --")
	proxy := newGouter()
	proxy.run(proxy.readConfig())
}
