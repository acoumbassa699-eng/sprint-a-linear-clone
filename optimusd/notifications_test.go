package optimus-ide-collabd_test

import (
	"net/http"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications/notificationstest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
	"github.com/optimus-ide-collab/serpent"
)

func createOpts(t *testing.T) *optimus-ide-collabdtest.Options {
	t.Helper()

	dt := optimus-ide-collabdtest.DeploymentValues(t)
	return &optimus-ide-collabdtest.Options{
		DeploymentValues: dt,
	}
}

func TestUpdateNotificationsSettings(t *testing.T) {
	t.Parallel()

	t.Run("Permissions denied", func(t *testing.T) {
		t.Parallel()

		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)
		anotherClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)

		// given
		expected := optimus-ide-collabsdk.NotificationsSettings{
			NotifierPaused: true,
		}

		ctx := testutil.Context(t, testutil.WaitShort)

		// when
		err := anotherClient.PutNotificationsSettings(ctx, expected)

		// then
		var sdkError *optimus-ide-collabsdk.Error
		require.Error(t, err)
		require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
		require.Equal(t, http.StatusForbidden, sdkError.StatusCode())
	})

	t.Run("Settings modified", func(t *testing.T) {
		t.Parallel()

		client := optimus-ide-collabdtest.New(t, createOpts(t))
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)

		// given
		expected := optimus-ide-collabsdk.NotificationsSettings{
			NotifierPaused: true,
		}

		ctx := testutil.Context(t, testutil.WaitShort)

		// when
		err := client.PutNotificationsSettings(ctx, expected)
		require.NoError(t, err)

		// then
		actual, err := client.GetNotificationsSettings(ctx)
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("Settings not modified", func(t *testing.T) {
		t.Parallel()

		// Empty state: notifications Settings are undefined now (default).
		client := optimus-ide-collabdtest.New(t, createOpts(t))
		_ = optimus-ide-collabdtest.CreateFirstUser(t, client)
		ctx := testutil.Context(t, testutil.WaitShort)

		// Change the state: pause notifications
		err := client.PutNotificationsSettings(ctx, optimus-ide-collabsdk.NotificationsSettings{
			NotifierPaused: true,
		})
		require.NoError(t, err)

		// Verify the state: notifications are paused.
		actual, err := client.GetNotificationsSettings(ctx)
		require.NoError(t, err)
		require.True(t, actual.NotifierPaused)

		// Change the stage again: notifications are paused.
		expected := actual
		err = client.PutNotificationsSettings(ctx, optimus-ide-collabsdk.NotificationsSettings{
			NotifierPaused: true,
		})
		require.NoError(t, err)

		// Verify the state: notifications are still paused, and there is no error returned.
		actual, err = client.GetNotificationsSettings(ctx)
		require.NoError(t, err)
		require.Equal(t, expected.NotifierPaused, actual.NotifierPaused)
	})
}

