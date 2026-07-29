package agent

import (
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/healthcheck/health"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/healthsdk"
)

func (a *agent) HandleNetcheck(rw http.ResponseWriter, r *http.Request) {
	ni := a.TailnetConn().GetNetInfo()

	ifReport, err := healthsdk.RunInterfacesReport()
	if err != nil {
		httpapi.Write(r.Context(), rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to run interfaces report",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(r.Context(), rw, http.StatusOK, healthsdk.AgentNetcheckReport{
		BaseReport: healthsdk.BaseReport{
			Severity: health.SeverityOK,
		},
		NetInfo:    ni,
		Interfaces: ifReport,
	})
}
