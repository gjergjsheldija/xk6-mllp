package mllp

import "go.k6.io/k6/stats"

var (
	WriterWrites   = stats.New("mllp.bytes.sent", stats.Counter)
	WriterReceived = stats.New("mllp.bytes.received", stats.Counter)
)
