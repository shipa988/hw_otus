package main

import (
	"context"
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
	flag.StringVar(&timeout, "-timeout", "10s", "enter connection timeout")
}

func main() {
	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGINT)

	wg:=&sync.WaitGroup{}

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
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case _ = <-c:
			fmt.Println("signal")
			cancel()
		}


	}()


	wg.Add(2)
	go readRoutine(ctx, client)
	go writeRoutine(ctx, client)
}

func readRoutine(ctx context.Context, telnetClient TelnetClient) {
	var err error
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			err = telnetClient.Receive()
			if err != nil {
				break OUTER
			}
		}
	}
	log.Printf("Finished readRoutine")
}

func writeRoutine(ctx context.Context, telnetClient TelnetClient) {
	var err error
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			err = telnetClient.Send()
			if err != nil {
				break OUTER
			}
		}

	}
	log.Printf("Finished writeRoutine")
}
