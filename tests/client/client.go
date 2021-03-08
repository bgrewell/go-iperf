package main

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/BGrewell/go-iperf"
	"os"
	"time"
)

func main() {

	includeServer := true
	proto := "tcp"
	runTime := 30
	omitSec := 10
	length := "65500"

	c := iperf.NewClient("10.23.42.18")
	c.SetIncludeServer(includeServer)
	c.SetTimeSec(runTime)
	c.SetOmitSec(omitSec)
	c.SetProto((iperf.Protocol)(proto))
	c.SetLength(length)

	fmt.Printf("blockcount: %s\n", c.BlockCount())

	err := c.Start()
	if err != nil {
		fmt.Println("failed to start client")
		os.Exit(-1)
	}

	for c.Running {
		time.Sleep(100 * time.Millisecond)
	}

	if c.Report().Error != "" {
		fmt.Println(c.Report().Error)
	} else {
		for _, entry := range c.Report().End.Streams {
			fmt.Println(entry.String())
		}
		for _, entry := range c.Report().ServerOutputJson.End.Streams {
			fmt.Println(entry.String())
		}
		fmt.Printf("DL Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumReceived.BitsPerSecond)))
		fmt.Printf("UL Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumSent.BitsPerSecond)))
	}
}
