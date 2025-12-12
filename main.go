package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	fileName := "./message.txt"
	file, err := os.Open(fileName)
	if err != nil {
		// Expecting fileName to be in the root directory
		panic(err)
	}
	readDataChan := getLinesChannel(file)
	for line := range readDataChan {
		// var line string
		// var ok bool
		// if line, ok = <-readDataChan; !ok {
		// 	return
		// }
		fmt.Printf("read: %s\n", line)
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
