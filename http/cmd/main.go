package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	// Open a server
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

	defer ln.Close()

	fmt.Println("Server running on 8080")

	for {
		// Accept the requests
		conn, err := ln.Accept()
		if err != nil {
			fmt.Print("Error accepting connection. Trying again")
			continue
		}

		fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

		// Read the request data
		go handleRequests(conn)
	}
}

func handleRequests(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	req, err := reader.ReadString('\n')
	if err != nil {
		fmt.Print("Error reading req", err)
		return
	}

	req = strings.TrimRight(req, "\r\n")

	fmt.Println("--------------------------")
	fmt.Println(req)

	parts := strings.Split(req, " ")
	method := parts[0]
	route := parts[1]
	route, params, _ := strings.Cut(route, "?")

	// Headers
	headers := map[string]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	fmt.Println(headers)

	if route == "/" {
		if params != "" {
			query := strings.Split(params, "&")
			respond(method+" at route /. It has params "+strings.Join(query, ","), conn)
		}
		respond(method+" at route /", conn)
	} else if route == "/profile" {
		if method != "GET" {
			if method == "POST" {
				cl, err := strconv.Atoi(headers["Content-Length"])
				if err != nil {
					respond("Invalid Content-Length", conn, 400)
					return
				}
				body := readBody(cl, reader)
				respond(body, conn)
			}
			respond("/profile not found", conn, 404)
		}
		respond("Mehul Parekh", conn)
	} else if strings.HasPrefix(route, "/static/") {
		respondStaticFiles(conn, route)
	}

	fmt.Println("--------------------------")
}

func readBody(contentLength int, reader *bufio.Reader) string {
	bodyBytes := make([]byte, contentLength)
	io.ReadFull(reader, bodyBytes)
	return string(bodyBytes)
}

func respond(message string, conn net.Conn, status ...int) {
	writer := bufio.NewWriter(conn)

	statusCode := 200
	statusText := "OK"

	if len(status) > 0 {
		statusCode = status[0]

		if statusCode == 404 {
			statusText = "NOT FOUND"
		} else if statusCode == 400 {
			statusText = "INTERNAL SERVER ERROR"
		}

	}

	// Headers, otherwise you won't get a response
	writer.WriteString(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusText))
	writer.WriteString("Content-Type: text/plain\r\n")
	writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(message))) // length of the body string
	writer.WriteString("\r\n")                                              // end of headers

	// Actual Message
	writer.WriteString(message)
	//Unless Flush is given it won't send back to the client and the client will get exhausted waiting for the reply.
	writer.Flush()
}

func respondStaticFiles(conn net.Conn, route string) {
	filePath := "." + route

	data, err := os.ReadFile(filePath)
	if err != nil {
		respond("File doesnt exist", conn, 404)
	}

	ext := filepath.Ext(route)
	contentType := "application/octet-stream"
	switch ext {
	case ".txt":
		contentType = "text/plain"
	case ".html":
		contentType = "text/html"
	case ".css":
		contentType = "text/css"
	case ".js":
		contentType = "application/javascript"

	}

	writer := bufio.NewWriter(conn)
	writer.WriteString("HTTP/1.1 200 OK\r\n")
	writer.WriteString(fmt.Sprintf("Content-Type: %s\r\n", contentType))
	writer.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(data))) // length of the body string
	writer.WriteString("\r\n")                                           // end of headers

	// Actual Message
	writer.Write(data)
	//Unless Flush is given it won't send back to the client and the client will get exhausted waiting for the reply.
	writer.Flush()

}
