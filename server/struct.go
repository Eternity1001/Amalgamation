package main

type HttpRequest struct {
	method  string
	path    string
	headers map[string]string
	body    []byte
}

type HttpResponse struct {
	httpVersion string
	statusCode  string
	statusMsg   string
	headers     map[string]string
	body        []byte
}

func createHttpResponse() *HttpResponse {
	response := &HttpResponse{
		httpVersion: "HTTP/1.1 ",
		statusCode:  "200 OK",
		headers:     make(map[string]string),
		body:        []byte{},
	}
	response.headers["Content-Length"] = "0"
	response.headers["Content-Type"] = "text/plain"
	response.headers["X-Content-Type-Options"] = "nosniff header"
	return response
}

func (response *HttpResponse) createStatus(statusCode string, msg string) {
	response.statusCode = statusCode
	response.statusMsg = msg
}
