package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	muxserver := http.Server{
		Addr:    ":8080",
		Handler: mux}

	err := http.ListenAndServe(muxserver.Addr, muxserver.Handler)

	if err != nil {
		fmt.Println(err) //check server started correctly.
	}

}
