package main

import (
	"fmt"
	"github.com/BGrewell/go-iperf"
	"os"
	"time"
)

func main() {
	s := iperf.NewServer()
	err := s.Start()
	if err != nil {
		fmt.Println("failed to start server")
		os.Exit(-1)
	}

	for s.Running {
		time.Sleep(1 * time.Second)
	}

	fmt.Println("server has exited")
}
