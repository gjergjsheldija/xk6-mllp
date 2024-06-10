package mllp

import (
	"errors"

	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type mllpMetrics struct {
	SentBytes        *metrics.Metric
	ReceivedBytes    *metrics.Metric
	SentMessages     *metrics.Metric
	ReceivedMessages *metrics.Metric
}

// registerMetrics registers the metrics for the mqtt module in the metrics registry
func registerMetrics(vu modules.VU) (mllpMetrics, error) {
	var err error
	m := mllpMetrics{}
	env := vu.InitEnv()
	if env == nil {
		return m, errors.New("missing env")
	}
	registry := env.Registry
	if registry == nil {
		return m, errors.New("missing registry")
	}

	m.SentBytes, err = registry.NewMetric("mllp_sent_bytes", metrics.Counter)
	if err != nil {
		return m, err
	}

	m.ReceivedBytes, err = registry.NewMetric("mllp_received_bytes", metrics.Counter)
	if err != nil {
		return m, err
	}

	m.SentMessages, err = registry.NewMetric("mllp_sent_messages_count", metrics.Counter)
	if err != nil {
		return m, err
	}

	m.ReceivedMessages, err = registry.NewMetric("mllp_received_messages_count", metrics.Counter)
	if err != nil {
		return m, err
	}

	return m, nil
}
