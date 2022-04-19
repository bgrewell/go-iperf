package main

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/BGrewell/go-iperf"
	"os"
)

func main() {

	includeServer := true
	proto := "tcp"
	runTime := 10
	omitSec := 0
	length := "65500"

	c := iperf.NewClient("127.0.0.1")
	c.SetIncludeServer(includeServer)
	c.SetTimeSec(runTime)
	c.SetOmitSec(omitSec)
	c.SetProto((iperf.Protocol)(proto))
	c.SetLength(length)
	c.SetJSON(true)
	c.SetIncludeServer(false)
	c.SetStreams(1)

	// Uncomment the below to get live results
	//reports := c.SetModeLive()

	//go func() {
	//	for report := range reports {
	//		fmt.Println(report.String())
	//	}
	//}()

	err := c.Start()
	if err != nil {
		fmt.Println("failed to start client")
		os.Exit(-1)
	}

	// Method 1: Wait for the test to finish by pulling from the 'Done' channel which will block until something is put in or it's closed
	<-c.Done

	// Method 2: Poll the c.Running state and wait for it to be 'false'
	//for c.Running {
	//	time.Sleep(100 * time.Millisecond)
	//}

	if c.Report() != nil && c.Report().Error != "" {
		fmt.Println(c.Report().Error)
	} else if c.Report() != nil {
		fmt.Println(c.Report().String())
		fmt.Println("----------------------------------------------------------------------------")
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
