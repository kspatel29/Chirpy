package main

import (
	"net/http"
)

func main() {
        // Create a new http.ServeMux
        mux := http.NewServeMux()

        // Create a new http.Server struct
        server := &http.Server{
                Addr:    ":8080",
                Handler: mux,
        }

        // Start the server
        err := server.ListenAndServe()
        if err != nil {
                panic(err)
        }
}