package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

)
// Convert a incoming requests to HttpRequests to process
type HttpRequest struct {
	method  string
	path    string
	httpVersion string
	headers map[string]string
	body    []byte
}

// Takes a incoming request and fill the struct HttpRequest
func parseHttpRequest(recvRequest []byte) *HttpRequest {

	// Generate the struct
	request := &HttpRequest{
		headers:make(map[string]string),
		body:  []byte{},
	}
	// Split by new line
	splitedRequest := bytes.Split(recvRequest, []byte("\r\n"))
	methodPath := (bytes.Split(splitedRequest[0], []byte(" ")))
	request.method = string(methodPath[0])
	request.path = string(methodPath[1])
	request.httpVersion = string(methodPath[2])

	// Since we got path, method, we do this to get headers
	for _, line := range splitedRequest[1:len(splitedRequest) - 1] {
		if len(line) <= 0 {
			break
		}
		splitLines := bytes.Split(line, []byte(": "))
		request.headers[string(splitLines[0])] = string(splitLines[1])
	}
	// We check if content length exist so we can get the body
	_, exist := request.headers["Content-Length"]
	
	if exist {
		request.body = bytes.Split(recvRequest, []byte("\r\n\r\n"))[1]
	}

	return request

}
// Take HttpRequest Struct and print out it details. Useful for debugging
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

// Contains functions call for request. Automatically sent the response back
type Router struct {
	request HttpRequest
	conn net.Conn
	functionMap map[string](map[string]func() int)
}


// Takes a request and Conn and send the proper response back
// 0: Success 1: Fail
func routeRequest(request HttpRequest, conn net.Conn) int {

	route := &Router {
		request: request,
		conn: conn,
		functionMap: make(map[string]map[string]func() int),
	}
	route.functionMap["/"] = make(map[string]func() int)
	route.functionMap["/"]["GET"] = route.sendHomePage
	
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

// Send a default fail response back 
func (route *Router) failedRequest() int {
	response := createDefaultHttpResponse()
	response.addStatus("404", "Not Found")
	msg :=fmt.Sprintf("Incorrect Method: %s Incorrect Path: %s", route.request.method, route.request.path)
	response.addBody("text/plain", []byte(msg), fmt.Sprintf("%d",len(msg) ))
	
	route.conn.Write(response.convertBytes())
	return 0
}

// Sent the home page for the website
func (route *Router) sendHomePage() int{

	file, err := os.Open("./html/index.html")
	if err != nil {
		fmt.Println("fatal error opening file:", err)
		return 0 
	}
	defer file.Close()

	// 2. Get file stats
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file stats:", err)
		return 0
	}

	fileSize := fileInfo.Size()
	fileSizeStr := fmt.Sprintf("%d", fileSize)
	response := createDefaultHttpResponse()
	response.addBody("Content-Type: text/html; charset=utf-8", []byte(""), fileSizeStr)
	route.conn.Write(response.convertBytes())
	_, err = io.Copy(route.conn, file)
	if err != nil {
		fmt.Println("Error streaming file:", err)
	}
	return 0
}


type HttpResponse struct {
	httpVersion string
	statusCode  string
	statusMsg   string
	headers     map[string]string
	body        []byte
}

// Create a generic http Response 
func createDefaultHttpResponse() *HttpResponse {
	response := &HttpResponse{
		httpVersion: "HTTP/1.1 ",
		statusCode:  "200 OK",
		headers:     make(map[string]string),
		body:        []byte(""),
	}
	response.headers["Content-Length"] = "0"
	response.headers["Content-Type"] = "text/plain"
	response.headers["X-Content-Type-Options"] = "nosniff header"
	return response
}

// Modify the response struct to add / change the status code
func (response *HttpResponse) addStatus(statusCode string, msg string) {
	response.statusCode = statusCode
	response.statusMsg = msg
}

// Add new header to response struct
func (response *HttpResponse) addHeaders(key string, value string) {
	response.headers[key] = value
}

// Add a body to the response struct
func (response *HttpResponse) addBody(dataType string, data []byte, length string) {
	response.headers["Content-Type"] = dataType
	response.body = data
	response.headers["Content-Length"] = length
}

// Fao
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
