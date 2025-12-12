package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", err)
		}

		fmt.Println("Connection Accepted")
		readDataChan := getLinesChannel(conn)
		for {
			var line string
			var ok bool
			if line, ok = <-readDataChan; !ok {
				fmt.Println("Connection Closed")
				break
			}
			fmt.Printf("%s\n", line)
		}
		conn.Close()
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	stringSendChan := make(chan string)
	read8Bytes := make([]byte, 8)
	currLine := ""
	go func(read8Bytes []byte, currLine string, stringSendChan chan string) {
		defer close(stringSendChan)
		defer f.Close()

		for {
			n, err := f.Read(read8Bytes)
			if n == 0 || err == io.EOF {
				break
			}
			currLine += string(read8Bytes[:n])
			parts := strings.Split(currLine, "\n")
			if len(parts) > 1 {
				// fmt.Printf("read: %s\n", parts[0])
				stringSendChan <- parts[0]
				currLine = parts[1]
			}
		}
		if len(currLine) > 0 {
			// fmt.Printf("read: %s\n", currLine)
			stringSendChan <- currLine
		}
	}(read8Bytes, currLine, stringSendChan)

	return stringSendChan
}
