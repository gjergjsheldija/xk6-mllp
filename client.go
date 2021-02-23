package mllp

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type MLLP struct {
	opts Options
}

type Options struct {
	Host string
	Port int
}

func NewClient(opts *Options) *MLLP {
	return &MLLP{
		opts: Options{
			Host: opts.Host,
			Port: opts.Port,
		}}
}

// Set the given key with the given value and expiration time.
func (m *MLLP) Send(file string) {
	err := m.sendFile(file)
	if err != nil {
		fmt.Println(err, "Impossible to send file")
	}
}

const (
	mllpStart = 0x0b
	mllpEnd   = 0x1c
	mllpEnd2  = 0x0d
)

//Send sends a file over MLLP
func (m *MLLP) sendFile(file string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", m.opts.Host, m.opts.Port))
	if err != nil {
		return err
	}
	defer conn.Close()
	fileContents := m.readFile(file)

	// write the actual message
	conn.Write([]byte{mllpStart})
	fmt.Fprintf(conn, fileContents)
	conn.Write([]byte{mllpEnd})
	conn.Write([]byte{mllpEnd2})

	// read response
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		return err
	}

	return nil
}

func (m *MLLP) readFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Reading file failed:", err.Error())
		os.Exit(1)
	}
	return string(content)
}
