package health_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/healthcheck/health"
)

func Test_MessageURL(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name     string
		code     health.Code
		base     string
		expected string
	}{
		{"empty", "", "", "https://optimus-ide-collab.com/docs/admin/monitoring/health-check#eunknown"},
		{"default", health.CodeAccessURLFetch, "", "https://optimus-ide-collab.com/docs/admin/monitoring/health-check#eacs03"},
		{"custom docs base", health.CodeAccessURLFetch, "https://example.com/docs", "https://example.com/docs/admin/monitoring/health-check#eacs03"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			uut := health.Message{Code: tt.code}
			actual := uut.URL(tt.base)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
