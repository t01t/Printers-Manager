package main

import "github.com/t01t/printers-manager/server"

func main() {
	server.Init()
	defer server.DefaultPrinter.Printer.Close()
}
