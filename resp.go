package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING = '+'
	ERROR = '-'
	INTEGER = ':'
	BULK = '$'
	ARRAY = '*'
)

type Value struct {
	typ string
	str string
	num int
	bulk string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) ReadLine() (line []byte, n int, err error) {
	fmt.Println("READ LINE")
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		if b == '\r' {
			continue
		}
		if b == '\n' {
			break
		}
		line = append(line, b)
		n++
	}

	fmt.Println("READ LINE", string(line))
	return
}

func (r *Resp) ReadInteger() (x int, n int, err error) {
	fmt.Println("READ INTEGER")
	line, n, err := r.ReadLine()
	if err != nil {
		return 0, 0, err
	}

	x, err = strconv.Atoi(string(line))
	if err != nil {
		return 0, 0, err
	}

	fmt.Println("INTEGER", x)
	return
}

func (r *Resp) readArray() (Value, error) {
	v := Value{
		typ: "array",
	}

	len, _, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("ARRAY LEN", len)
	for i := 0; i < len; i++ {
		fmt.Println("ARRAY ITER", i)
		val, err := r.Read()
		if err != nil {
			return Value{}, err
		}
		v.array = append(v.array, val)
		fmt.Println("ARRAY APPEND", val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{
		typ: "bulk",
	}

	len, _, err := r.ReadInteger()
	if err != nil {
		return Value{}, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)
	v.bulk = string(bulk)
	r.ReadLine()

	return v, nil

}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("TYPE", string(_type))
	switch _type {
		case ARRAY:
			fmt.Println("ARRAY")
			return r.readArray()
		case BULK:
			fmt.Println("BULK")
			return r.readBulk()
		default:
			fmt.Printf("Unknown type: %v", string(_type))
			return Value{}, nil
	}
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, '\r', '\n')

	for _, val := range v.array {
		bytes = append(bytes, val.Marshal()...)
	}

	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

func NewWriter(wr io.Writer) *Writer {
	return &Writer{
		writer: wr,
	}
}

func (w *Writer) Write (v Value) error {
	bytes := v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
