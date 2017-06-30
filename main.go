package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"utilities"
)

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}
	currentUrl := args[0]
	ebooks := utilities.NewEbooks()
	err := ebooks.GetTotalPages(currentUrl)
	checkError(err)
	err = ebooks.GetAllEbooks(currentUrl)
	checkError(err)
	ebooksJson, err := json.Marshal(ebooks)
	checkError(err)
	err = ioutil.WriteFile("output.json", ebooksJson, 0644)
	checkError(err)
}
