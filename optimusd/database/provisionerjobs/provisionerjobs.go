package provisionerjobs

import (
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/pubsub"
)

const EventJobPosted = "provisioner_job_posted"

type JobPosting struct {
	OrganizationID  uuid.UUID                `json:"organization_id"`
	ProvisionerType database.ProvisionerType `json:"type"`
	Tags            map[string]string        `json:"tags"`
}

func PostJob(ps pubsub.Pubsub, job database.ProvisionerJob) error {
	msg, err := json.Marshal(JobPosting{
		OrganizationID:  job.OrganizationID,
		ProvisionerType: job.Provisioner,
		Tags:            job.Tags,
	})
	if err != nil {
		return xerrors.Errorf("marshal job posting: %w", err)
	}
	err = ps.Publish(EventJobPosted, msg)
	return err
}
