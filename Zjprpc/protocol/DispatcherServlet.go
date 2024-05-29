package protocol

import "net/http"

type DispatcherServlet struct {
}

func (d *DispatcherServlet) service(resp http.ResponseWriter, req *http.Request) {
	new(HttpServerHandler).Handler(resp, req)
}
