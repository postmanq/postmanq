package main

import "fmt"

var Constructor = func() error {
	fmt.Println("run config plugin")
	return nil
}
