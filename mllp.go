package mllp

import (
	"context"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modules"
)

// Register the extension on module initialization, available to
// import from JS as "k6/x/mllp".
func init() {
	modules.Register("k6/x/mllp", new(MLLP))
}

func (m *MLLP) XClient(ctxPtr *context.Context, opts *Options) interface{} {
	rt := common.GetRuntime(*ctxPtr)
	return common.Bind(rt, NewClient(opts), ctxPtr)
}
