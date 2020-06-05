package main

import (
	//"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)
var timeoutError="Timeout while connecting to server"
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
	client := NewTelnetClient(net.JoinHostPort( host, port), t, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatal(timeoutError)
	}
	fmt.Fprintf(os.Stderr,"...Connected to %v\n",net.JoinHostPort( host, port))
	//ctx, cancel := context.WithCancel(context.Background())
	wg.Add(2)
	go func() {
		select {
		case _ = <-c:
			os.Exit(1)
		}
	}()
	go readRoutine(client)
	go writeRoutine(client)
	wg.Wait()
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
		fmt.Fprintf(os.Stderr,"...Connection was closed by peer\n")
		return
	}
	fmt.Fprintf(os.Stderr,"...EOF\n")
	telnetClient.Close()
}
