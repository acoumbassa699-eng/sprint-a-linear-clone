package optimus-ide-collabd

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

func TestInboxNotifications_ensureNotificationIcon(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		icon         string
		templateID   uuid.UUID
		expectedIcon string
	}{
		{"WorkspaceCreated", "", notifications.TemplateWorkspaceCreated, optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace},
		{"UserAccountCreated", "", notifications.TemplateUserAccountCreated, optimus-ide-collabsdk.InboxNotificationFallbackIconAccount},
		{"TemplateDeleted", "", notifications.TemplateTemplateDeleted, optimus-ide-collabsdk.InboxNotificationFallbackIconTemplate},
		{"ChatShared", "", notifications.TemplateChatShared, optimus-ide-collabsdk.InboxNotificationFallbackIconOther},
		{"TestNotification", "", notifications.TemplateTestNotification, optimus-ide-collabsdk.InboxNotificationFallbackIconOther},
		{"TestExistingIcon", "https://cdn.optimus-ide-collab.com/icon_notif.png", notifications.TemplateTemplateDeleted, "https://cdn.optimus-ide-collab.com/icon_notif.png"},
		{"UnknownTemplate", "", uuid.New(), optimus-ide-collabsdk.InboxNotificationFallbackIconOther},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			notif := optimus-ide-collabsdk.InboxNotification{
				ID:         uuid.New(),
				UserID:     uuid.New(),
				TemplateID: tt.templateID,
				Title:      "notification title",
				Content:    "notification content",
				Icon:       tt.icon,
				CreatedAt:  time.Now(),
			}

			notif = ensureNotificationIcon(notif)
			require.Equal(t, tt.expectedIcon, notif.Icon)
		})
	}
}
