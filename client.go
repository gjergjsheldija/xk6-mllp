package mllp

import (
	"errors"
	"fmt"
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/metrics"
	"io/ioutil"
	"net"
	"os"
	"time"
)

func (m *HL7) client(c goja.ConstructorCall) *goja.Object {

	rt := m.vu.Runtime()

	var cfg Options
	err := rt.ExportTo(c.Argument(0), &cfg)
	if err != nil {
		common.Throw(rt, fmt.Errorf("HL7 constructor expect Options as it's argument: %w", err))
	}

	return rt.ToValue(&HL7{
		opts: Options{
			Host: cfg.Host,
			Port: cfg.Port,
		}}).ToObject(rt)
}

// Send Set the given key with the given value and expiration time.
func (m *HL7) Send(file string) error {
	err := m.sendFile(file)
	if err != nil {
		return err
	}
	return nil
}

// sendFile sends a file over MLLP
func (m *HL7) sendFile(file string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", m.opts.Host, m.opts.Port))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()

	// send data
	fileContents := m.readFile(file)
	sent, err := conn.Write(encapsulate([]byte(fileContents)))
	if err != nil {
		return err
	}

	// read response
	reply := make([]byte, 1024)
	resp, err := conn.Read(reply)
	if err != nil {
		return err
	}

	err = m.publishMetrics(sent, resp)
	if err != nil {
		return err
	}

	return nil
}

func (m *HL7) publishMetrics(sent int, resp int) error {
	// publish metrics
	now := time.Now()
	state := m.vu.State()
	if state == nil {
		return errors.New("state is nil")
	}

	ctx := m.vu.Context()
	if ctx == nil {
		return errors.New("context is nil")
	}
	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{Metric: m.metrics.SentMessages},
		Time:       now,
		Value:      float64(sent),
	})
	metrics.PushIfNotDone(ctx, state.Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{Metric: m.metrics.SentBytes},
		Time:       now,
		Value:      float64(resp),
	})

	return nil
}

const (
	startBlock = '\x0b'
	endBlock   = '\x1c'
	cr         = '\x0d'
)

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

// readFile reads the file to be used as a test
func (m *HL7) readFile(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Reading file failed:", err.Error())
		os.Exit(1)
	}
	return string(content)
}
