package utils

import (
	"flag"
	"log"
)

func GetRootDirectory() string {
	prof := flag.Bool("prof", false, "Start pprof server")
	flag.Parse()

	if *prof {
		go StartProfServer()
	}

	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatal("Usage: ./indexer <directory_path>")
	}
	return flag.Arg(0)
}
