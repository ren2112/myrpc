package protocol

import (
	"fmt"
	"net/http"
)

type HttpServer struct {
}

func (server *HttpServer) Start(host string, port int) {
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Printf("Starting server at %s\n", address)

	http.HandleFunc("/", new(DispatcherServlet).service)

	if err := http.ListenAndServe(address, nil); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
