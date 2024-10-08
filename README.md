# Go Redis-like Server

This project implements a simple Redis-like server in Go, supporting basic Redis commands and persistence through an Append-Only File (AOF).

## Features

- TCP server listening on port 6379
- Support for essential Redis-like commands
- String and Hash data structures
- Append-Only File (AOF) persistence
- RESP (Redis Serialization Protocol) implementation

## Project Structure

- `main.go`: Entry point of the application, sets up the TCP server
- `resp.go`: Implements the RESP protocol for reading and writing Redis commands
- `handler.go`: Contains command handlers for supported Redis commands
- `aof.go`: Implements the Append-Only File (AOF) persistence mechanism

## Getting Started

1. Clone the repository
2. Run the server:
   ```
   go run *.go
   ```
3. Connect to the server using `redis-cli`:
   ```
   redis-cli
   ```

## Supported Commands

The server supports a subset of Redis commands, including basic operations for string and hash data structures. Refer to the code for the full list of implemented commands.

## Persistence

The server uses an Append-Only File (AOF) for persistence. All write operations are logged to the `aof.txt` file.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## References

Thanks to [ahmedash.dev](https://ahmedash.dev) for the inspiration!