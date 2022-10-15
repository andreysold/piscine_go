package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Flags struct {
	Directory *bool
	File      *bool
	Link      *bool
	All       *bool
}

func recursiveTree(mainPath string, flags Flags) {
	if !strings.HasPrefix(mainPath, "/") {
		mainPath += "/"
	}
	data, err := ioutil.ReadDir(mainPath)
	if err != nil {
		panic(err)
	}
	for i := range data {
		if data[i].IsDir() {
			if *flags.Directory || *flags.All {
				fmt.Println("/" + mainPath + data[i].Name())
			}
			recursiveTree(mainPath+data[i].Name(), flags)
		} else {
			link, errLink := os.Readlink(mainPath + data[i].Name())
			if errLink == nil {
				_, err := os.Open(mainPath + data[i].Name())
				if err != nil && (*flags.Link || *flags.All) {
					fmt.Println(link + "-> [broken]")
				} else {
					fmt.Println(mainPath + data[i].Name() + " -> " + mainPath + link)
				}
			}
			if *flags.File || *flags.All {
				fmt.Println("/" + mainPath + data[i].Name())
			}
		}
	}
}

func getPath(args []string) (string, error) {
	var filename string
	for i := range args {
		if !strings.HasPrefix(args[i], "-") {
			filename = args[i]
		}
	}
	if filename == "" {
		return "", errors.New("uncorrected path name")
	}
	return filename, nil
}

func main() {
	var (
		flags    Flags
		mainPath string
		err      error
	)
	flags.File = flag.Bool("f", false, "files flag")
	flags.Directory = flag.Bool("d", false, "directories flag")
	flags.Link = flag.Bool("sl", false, "links flag")
	flags.All = flag.Bool("all", false, "all flags")
	flag.Parse()
	mainPath, err = getPath(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	recursiveTree(mainPath, flags)
}
