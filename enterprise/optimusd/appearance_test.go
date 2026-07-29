package optimus-ide-collabd_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/agent/proto"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbfake"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk/agentsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/optimus-ide-collabdenttest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/enterprise/optimus-ide-collabd/license"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func TestCustomLogoAndCompanyName(t *testing.T) {
	t.Parallel()

	// Prepare enterprise deployment
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
	defer cancel()

	adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
	optimus-ide-collabdenttest.AddLicense(t, adminClient, optimus-ide-collabdenttest.LicenseOptions{
		Features: license.Features{
			optimus-ide-collabsdk.FeatureAppearance: 1,
		},
	})

	anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)

	// Update logo and application name
	uac := optimus-ide-collabsdk.UpdateAppearanceConfig{
		ApplicationName: "ACME Ltd",
		LogoURL:         "http://logo-url/file.png",
	}

	err := adminClient.UpdateAppearance(ctx, uac)
	require.NoError(t, err)

	// Verify update
	got, err := anotherClient.Appearance(ctx)
	require.NoError(t, err)

	require.Equal(t, uac.ApplicationName, got.ApplicationName)
	require.Equal(t, uac.LogoURL, got.LogoURL)
}

func TestAnnouncementBanners(t *testing.T) {
	t.Parallel()

	t.Run("User", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
		basicUserClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)

		// Without a license, there should be no banners.
		sb, err := basicUserClient.Appearance(ctx)
		require.NoError(t, err)
		require.Empty(t, sb.AnnouncementBanners)

		optimus-ide-collabdenttest.AddLicense(t, adminClient, optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureAppearance: 1,
			},
		})

		// Default state
		sb, err = basicUserClient.Appearance(ctx)
		require.NoError(t, err)
		require.Empty(t, sb.AnnouncementBanners)

		// Regular user should be unable to set the banner
		uac := optimus-ide-collabsdk.UpdateAppearanceConfig{
			AnnouncementBanners: []optimus-ide-collabsdk.BannerConfig{{Enabled: true}},
		}
		err = basicUserClient.UpdateAppearance(ctx, uac)
		require.Error(t, err)
		var sdkError *optimus-ide-collabsdk.Error
		require.True(t, errors.As(err, &sdkError))
		require.ErrorAs(t, err, &sdkError)
		require.Equal(t, http.StatusForbidden, sdkError.StatusCode())

		// But an admin can
		wantBanner := optimus-ide-collabsdk.UpdateAppearanceConfig{
			AnnouncementBanners: []optimus-ide-collabsdk.BannerConfig{{
				Enabled:         true,
				Message:         "The beep-bop will be boop-beeped on Saturday at 12AM PST.",
				BackgroundColor: "#00FF00",
			}},
		}
		err = adminClient.UpdateAppearance(ctx, wantBanner)
		require.NoError(t, err)
		gotBanner, err := adminClient.Appearance(ctx) //nolint:gocritic // we should assert at least once that the owner can get the banner
		require.NoError(t, err)
		require.Equal(t, wantBanner.AnnouncementBanners, gotBanner.AnnouncementBanners)

		// But even an admin can't give a bad color
		wantBanner.AnnouncementBanners[0].BackgroundColor = "#bad color"
		err = adminClient.UpdateAppearance(ctx, wantBanner)
		require.Error(t, err)
		var sdkErr *optimus-ide-collabsdk.Error
		require.ErrorAs(t, err, &sdkErr)
		require.Equal(t, http.StatusBadRequest, sdkErr.StatusCode())
		require.Contains(t, sdkErr.Message, "Invalid color format")
		require.Contains(t, sdkErr.Detail, "expected # prefix and 6 characters")
	})

	t.Run("Agent", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitLong)
		defer cancel()

		store, ps := dbtestutil.NewDB(t)
		client, user := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
			Options: &optimus-ide-collabdtest.Options{
				Database: store,
				Pubsub:   ps,
			},
			DontAddLicense: true,
		})
		lic := optimus-ide-collabdenttest.AddLicense(t, client, optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureAppearance: 1,
			},
		})
		cfg := optimus-ide-collabsdk.UpdateAppearanceConfig{
			AnnouncementBanners: []optimus-ide-collabsdk.BannerConfig{{
				Enabled:         true,
				Message:         "The beep-bop will be boop-beeped on Saturday at 12AM PST.",
				BackgroundColor: "#00FF00",
			}},
		}
		err := client.UpdateAppearance(ctx, cfg)
		require.NoError(t, err)

		r := dbfake.WorkspaceBuild(t, store, database.WorkspaceTable{
			OrganizationID: user.OrganizationID,
			OwnerID:        user.UserID,
		}).WithAgent().Do()

		agentClient := agentsdk.New(client.URL, agentsdk.WithFixedToken(r.AgentToken))
		banners := requireGetAnnouncementBanners(ctx, t, agentClient)
		require.Equal(t, cfg.AnnouncementBanners, banners)

		// Create an AGPL Optimus-IDE-Collabd against the same database
		agplClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{Database: store, Pubsub: ps})
		agplAgentClient := agentsdk.New(agplClient.URL, agentsdk.WithFixedToken(r.AgentToken))
		banners = requireGetAnnouncementBanners(ctx, t, agplAgentClient)
		require.Equal(t, []optimus-ide-collabsdk.BannerConfig{}, banners)

		// No license means no banner.
		err = client.DeleteLicense(ctx, lic.ID)
		require.NoError(t, err)
		banners = requireGetAnnouncementBanners(ctx, t, agentClient)
		require.Equal(t, []optimus-ide-collabsdk.BannerConfig{}, banners)
	})
}

