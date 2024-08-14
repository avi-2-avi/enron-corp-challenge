package utils

import (
	"flag"
	"fmt"
	"indexing/internal/profiling"
	"log"
)

func GetRootDirectory() string {
	prof := flag.Bool("prof", false, "Start pprof server")
	flag.Parse()

	if *prof {
		go profiling.StartPprofServer()
	}

	if flag.NArg() < 1 {
		fmt.Println("Usage: ./indexer <directory_path>")
		log.Fatal("Try again with the correct directory path")
	}

	return flag.Arg(0)
}
