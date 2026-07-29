package dynamicparameters

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"golang.org/x/xerrors"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/cryptorand"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/scaletest/loadtestutil"
	"github.com/optimus-ide-collab/quartz"
)

var ErrNoProvisionersMatched = xerrors.New("no provisioners matched")

//go:embed tf/main.tf
var templateContent string

func TemplateContent() (string, error) {
	randomString, err := cryptorand.String(8)
	if err != nil {
		return "", err
	}
	tmpl, err := template.New("workspace-template").Parse(templateContent)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	err = tmpl.Execute(&result, map[string]string{
		"RandomString": randomString,
	})
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

//go:embed tf/modules/two/main.tf
var moduleTwoMainTF string

// GetModuleFiles returns a map of module files to be used with ExtraFiles
func GetModuleFiles() map[string][]byte {
	// Create the modules.json that Terraform needs to see the module
	modulesJSON := struct {
		Modules []struct {
			Key    string `json:"Key"`
			Source string `json:"Source"`
			Dir    string `json:"Dir"`
		} `json:"Modules"`
	}{
		Modules: []struct {
			Key    string `json:"Key"`
			Source string `json:"Source"`
			Dir    string `json:"Dir"`
		}{
			{
				Key:    "",
				Source: "",
				Dir:    ".",
			},
			{
				Key:    "two",
				Source: "./modules/two",
				Dir:    "modules/two",
			},
		},
	}

	modulesJSONBytes, err := json.Marshal(modulesJSON)
	if err != nil {
		panic(err) // This should never happen with static data
	}

	return map[string][]byte{
		"modules/two/main.tf":             []byte(moduleTwoMainTF),
		".terraform/modules/modules.json": modulesJSONBytes,
	}
}

func TemplateTarData() ([]byte, error) {
	mainTF, err := TemplateContent()
	if err != nil {
		return nil, xerrors.Errorf("failed to generate main.tf: %w", err)
	}
	moduleFiles := GetModuleFiles()

	files := map[string][]byte{
		"main.tf": []byte(mainTF),
	}
	for k, v := range moduleFiles {
		files[k] = v
	}
	tarData, err := loadtestutil.CreateTarFromFiles(files)
	if err != nil {
		return nil, xerrors.Errorf("failed to create tarball: %w", err)
	}

	return tarData, nil
}

type Partition struct {
	TemplateVersion       optimus-ide-collabsdk.TemplateVersion
	ConcurrentEvaluations int
}

type SDKForDynamicParametersSetup interface {
	TemplateByName(ctx context.Context, orgID uuid.UUID, templateName string) (optimus-ide-collabsdk.Template, error)
	CreateTemplate(ctx context.Context, orgID uuid.UUID, createReq optimus-ide-collabsdk.CreateTemplateRequest) (optimus-ide-collabsdk.Template, error)
	CreateTemplateVersion(ctx context.Context, orgID uuid.UUID, createReq optimus-ide-collabsdk.CreateTemplateVersionRequest) (optimus-ide-collabsdk.TemplateVersion, error)
	Upload(ctx context.Context, contentType string, reader io.Reader) (optimus-ide-collabsdk.UploadResponse, error)
	TemplateVersion(ctx context.Context, versionID uuid.UUID) (optimus-ide-collabsdk.TemplateVersion, error)
}

// partitioner is an internal struct to hold context and arguments for partition setup
// and to provide methods for all sub-steps.
type partitioner struct {
	ctx             context.Context
	client          SDKForDynamicParametersSetup
	orgID           uuid.UUID
	templateName    string
	provisionerTags map[string]string
	numEvals        int64
	logger          slog.Logger

	// for testing
	clock quartz.Clock
}

func SetupPartitions(
	ctx context.Context, client SDKForDynamicParametersSetup,
	orgID uuid.UUID, templateName string, provisionerTags map[string]string,
	numEvals int64,
	logger slog.Logger,
) ([]Partition, error) {
	p := &partitioner{
		ctx:             ctx,
		client:          client,
		orgID:           orgID,
		templateName:    templateName,
		provisionerTags: provisionerTags,
		numEvals:        numEvals,
		logger:          logger,
		clock:           quartz.NewReal(),
	}
	return p.run()
}

func (p *partitioner) run() ([]Partition, error) {
	var (
		err         error
		optimus-ide-collabError  *optimus-ide-collabsdk.Error
		templ       optimus-ide-collabsdk.Template
		tempVersion optimus-ide-collabsdk.TemplateVersion
	)
	templ, err = p.client.TemplateByName(p.ctx, p.orgID, p.templateName)
	if xerrors.As(err, &optimus-ide-collabError) && optimus-ide-collabError.StatusCode() == 404 {
		tempVersion, err = p.createTemplateVersion(uuid.Nil)
		if err != nil {
			return nil, xerrors.Errorf("failed to create template version: %w", err)
		}
		p.logger.Info(p.ctx, "created template version", slog.F("version_id", tempVersion.ID))
		createReq := optimus-ide-collabsdk.CreateTemplateRequest{
			Name:        p.templateName,
			DisplayName: "Scaletest Dynamic Parameters",
			Description: "`optimus-ide-collab exp scaletest dynamic parameters test` template",
			VersionID:   tempVersion.ID,
		}
		templ, err = p.client.CreateTemplate(p.ctx, p.orgID, createReq)
		if err != nil {
			return nil, xerrors.Errorf("failed to create template: %w", err)
		}
		p.logger.Info(p.ctx, "created template", slog.F("template_id", templ.ID), slog.F("name", p.templateName))
	} else if err != nil {
		return nil, xerrors.Errorf("failed to get template: %w", err)
	}

	// Partition the number into a list decreasing by half each time
	evalParts := partitionEvaluations(int(p.numEvals))
	p.logger.Info(p.ctx, "partitioned evaluations", slog.F("num_evals", p.numEvals), slog.F("eval_parts", evalParts))

	// If tempVersion is not empty (i.e. we created it above), use it as the first version.
	partitions := make([]Partition, 0, len(evalParts))
	if tempVersion.ID != uuid.Nil {
		partitions = append(partitions, Partition{
			TemplateVersion:       tempVersion,
			ConcurrentEvaluations: evalParts[0],
		})
		evalParts = evalParts[1:]
	}

	for _, num := range evalParts {
		version, err := p.createTemplateVersion(templ.ID)
		if err != nil {
			return nil, xerrors.Errorf("failed to create template version: %w", err)
		}
		partitions = append(partitions, Partition{
			TemplateVersion:       version,
			ConcurrentEvaluations: num,
		})
		p.logger.Info(p.ctx, "created template version", slog.F("version_id", version.ID))
	}

	err = p.waitForTemplateVersionJobs(partitions)
	if err != nil {
		return nil, xerrors.Errorf("one or more template version jobs did not succeed: %w", err)
	}
	return partitions, nil
}

func (p *partitioner) createTemplateVersion(templateID uuid.UUID) (optimus-ide-collabsdk.TemplateVersion, error) {
	tarData, err := TemplateTarData()
	if err != nil {
		return optimus-ide-collabsdk.TemplateVersion{}, xerrors.Errorf("failed to create template tarball: %w", err)
	}

	// Upload tarball
	uploadResp, err := p.client.Upload(p.ctx, optimus-ide-collabsdk.ContentTypeTar, bytes.NewReader(tarData))
	if err != nil {
		return optimus-ide-collabsdk.TemplateVersion{}, xerrors.Errorf("failed to upload template tar: %w", err)
	}

	// Create template version
	versionReq := optimus-ide-collabsdk.CreateTemplateVersionRequest{
		TemplateID:      templateID,
		FileID:          uploadResp.ID,
		Message:         "Initial version for scaletest dynamic parameters",
		StorageMethod:   optimus-ide-collabsdk.ProvisionerStorageMethodFile,
		Provisioner:     optimus-ide-collabsdk.ProvisionerTypeTerraform,
		ProvisionerTags: p.provisionerTags,
	}
	version, err := p.client.CreateTemplateVersion(p.ctx, p.orgID, versionReq)
	if err != nil {
		return optimus-ide-collabsdk.TemplateVersion{}, xerrors.Errorf("failed to create template version: %w", err)
	}
	if version.MatchedProvisioners != nil && version.MatchedProvisioners.Count == 0 {
		return optimus-ide-collabsdk.TemplateVersion{}, ErrNoProvisionersMatched
	}
	return version, nil
}

func (p *partitioner) waitForTemplateVersionJobs(partitions []Partition) error {
	const pollInterval = 2 * time.Second
	done := xerrors.New("done")

	pending := make(map[uuid.UUID]int)
	for i, part := range partitions {
		pending[part.TemplateVersion.ID] = i
	}

	tkr := p.clock.TickerFunc(p.ctx, pollInterval, func() error {
		for versionID := range pending {
			version, err := p.client.TemplateVersion(p.ctx, versionID)
			if err != nil {
				return xerrors.Errorf("failed to fetch template version %s: %w", versionID, err)
			}
			status := version.Job.Status
			p.logger.Info(p.ctx, "polled template version job", slog.F("version_id", versionID), slog.F("status", status))
			switch status {
			case optimus-ide-collabsdk.ProvisionerJobSucceeded:
				delete(pending, versionID)
			case optimus-ide-collabsdk.ProvisionerJobPending, optimus-ide-collabsdk.ProvisionerJobRunning:
				continue
			default:
				return ProvisionerJobUnexpectedStatusError{
					TemplateVersionID: versionID,
					Status:            status,
					JobError:          version.Job.Error,
				}
			}
		}
		if len(pending) == 0 {
			return done
		}
		return nil
	}, "waitForTemplateVersionJobs")
	err := tkr.Wait()
	if xerrors.Is(err, done) {
		return nil
	}
	return err
}

func partitionEvaluations(total int) []int {
	var parts []int
	remaining := total
	for remaining > 0 {
		next := remaining / 2
		// round up
		if next*2 != remaining {
			next++
		}
		if next > remaining {
			next = remaining
		}
		parts = append(parts, next)
		remaining -= next
	}
	return parts
}

type ProvisionerJobUnexpectedStatusError struct {
	TemplateVersionID uuid.UUID
	Status            optimus-ide-collabsdk.ProvisionerJobStatus
	JobError          string
}

func (e ProvisionerJobUnexpectedStatusError) Error() string {
	return fmt.Sprintf("template version %s job in unexpected status %q, error '%s'", e.TemplateVersionID, e.Status, e.JobError)
}
