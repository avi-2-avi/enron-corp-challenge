package utils

import (
	"fmt"
	"net/http"
)

func StartProfServer() {
	fmt.Println("Starting pprof server on http://localhost:6060")
	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		fmt.Println("error starting pprof server:", err)
	}
}
