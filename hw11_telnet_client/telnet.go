package main

import (
	"io"
	"net"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=./telnet_mock.go -package=main TelnetClient

// TelnetClient ...
type TelnetClient interface {
	Connect() error
	Send() error
	Receive() error
	Close() error
}

// NewTelnetClient ...
func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TelnetClientImpl{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

// TelnetClientImpl ...
type TelnetClientImpl struct {
	address    string
	timeout    time.Duration
	in         io.ReadCloser
	out        io.Writer
	connection net.Conn
}

// Connect ...
func (t *TelnetClientImpl) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return err
	}
	t.connection = conn

	return nil
}

// Send ...
func (t *TelnetClientImpl) Send() error {
	_, err := io.Copy(t.connection, t.in)

	return err
}

// Receive ...
func (t *TelnetClientImpl) Receive() error {
	_, err := io.Copy(t.out, t.connection)

	return err
}

// Close ...
func (t *TelnetClientImpl) Close() error {
	return t.connection.Close()
}
