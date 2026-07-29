package chatadvisor_test

import (
	"testing"

	"go.uber.org/goleak"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m, testutil.GoleakOptions...)
}
