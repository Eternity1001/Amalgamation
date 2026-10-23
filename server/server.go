package main

import (
	"fmt"
	"net"
)


func main(){
	socket, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error: ", err)
        return
	}
	defer socket.Close()
    
    for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go sendHello(conn)
	}

}


func sendHello(conn net.Conn) {
    defer conn.Close()

    // Read the browser's HTTP request
    buffer := make([]byte, 2048)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Read error:", err)
        return
    }

    fmt.Println("--- INCOMING REQUEST ---")
	request := parseHttpRequest(buffer[:n])
	route := routeRequest(*request, conn)
	// request.print(true)
    fmt.Println("\n------------------------")

    // requestString := string(buffer[:n])
    // fmt.Println("--- INCOMING REQUEST ---")
    // fmt.Println(requestString)

    // Send a valid HTTP response

	if route != 0 {
		fmt.Println("data Broken")
	}
	// response := createDefaultHttpResponse()
	// response.addBody("text/plain", []byte("Hello Froms server"))
    // conn.Write(response.convertBytes())
    // conn.Write(route.convertBytes())

}