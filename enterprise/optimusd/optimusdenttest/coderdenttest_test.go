package optimus-ide-collabdenttest_test

import (
	"testing"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
)

func TestNew(t *testing.T) {
	t.Parallel()
	_, _ = optimus-ide-collabdenttest.New(t, nil)
}
