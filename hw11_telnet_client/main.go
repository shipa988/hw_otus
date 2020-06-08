package main

import (
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

var timeoutError = "Timeout while connecting to server"
var sTimeout string

func init() {
	flag.StringVar(&sTimeout, "timeout", "10s", "enter connection timeout")
}

var wg *sync.WaitGroup

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) <= 1 {
		log.Fatal("Please enter server address and port")
	}

	host := args[0]
	port := args[1]
	timeout, e := time.ParseDuration(sTimeout)
	if e != nil {
		log.Fatal("Please enter correct timeout value")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT) //ctrl+d it's EOF
	go func() {
		<-c
		os.Exit(1) //only SIGINT
	}()

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatal(timeoutError)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %v\n", net.JoinHostPort(host, port))

	wg = &sync.WaitGroup{}
	wg.Add(2)
	go readRoutine(client)
	//ох..., я думаю дело в этом.
	//В bash тестах события записи на сервер и чтения с сервера возникают для клиента одновременно и
	//когда шедулер go решает запустить первой writeRoutine, то после отправки на сервер сообщения и получения EOF
	//клиент закрывает соединение согласно ТЗ-и тогда readRoutine не успевает принять сообщение с сервера.-тест фейлится (это происходит в случайном порядке)
	//
	//*однако если  readRoutine запустить первой-тесты проходят всегда-так как сообщение успевает приняться с сервера...
	time.Sleep(time.Millisecond * 100)
	go writeRoutine(client)
	wg.Wait()
}

func readRoutine(telnetClient TelnetClient) {
	defer wg.Done()
	if e := telnetClient.Receive(); e != nil { //if server close connect this routine is exit but we wait some unsuccessful attempts to send in writeRoutine
		fmt.Fprintf(os.Stderr, "%v\n", e)
		return
	}
}

func writeRoutine(telnetClient TelnetClient) {
	defer wg.Done()
	if e := telnetClient.Send(); e != nil {
		fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n") //an error occurs if server sent ctrl + c (close) and client execute some unsuccessful attempts to send
		return
	}
	fmt.Fprintf(os.Stderr, "...EOF\n") //if client sent ctrl+d
	telnetClient.Close()
}