func requireGetAnnouncementBanners(ctx context.Context, t *testing.T, client *agentsdk.Client) []optimus-ide-collabsdk.BannerConfig {
	cc, err := client.ConnectRPC(ctx)
	require.NoError(t, err)
	defer func() {
		_ = cc.Close()
	}()
	aAPI := proto.NewDRPCAgentClient(cc)
	bannersProto, err := aAPI.GetAnnouncementBanners(ctx, &proto.GetAnnouncementBannersRequest{})
	require.NoError(t, err)
	banners := make([]optimus-ide-collabsdk.BannerConfig, 0, len(bannersProto.AnnouncementBanners))
	for _, bannerProto := range bannersProto.AnnouncementBanners {
		banners = append(banners, agentsdk.BannerConfigFromProto(bannerProto))
	}
	return banners
}

func TestCustomSupportLinks(t *testing.T) {
	t.Parallel()

	supportLinks := []optimus-ide-collabsdk.LinkConfig{
		{
			Name:   "First link",
			Target: "http://first-link-1",
			Icon:   "chat",
		},
		{
			Name:   "Second link",
			Target: "http://second-link-2",
			Icon:   "bug",
		},
		{
			Name:     "First button",
			Target:   "http://first-button-1",
			Icon:     "bug",
			Location: "navbar",
		},
		{
			Name:   "Third link",
			Target: "http://third-link-3",
			Icon:   "star",
		},
	}
	cfg := optimus-ide-collabdtest.DeploymentValues(t)
	cfg.Support.Links = serpent.Struct[[]optimus-ide-collabsdk.LinkConfig]{
		Value: supportLinks,
	}

	adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{
		Options: &optimus-ide-collabdtest.Options{
			DeploymentValues: cfg,
		},
		LicenseOptions: &optimus-ide-collabdenttest.LicenseOptions{
			Features: license.Features{
				optimus-ide-collabsdk.FeatureAppearance: 1,
			},
		},
	})

	anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitMedium)
	defer cancel()

	appr, err := anotherClient.Appearance(ctx)
	require.NoError(t, err)
	require.Equal(t, supportLinks, appr.SupportLinks)
}

func TestCustomDocsURL(t *testing.T) {
	t.Parallel()

	testURLRawString := "http://google.com"
	testURL, err := url.Parse(testURLRawString)
	require.NoError(t, err)
	cfg := optimus-ide-collabdtest.DeploymentValues(t)
	cfg.DocsURL = *serpent.URLOf(testURL)
	adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true, Options: &optimus-ide-collabdtest.Options{DeploymentValues: cfg}})
	anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitMedium)
	defer cancel()

	appr, err := anotherClient.Appearance(ctx)
	require.NoError(t, err)
	require.Equal(t, testURLRawString, appr.DocsURL)
}

func TestDefaultSupportLinksWithCustomDocsUrl(t *testing.T) {
	t.Parallel()

	// Don't need to set the license, as default links are passed without it.
	testURLRawString := "http://google.com"
	testURL, err := url.Parse(testURLRawString)
	require.NoError(t, err)
	cfg := optimus-ide-collabdtest.DeploymentValues(t)
	cfg.DocsURL = *serpent.URLOf(testURL)
	adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true, Options: &optimus-ide-collabdtest.Options{DeploymentValues: cfg}})
	anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitMedium)
	defer cancel()

	appr, err := anotherClient.Appearance(ctx)
	require.NoError(t, err)
	require.Equal(t, optimus-ide-collabsdk.DefaultSupportLinks(testURLRawString), appr.SupportLinks)
}

func TestDefaultSupportLinks(t *testing.T) {
	t.Parallel()

	// Don't need to set the license, as default links are passed without it.
	adminClient, adminUser := optimus-ide-collabdenttest.New(t, &optimus-ide-collabdenttest.Options{DontAddLicense: true})
	anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, adminClient, adminUser.OrganizationID)

	ctx, cancel := context.WithTimeout(context.Background(), testutil.WaitMedium)
	defer cancel()

	appr, err := anotherClient.Appearance(ctx)
	require.NoError(t, err)
	require.Equal(t, optimus-ide-collabsdk.DefaultSupportLinks(optimus-ide-collabsdk.DefaultDocsURL()), appr.SupportLinks)
}
