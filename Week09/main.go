package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	closeChan := make(chan struct{})

	s := NewServer()
	go func() {
		err := s.Serve(":8000")
		if err != nil {
			fmt.Println("listen failed, addr: :8000, err: ", err)
		}
		closeChan <- struct{}{}
	}()

	select {
	case <-sigChan:
		s.Close()
	case <-closeChan:
	}

	fmt.Println("server is exit")
}
