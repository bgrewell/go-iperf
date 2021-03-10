package iperf

import (
	"bufio"
	"fmt"
	"io"
)

type TestMode string

const (
	MODE_JSON TestMode = "json"
	MODE_LIVE TestMode = "live"
)

type DebugScanner struct {
	Silent bool
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
		if !ds.Silent {
			fmt.Println(text)
		}
	}
}
