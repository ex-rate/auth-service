package main

import (
	"flag"

	app "github.com/ex-rate/auth-service/cmd"
)

func main() {
	var path, fileName string

	flag.StringVar(&path, "path", ".", "path to config file")
	flag.StringVar(&fileName, "name", ".env", "config file name")
	flag.Parse()

	app.Start(path, fileName)
}
