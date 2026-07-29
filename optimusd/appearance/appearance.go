package appearance

import (
	"context"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type Fetcher interface {
	Fetch(ctx context.Context) (optimus-ide-collabsdk.AppearanceConfig, error)
}

type AGPLFetcher struct {
	docsURL string
}

func (f AGPLFetcher) Fetch(context.Context) (optimus-ide-collabsdk.AppearanceConfig, error) {
	return optimus-ide-collabsdk.AppearanceConfig{
		AnnouncementBanners: []optimus-ide-collabsdk.BannerConfig{},
		SupportLinks:        optimus-ide-collabsdk.DefaultSupportLinks(f.docsURL),
		DocsURL:             f.docsURL,
	}, nil
}

func NewDefaultFetcher(docsURL string) Fetcher {
	if docsURL == "" {
		docsURL = optimus-ide-collabsdk.DefaultDocsURL()
	}
	return &AGPLFetcher{
		docsURL: docsURL,
	}
}
