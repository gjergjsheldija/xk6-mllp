package mllp

import "github.com/loadimpact/k6/stats"

var (
	WriterWrites   = stats.New("mllp.bytes.sent", stats.Counter)
	WriterReceived = stats.New("mllp.bytes.received", stats.Counter)
)
