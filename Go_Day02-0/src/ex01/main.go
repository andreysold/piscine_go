package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Flags struct {
	Words      *bool
	Lines      *bool
	Characters *bool
}

func getFiles(args []string) []string {
	var files []string
	for i := range args {
		if strings.HasSuffix(args[i], "txt") {
			files = append(files, args[i])
		}
	}
	return files
}

func filesCommand(flags *Flags, fileName string) {
	var (
		fd    *os.File
		err   error
		count int
	)

	fd, err = os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return
	}

	scanner := bufio.NewScanner(fd)
	if *flags.Words {
		scanner.Split(bufio.ScanWords)
	} else if *flags.Lines {
		scanner.Split(bufio.ScanLines)
	} else {
		scanner.Split(bufio.ScanBytes)
	}
	for scanner.Scan() {
		count++
	}

	fmt.Println(count, fileName)
}

func concurrentlyOutput(filesName *[]string, flags *Flags) {
	var wait sync.WaitGroup

	for i := range *filesName {
		wait.Add(1)
		go func(i int) {
			defer wait.Done()
			filesCommand(flags, (*filesName)[i])
		}(i)
	}
	wait.Wait()
}

func main() {
	var (
		flags     Flags
		filesName []string
	)

	flags.Words = flag.Bool("w", false, "worlds")
	flags.Lines = flag.Bool("l", false, "lines")
	flags.Characters = flag.Bool("m", false, "characters")
	flag.Parse()

	filesName = getFiles(os.Args[1:])
	concurrentlyOutput(&filesName, &flags)
}
