package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 || len(os.Args) == 1 {
		fmt.Println("Please provide one argument: path to file or directory")
		return
	}

	arg := os.Args[1]

	zipManager := &ZipManager{}
	docxConverter := &DocxConverter{}

	processor := FileProcessor{docxConverter, zipManager}

	if err := processor.Process(arg); err != nil {
		log.Fatal(err)
	}
}
