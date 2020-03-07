package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	extime, err := ntp.Time("www.ntp5.stratum2.ru")
	if err != nil {
		log.Fatal("error while connect to ntp server:", err)
	}
	fmt.Fprintf(os.Stdout, "current time: %v\n", time.Now())
	fmt.Fprintf(os.Stdout, "exact time: %v\n", extime)
}
