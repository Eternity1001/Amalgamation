package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type HttpRequest struct {
	method  string
	path    string
	httpVersion string
	headers map[string]string
	body    []byte
}

func (request *HttpRequest) parseHttpRequest(recvRequest []byte){
	fmt.Println(string(recvRequest))
}

func parseHttpRequest(recvRequest []byte) *HttpRequest {
	request := &HttpRequest{
		headers:make(map[string]string),
		body:  []byte{},
	}

	splitedRequest := bytes.Split(recvRequest, []byte("\r\n"))
	methodPath := (bytes.Split(splitedRequest[0], []byte(" ")))

	request.method = string(methodPath[0])
	request.path = string(methodPath[1])
	request.httpVersion = string(methodPath[2])


	for _, line := range splitedRequest[1:len(splitedRequest) - 1] {
		if len(line) <= 0 {
			break
		}
		splitLines := bytes.Split(line, []byte(": "))
		request.headers[string(splitLines[0])] = string(splitLines[1])
	}

	_, exist := request.headers["Content-Length"]
	
	if exist {
		request.body = bytes.Split(recvRequest, []byte("\r\n\r\n"))[1]
	}

	return request

}

func (request *HttpRequest) print(body bool){


	var requestString strings.Builder
	requestString.WriteString(request.method)
	requestString.WriteString(" ")
	requestString.WriteString(request.path)
	requestString.WriteString(" ")
	requestString.WriteString(request.httpVersion)
	
	for key, value := range request.headers {
		requestString.WriteString(" ")
		requestString.WriteString(key)
		requestString.WriteString(": ")
		requestString.WriteString(value)
	}

	if body {
		requestString.Write(request.body)

	}
	fmt.Printf("%s", requestString.String())
}

type Router struct {
	request HttpRequest
	conn net.Conn
	functionMap map[string](map[string]func() *HttpResponse)
}


func routeRequest(request HttpRequest, conn net.Conn) *HttpResponse {

	route := &Router {
		request: request,
		functionMap: make(map[string]map[string]func() *HttpResponse),
	}
	route.functionMap["/"] = make(map[string]func() *HttpResponse)
	route.functionMap["/"]["GET"] = route.rootPath
	
	fmt.Printf("Method: %s, Path: %s", route.request.method, route.request.path)
	path, pathExist := route.functionMap[request.path]

	if !pathExist {
		fmt.Printf("happened here first")
		return route.failedRequest()
	}

	method, methodExist := path[request.method]

	if !methodExist {
		fmt.Printf("happened here second")
		return route.failedRequest()
	}
	return method()
}


func (route *Router) failedRequest() *HttpResponse {
	response := createDefaultHttpResponse()
	response.addStatus("404", "Not Found")
	body := make([]byte, 0,  1024)
	body = (fmt.Appendf(body, "Incorrect Method: %s Incorrect Path: %s", route.request.method, route.request.path))
	return response
}

func (route *Router) rootPath() *HttpResponse{


	file, err := os.Open("../html/index.html")

	if err != nil {
		fmt.Println("fatal erorr");
	}

	defer file.Close()

	fileInfo, err := file.Stat()

	if err != nil {
		fmt.Println("Error")
	}

	fileSize := fileInfo.Size()
	data 
	_,err = io.Copy()
	response := createDefaultHttpResponse()
	response.addBody("Content-Type: text/html; charset=utf-8", string(fileSize))
	return response
}


type HttpResponse struct {
	httpVersion string
	statusCode  string
	statusMsg   string
	headers     map[string]string
	body        []byte
}

func createDefaultHttpResponse() *HttpResponse {
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

func (response *HttpResponse) addStatus(statusCode string, msg string) {
	response.statusCode = statusCode
	response.statusMsg = msg
}

func (response *HttpResponse) addHeaders(key string, value string) {
	response.headers[key] = value
}

func (response *HttpResponse) addBody(dataType string, data []byte, length string) {
	response.headers["Content-Type"] = dataType
	response.body = data
	response.headers["Content-Length"] = length
}

func (response *HttpResponse) convertBytes() []byte {
	var byteResponse strings.Builder

	byteResponse.WriteString(response.httpVersion)
	byteResponse.WriteString(response.statusCode)
	byteResponse.WriteString("\r\n")

	for header, value := range response.headers {
		byteResponse.WriteString(header)
		byteResponse.WriteString(": ")
		byteResponse.WriteString(value)
		byteResponse.WriteString("\r\n")
	}
	byteResponse.WriteString("\r\n")
	byteResponse.Write(response.body)
	return []byte(byteResponse.String())
}
