package mllp

import (
	"context"
	"errors"
	"fmt"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/stats"
	"io/ioutil"
	"net"
	"os"
	"time"
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
func (m *MLLP) Send(ctx context.Context, file string) error {
	err := m.sendFile(ctx, file)
	if err != nil {
		return err
	}
	return nil
}

const (
	startBlock = '\x0b'
	endBlock   = '\x1c'
	cr         = '\x0d'
)

//Send sends a file over MLLP
func (m *MLLP) sendFile(ctx context.Context, file string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", m.opts.Host, m.opts.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	// send data
	fileContents := m.readFile(file)
	sent, err := fmt.Fprint(conn, encapsulate([]byte(fileContents)))
	if err != nil {
		return err
	}

	// read response
	reply := make([]byte, 1024)
	resp, err := conn.Read(reply)
	if err != nil {
		return err
	}

	state := lib.GetState(ctx)
	if state == nil {
		return errors.New("state is nil")
	}

	stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
		Metric: WriterWrites,
		Time:   time.Time{},
		Value:  float64(sent),
	})
	stats.PushIfNotDone(ctx, state.Samples, stats.Sample{
		Metric: WriterReceived,
		Time:   time.Time{},
		Value:  float64(resp),
	})

	return nil
}

// encapsulate adds startBlock, endBlock and cr to the message to make it hl7v2 conformant
func encapsulate(in []byte) []byte {
	out := make([]byte, len(in)+3)
	out[0] = startBlock
	for i, b := range in {
		out[i+1] = b
	}
	out[len(out)-2] = endBlock
	out[len(out)-1] = cr
	return out
}

func (m *MLLP) readFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Reading file failed:", err.Error())
		os.Exit(1)
	}
	return string(content)
}
