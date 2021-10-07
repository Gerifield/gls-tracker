package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gerifield/gls-tracker/tracker"
)

func main() {
	pkg := flag.String("pkg", "", "Package ID")
	flag.Parse()

	if *pkg == "" {
		log.Fatalln("Set a package id please!")
	}

	t := tracker.New()

	status, err := t.Get(*pkg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Now: %s (%s)\n", status.StatusText, status.ImageText)
}
