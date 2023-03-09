package wallet

import (
	"runtime"
	"testing"
)

func TestFactory_Create(t *testing.T) {
	fa := NewFactory(NewProduceCreator(10000))
	err := fa.Create(runtime.NumCPU())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("got: %d", len(fa.GetRet()))
}