func TestNotificationPreferences(t *testing.T) {
	t.Parallel()

	t.Run("Initial state", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitSuperLong)
		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)

		// Given: a member in its initial state.
		memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)

		// When: calling the API.
		prefs, err := memberClient.GetUserNotificationPreferences(ctx, member.ID)
		require.NoError(t, err)

		// Then: no preferences will be returned.
		require.Len(t, prefs, 0)
	})

	t.Run("Insufficient permissions", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitSuperLong)
		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)

		// Given: 2 members.
		_, member1 := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)
		member2Client, _ := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)

		// When: attempting to retrieve the preferences of another member.
		_, err := member2Client.GetUserNotificationPreferences(ctx, member1.ID)

		// Then: the API should reject the request.
		var sdkError *optimus-ide-collabsdk.Error
		require.Error(t, err)
		require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
		// NOTE: ExtractUserParam gets in the way here, and returns a 400 Bad Request instead of a 403 Forbidden.
		// This is not ideal, and we should probably change this behavior.
		require.Equal(t, http.StatusNotFound, sdkError.StatusCode())
	})

	t.Run("Admin may read any users' preferences", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitSuperLong)
		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)

		// Given: a member.
		_, member := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)

		// When: attempting to retrieve the preferences of another member as an admin.
		prefs, err := api.GetUserNotificationPreferences(ctx, member.ID)

		// Then: the API should not reject the request.
		require.NoError(t, err)
		require.Len(t, prefs, 0)
	})

	t.Run("Admin may update any users' preferences", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitSuperLong)
		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)

		// Given: a member.
		memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)

		// When: attempting to modify and subsequently retrieve the preferences of another member as an admin.
		prefs, err := api.UpdateUserNotificationPreferences(ctx, member.ID, optimus-ide-collabsdk.UpdateUserNotificationPreferences{
			TemplateDisabledMap: map[string]bool{
				notifications.TemplateWorkspaceMarkedForDeletion.String(): true,
			},
		})

		// Then: the request should succeed and the user should be able to query their own preferences to see the same result.
		require.NoError(t, err)
		require.Len(t, prefs, 1)

		memberPrefs, err := memberClient.GetUserNotificationPreferences(ctx, member.ID)
		require.NoError(t, err)
		require.Len(t, memberPrefs, 1)
		require.Equal(t, prefs[0].NotificationTemplateID, memberPrefs[0].NotificationTemplateID)
		require.Equal(t, prefs[0].Disabled, memberPrefs[0].Disabled)
	})

	t.Run("Add preferences", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitSuperLong)
		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)

		// Given: a member with no preferences.
		memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)
		prefs, err := memberClient.GetUserNotificationPreferences(ctx, member.ID)
		require.NoError(t, err)
		require.Len(t, prefs, 0)

		// When: attempting to add new preferences.
		template := notifications.TemplateWorkspaceDeleted
		prefs, err = memberClient.UpdateUserNotificationPreferences(ctx, member.ID, optimus-ide-collabsdk.UpdateUserNotificationPreferences{
			TemplateDisabledMap: map[string]bool{
				template.String(): true,
			},
		})

		// Then: the returning preferences should be set as expected.
		require.NoError(t, err)
		require.Len(t, prefs, 1)
		require.Equal(t, prefs[0].NotificationTemplateID, template)
		require.True(t, prefs[0].Disabled)
	})

	t.Run("Modify preferences", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitSuperLong)
		api := optimus-ide-collabdtest.New(t, createOpts(t))
		firstUser := optimus-ide-collabdtest.CreateFirstUser(t, api)

		// Given: a member with preferences.
		memberClient, member := optimus-ide-collabdtest.CreateAnotherUser(t, api, firstUser.OrganizationID)
		prefs, err := memberClient.UpdateUserNotificationPreferences(ctx, member.ID, optimus-ide-collabsdk.UpdateUserNotificationPreferences{
			TemplateDisabledMap: map[string]bool{
				notifications.TemplateWorkspaceDeleted.String(): true,
				notifications.TemplateWorkspaceDormant.String(): true,
			},
		})
		require.NoError(t, err)
		require.Len(t, prefs, 2)

		// When: attempting to modify their preferences.
		prefs, err = memberClient.UpdateUserNotificationPreferences(ctx, member.ID, optimus-ide-collabsdk.UpdateUserNotificationPreferences{
			TemplateDisabledMap: map[string]bool{
				notifications.TemplateWorkspaceDeleted.String(): true,
				notifications.TemplateWorkspaceDormant.String(): false, // <--- this one was changed
			},
		})
		require.NoError(t, err)
		require.Len(t, prefs, 2)

		// Then: the modified preferences should be set as expected.
		var found bool
		for _, p := range prefs {
			switch p.NotificationTemplateID {
			case notifications.TemplateWorkspaceDormant:
				found = true
				require.False(t, p.Disabled)
			case notifications.TemplateWorkspaceDeleted:
				require.True(t, p.Disabled)
			}
		}
		require.True(t, found, "dormant notification preference was not found")
	})
}

func TestNotificationDispatchMethods(t *testing.T) {
	t.Parallel()

	defaultOpts := createOpts(t)
	webhookOpts := createOpts(t)
	webhookOpts.DeploymentValues.Notifications.Method = serpent.String(database.NotificationMethodWebhook)

	tests := []struct {
		name            string
		opts            *optimus-ide-collabdtest.Options
		expectedDefault string
	}{
		{
			name:            "default",
			opts:            defaultOpts,
			expectedDefault: string(database.NotificationMethodSmtp),
		},
		{
			name:            "non-default",
			opts:            webhookOpts,
			expectedDefault: string(database.NotificationMethodWebhook),
		},
	}

	var allMethods []string
	for _, nm := range database.AllNotificationMethodValues() {
		if nm == database.NotificationMethodInbox {
			continue
		}
		allMethods = append(allMethods, string(nm))
	}
	slices.Sort(allMethods)

	// nolint:paralleltest // Not since Go v1.22.
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := testutil.Context(t, testutil.WaitSuperLong)
			api := optimus-ide-collabdtest.New(t, tc.opts)
			_ = optimus-ide-collabdtest.CreateFirstUser(t, api)

			resp, err := api.GetNotificationDispatchMethods(ctx)
			require.NoError(t, err)

			slices.Sort(resp.AvailableNotificationMethods)
			require.EqualValues(t, resp.AvailableNotificationMethods, allMethods)
			require.Equal(t, tc.expectedDefault, resp.DefaultNotificationMethod)
		})
	}
}

