package main

import (
	"github.com/spf13/pflag"
)

func main() {
	var filename string
	pflag.StringP(&filename, "config", "c", "config filename")
	if !pflag.CommandLine.Parsed() {
		pflag.CommandLine.PrintDefaults()
	}

	//pflag.StringP("config", "c", "./config.yml", "config filename")
	//
	//pflag.Parse()
	//if !pflag.Parsed() {
	//	pflag.PrintDefaults()
	//}
}
