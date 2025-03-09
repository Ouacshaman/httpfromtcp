package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt")
	defer file.Close()
	if err != nil {
		fmt.Printf("Unable to Open File: %v", err)
		return
	}
	ch := getLinesChannel(file)
	for {
		v, ok := <-ch
		if !ok {
			break
		}
		fmt.Printf("read: %s\n", v)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		/*
			channel is closed inside as the go routine is the sender
			The goroutine handles everything related to sending values
			to the channel.
			It doesn't rely on external code to manage the channel's lifecycle.
			This will ensure channel does not close prematurely
			This will ensure the channel close after the go routine is processed and no data is sent to a closed channel
		*/
		defer close(ch)
		res := ""

		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				break
			}
			parts := strings.Split(string(b[:n]), "\n")
			for _, v := range parts[:len(parts)-1] {
				res += v
				ch <- res
				res = ""
			}
			res += parts[len(parts)-1]
		}
	}()
	return ch
}
