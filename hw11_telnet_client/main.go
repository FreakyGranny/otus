package main

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

func main() {
	var timeoutVar string

	pflag.StringVarP(&timeoutVar, "timeout", "t", "10s", "timeout for create connection")
	pflag.Parse()
	args := pflag.Args()
	if len(args) != 2 {
		println("Wrong arguments count")
		os.Exit(111)
	}
	timeout, err := time.ParseDuration(timeoutVar)
	if err != nil {
		println("unable to parse timeout")
		os.Exit(1)
	}
	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)
	if err := run(client); err != nil {
		println("bad host")
		os.Exit(1)
	}
}

func run(client TelnetClient) error {
	err := client.Connect()
	if err != nil {
		return err
	}

	println("connected ...")
	wg := &sync.WaitGroup{}
	wg.Add(2)

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)

	sendCh := make(chan struct{}, 1)

	go func() {
		err := client.Send()
		if err != nil {
			println("error while sending")
		}
		sendCh <- struct{}{}
	}()

	go func() {
		defer wg.Done()
		defer client.Close()
		select {
		case <-terminate:
			println("interraupted...")
		case <-sendCh:
			println("exiting...")
		}
	}()

	go func() {
		defer wg.Done()
		err := client.Receive()
		if err != nil {
			println("error while receiving")

			return
		}
	}()
	wg.Wait()

	return nil
}
