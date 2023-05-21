package trinity

import (
	"testing"

	"github.com/chiyoi/go/pkg/logs"
)

func TestTrinity(t *testing.T) {
	cg := Command()
	cg.FlagSet.Usage()
}

func TestRedirectStderr(t *testing.T) {
	logs.Info("nyan")
	logs.SetOutput(LogFile())
	logs.Info("nyan1")
}
