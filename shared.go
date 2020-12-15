package iperf

import (
	"bufio"
	"fmt"
	"io"
)

type SharedOptions struct {
	Id string
	Port *int
	Format *rune
	Interval *int
}

type DebugScanner struct {

}

func (ds *DebugScanner) Scan(buff io.ReadCloser) {
	if buff == nil {
		fmt.Println("unable to read, ReadCloser is nil")
		return
	}
	scanner := bufio.NewScanner(buff)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text)
	}
}