func TestNotificationTest(t *testing.T) {
	t.Parallel()

	t.Run("OwnerCanSendTestNotification", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)

		notifyEnq := &notificationstest.FakeEnqueuer{}
		ownerClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues:      optimus-ide-collabdtest.DeploymentValues(t),
			NotificationsEnqueuer: notifyEnq,
		})

		// Given: A user with owner permissions.
		_ = optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)

		// When: They attempt to send a test notification.
		err := ownerClient.PostTestNotification(ctx)
		require.NoError(t, err)

		// Then: We expect a notification to have been sent.
		sent := notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateTestNotification))
		require.Len(t, sent, 1)
	})

	t.Run("MemberCannotSendTestNotification", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)

		notifyEnq := &notificationstest.FakeEnqueuer{}
		ownerClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues:      optimus-ide-collabdtest.DeploymentValues(t),
			NotificationsEnqueuer: notifyEnq,
		})

		// Given: A user without owner permissions.
		ownerUser := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, ownerUser.OrganizationID)

		// When: They attempt to send a test notification.
		err := memberClient.PostTestNotification(ctx)

		// Then: We expect a forbidden error with no notifications sent
		var sdkError *optimus-ide-collabsdk.Error
		require.Error(t, err)
		require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
		require.Equal(t, http.StatusForbidden, sdkError.StatusCode())

		sent := notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateTestNotification))
		require.Len(t, sent, 0)
	})
}

func TestCustomNotification(t *testing.T) {
	t.Parallel()

	t.Run("BadRequest", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)

		notifyEnq := &notificationstest.FakeEnqueuer{}
		ownerClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues:      optimus-ide-collabdtest.DeploymentValues(t),
			NotificationsEnqueuer: notifyEnq,
		})

		// Given: A member user
		ownerUser := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
		memberClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, ownerUser.OrganizationID)

		// When: The member user attempts to send a custom notification with empty title and message
		err := memberClient.PostCustomNotification(ctx, optimus-ide-collabsdk.CustomNotificationRequest{
			Content: &optimus-ide-collabsdk.CustomNotificationContent{
				Title:   "",
				Message: "",
			},
		})

		// Then: a bad request error is expected with no notifications sent
		var sdkError *optimus-ide-collabsdk.Error
		require.Error(t, err)
		require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
		require.Equal(t, http.StatusBadRequest, sdkError.StatusCode())
		require.Equal(t, "Invalid request body", sdkError.Message)

		sent := notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateTestNotification))
		require.Len(t, sent, 0)
	})

	t.Run("SystemUserNotAllowed", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)

		notifyEnq := &notificationstest.FakeEnqueuer{}
		ownerClient, db := optimus-ide-collabdtest.NewWithDatabase(t, &optimus-ide-collabdtest.Options{
			DeploymentValues:      optimus-ide-collabdtest.DeploymentValues(t),
			NotificationsEnqueuer: notifyEnq,
		})

		// Given: A system user (prebuilds system user)
		_, token := dbgen.APIKey(t, db, database.APIKey{
			UserID:    database.PrebuildsSystemUserID,
			LoginType: database.LoginTypeNone,
		})
		systemUserClient := optimus-ide-collabsdk.New(ownerClient.URL)
		systemUserClient.SetSessionToken(token)

		// When: The system user attempts to send a custom notification
		err := systemUserClient.PostCustomNotification(ctx, optimus-ide-collabsdk.CustomNotificationRequest{
			Content: &optimus-ide-collabsdk.CustomNotificationContent{
				Title:   "Custom Title",
				Message: "Custom Message",
			},
		})

		// Then: a forbidden error is expected with no notifications sent
		var sdkError *optimus-ide-collabsdk.Error
		require.Error(t, err)
		require.ErrorAsf(t, err, &sdkError, "error should be of type *optimus-ide-collabsdk.Error")
		require.Equal(t, http.StatusForbidden, sdkError.StatusCode())
		require.Equal(t, "Forbidden", sdkError.Message)

		sent := notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateTestNotification))
		require.Len(t, sent, 0)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		ctx := testutil.Context(t, testutil.WaitShort)

		notifyEnq := &notificationstest.FakeEnqueuer{}
		ownerClient := optimus-ide-collabdtest.New(t, &optimus-ide-collabdtest.Options{
			DeploymentValues:      optimus-ide-collabdtest.DeploymentValues(t),
			NotificationsEnqueuer: notifyEnq,
		})

		// Given: A member user
		ownerUser := optimus-ide-collabdtest.CreateFirstUser(t, ownerClient)
		memberClient, memberUser := optimus-ide-collabdtest.CreateAnotherUser(t, ownerClient, ownerUser.OrganizationID)

		// When: The member user attempts to send a custom notification
		err := memberClient.PostCustomNotification(ctx, optimus-ide-collabsdk.CustomNotificationRequest{
			Content: &optimus-ide-collabsdk.CustomNotificationContent{
				Title:   "Custom Title",
				Message: "Custom Message",
			},
		})
		require.NoError(t, err)

		// Then: we expect a custom notification to be sent to the member user
		sent := notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateCustomNotification))
		require.Len(t, sent, 1)
		require.Equal(t, memberUser.ID, sent[0].UserID)
		require.Len(t, sent[0].Labels, 2)
		require.Equal(t, "Custom Title", sent[0].Labels["custom_title"])
		require.Equal(t, "Custom Message", sent[0].Labels["custom_message"])
		require.Equal(t, memberUser.ID.String(), sent[0].CreatedBy)
	})
}
