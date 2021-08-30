package main

import (
//"fmt"
   "github.com/BGrewell/go-conversions"
//"github.com/BGrewell/go-iperf"
//"time"
	"fmt"
<<<<<<< HEAD
	"github.com/BGrewell/go-conversions"
=======
>>>>>>> 328913249f87399ed1ce133fec58df85a24aa9b0
	"github.com/BGrewell/go-iperf"
	"time"
)

func main() {

	s := iperf.NewServer()
	c := iperf.NewClient("127.0.0.1")
	c.SetIncludeServer(true)
	fmt.Println(s.Id)
	fmt.Println(c.Id)

	err := s.Start()
	if err != nil {
		fmt.Println(err)
	}

	err = c.Start()
	if err != nil {
		fmt.Println(err)
	}

	for c.Running {
		time.Sleep(1 * time.Second)
	}

	fmt.Println("stopping server")
	s.Stop()

	fmt.Printf("Client exit code: %d\n", *c.ExitCode())
	fmt.Printf("Server exit code: %d\n", *s.ExitCode)
	iperf.Cleanup()
	if c.Report().Error != "" {
		fmt.Println(c.Report().Error)
	} else {
		fmt.Printf("Recv Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumReceived.BitsPerSecond)))
		fmt.Printf("Send Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumSent.BitsPerSecond)))
	}

}
