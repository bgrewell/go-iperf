package main

import (
	"fmt"
	"github.com/BGrewell/go-iperf"
	"time"
)

func main() {
	s := iperf.NewServer()
	c := iperf.NewClient()
	fmt.Println(s.Id)
	fmt.Println(c.Id)

	err := s.Start()
	if err != nil {
		fmt.Println(err)
	}
	for s.Running {
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Server exit code: %d\n", *s.ExitCode)
	iperf.Cleanup()
}
