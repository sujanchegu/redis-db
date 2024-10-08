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

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		fmt.Println("NEW-----------------------------------------------")
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

		fmt.Println("Received data:", val)

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

		command := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid request, command not found")
			writer.Write(Value{typ: "error", str: "Invalid request, command not found"})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(val)
		}

		result := handler(args)

		writer.Write(result)
	}

}