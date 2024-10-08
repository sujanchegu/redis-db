package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	mu sync.Mutex
	rd *bufio.Reader
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening AOF file:", err)
		return nil, err
	}

	aof := &Aof{file: file, rd: bufio.NewReader(file)}

	go func () {
		for {
			time.Sleep(1 * time.Second)
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
		}
	}()

	return aof, nil
}

func (aof *Aof) Close () error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(value Value) (error) {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}