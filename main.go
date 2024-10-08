package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Redis server listening on port 6379")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Creating a new goroutine for each connection so that the server can handle multiple connections concurrently.
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		resp := NewResp(conn)
		writer := NewWriter(conn)
		aof, err := NewAof("aof.txt")
		if err != nil {
			fmt.Println("Error creating AOF:", err)
			return
		}

		val, err := resp.Read()
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}

		// Comment this out to see the raw data
		// fmt.Println("Received data:", val)

		// The type is expected to be an array because commands are sent over RESP as arrays.
		if val.typ != "array" {
			fmt.Println("Invalid request, array expected")
			writer.Write(Value{typ: "error", str: "Invalid request, array expected"})
			continue
		}

		if len(val.array) == 0 {
			fmt.Println("Invalid request, not enough args")
			writer.Write(Value{typ: "error", str: "Invalid request, not enough args"})
			continue
		}

		// RESP stores the command as a bulk string.
		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			if command != "COMMAND" && command != "INFO" {
				fmt.Printf("Invalid request, command not found: %s\n", command)
			}
			writer.Write(Value{typ: "error", str: fmt.Sprintf("Invalid request, command not found: %s", command)})
			continue
		}

		if command == "SET" || command == "HSET" {
			// Adding only write commands to the AOF file, as read commands doesn't change state of data.
			aof.Write(val)
		}

		result := handler(args)

		writer.Write(result)
	}

}