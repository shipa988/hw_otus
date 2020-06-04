package main

import (
	//"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout string

func init() {
	flag.StringVar(&timeout, "timeout", "10s", "enter connection timeout")
}
var wg *sync.WaitGroup
func main() {
	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT)

	wg=&sync.WaitGroup{}

	args := flag.Args()
	if len(args) <= 1 {
		log.Fatal("Please enter server address and port")
	}
	host := args[0]
	port := args[1]
	t, e := time.ParseDuration(timeout)
	if e != nil {
		log.Fatal("Please enter correct timeout value")
	}
	client := NewTelnetClient(fmt.Sprintf("%s:%s", host, port), t, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatal("Connect err")
	}
	//fmt.Println("Connect to server: ",fmt.Sprintf("%s:%s", host, port))
	//ctx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case _ = <-c:
			fmt.Println("signal")
			client.Close()
			//cancel()
		}


	}()
	go readRoutine( client)
	go writeRoutine(client)
	wg.Wait()
	//client.Close()
}

func readRoutine(telnetClient TelnetClient) {
	defer wg.Done()
	err:= telnetClient.Receive()
	if err != nil {
		return
	}
}

func writeRoutine(telnetClient TelnetClient) {
	defer wg.Done()
	err:= telnetClient.Send()
	if err != nil {
		fmt.Println(err)
		fmt.Println("...Connection was closed by peer")
		return
	}

}
