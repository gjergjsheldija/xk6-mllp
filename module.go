package mllp

import (
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/mllp".
func init() {
	modules.Register("k6/x/mllp", new(RootModule))
}

// RootModule is the global module object type. It is instantiated once per test
// run and will be used to create k6/x/mllp module instances for each VU.
type RootModule struct{}

// HL7 is the k6 extension implementing Hl7
type HL7 struct {
	vu      modules.VU
	metrics mllpMetrics
	exports map[string]interface{}
	opts    Options
}

type Options struct {
	Host string
	Port int
}

// Ensure the interfaces are implemented correctly.
var (
	_ modules.Instance = &HL7{}
	_ modules.Module   = &RootModule{}
)

// NewModuleInstance implements the modules.Module interface and returns
// a new instance for each VU.
func (r *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	metrics, err := registerMetrics(vu)
	if err != nil {
		common.Throw(vu.Runtime(), err)
	}

	hl7 := &HL7{
		vu:      vu,
		metrics: metrics,
		exports: make(map[string]interface{}),
	}

	hl7.exports["Hl7"] = hl7.client

	return hl7
}

// Exports implements the modules.Instance interface and returns the exports
// of the JS module.
func (m *HL7) Exports() modules.Exports {
	return modules.Exports{
		Named: m.exports,
	}
}
