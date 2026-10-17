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
    
    fmt.Println("Server listening on http://localhost:8080")
	
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

    requestString := string(buffer[:n])
    fmt.Println("--- INCOMING REQUEST ---")
    fmt.Println(requestString)
    fmt.Println("------------------------")

    // Send a valid HTTP response
    response := "HTTP/1.1 200 OK\r\n" +
        "Content-Type: text/plain\r\n" +
        "Content-Length: 21\r\n" +
        "Connection: close\r\n" +
        "\r\n" +
        "Hell0 fr0m Go server!"

    conn.Write([]byte(response))
}