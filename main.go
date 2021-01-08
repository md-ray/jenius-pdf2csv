// main.go

package main

import (
	"log"
	"os"
)

func main() {
	a := App{}

	if len(os.Args) <= 1 {
		log.Fatal("Usage: jenius-pdf2csv <pdf-filename>")
	}

	a.Initialize(os.Args[1])
	a.Run()
}
