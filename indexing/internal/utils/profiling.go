package utils

import (
	"fmt"
	"log"
	"net/http"
)

func StartPprofServer() {
	fmt.Println("Starting pprof server on http://localhost:6060")

	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		log.Fatalf("error starting pprof server: %v", err)
	}
}
