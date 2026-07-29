package optimus-ide-collabd_test

import (
	"bytes"
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/optimus-ide-collabdtest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbauthz"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbgen"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtestutil"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications/notificationstest"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/rbac/policy"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/util/ptr"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/testutil"
)

func TestChatACLSharingLifecycle(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	mAudit := audit.NewMock()
	notifyEnq := &notificationstest.FakeEnqueuer{}
	client, db := newChatClientWithDatabase(t, func(opts *optimus-ide-collabdtest.Options) {
		opts.Auditor = mAudit
		opts.NotificationsEnqueuer = notifyEnq
	})
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	_ = createChatModelConfig(t, client)

	sharedClient, sharedUser := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)
	sharedClientExp := optimus-ide-collabsdk.NewExperimentalClient(sharedClient)
	nonSharedClient, _ := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)
	nonSharedClientExp := optimus-ide-collabsdk.NewExperimentalClient(nonSharedClient)
	groupMemberClient, groupMember := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)
	groupMemberClientExp := optimus-ide-collabsdk.NewExperimentalClient(groupMemberClient)
	sharedGroup := dbgen.Group(t, db, database.Group{OrganizationID: firstUser.OrganizationID})
	dbgen.GroupMember(t, db, database.GroupMemberTable{GroupID: sharedGroup.ID, UserID: groupMember.ID})

	data := []byte("chat sharing file")
	uploaded, err := client.UploadChatFile(ctx, firstUser.OrganizationID, "text/plain", "shared.txt", bytes.NewReader(data))
	require.NoError(t, err)
	chat := createChatForSharing(ctx, t, client, firstUser.OrganizationID, "shared chat", uploaded.ID)

	_, err = sharedClientExp.GetChat(ctx, chat.ID)
	requireSDKError(t, err, http.StatusNotFound)
	_, _, err = nonSharedClientExp.GetChatFile(ctx, uploaded.ID)
	requireSDKError(t, err, http.StatusNotFound)

	err = client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedUser.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
		GroupRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedGroup.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	require.NoError(t, err)
	require.True(t, mAudit.Contains(t, database.AuditLog{
		Action:       database.AuditActionWrite,
		ResourceType: database.ResourceTypeChat,
		ResourceID:   chat.ID,
		UserID:       firstUser.UserID,
	}))
	// Only the direct user-ACL grant is notified. The group member gains
	// access but is not notified, because group grants are not expanded.
	var sent []*notificationstest.FakeNotification
	testutil.Eventually(ctx, t, func(context.Context) bool {
		sent = notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateChatShared))
		return len(sent) == 1
	}, testutil.IntervalFast)
	require.Equal(t, sharedUser.ID, sent[0].UserID)
	require.Equal(t, firstUser.UserID.String(), sent[0].CreatedBy)
	require.Equal(t, map[string]string{
		"chat_id":    chat.ID.String(),
		"chat_title": chat.Title,
		"initiator":  optimus-ide-collabdtest.FirstUserParams.Username,
	}, sent[0].Labels)
	require.Equal(t, []uuid.UUID{chat.ID}, sent[0].Targets)
	for _, notification := range sent {
		require.NotEqual(t, groupMember.ID, notification.UserID)
	}

	notifyEnq.Clear()
	err = client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedUser.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
		GroupRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedGroup.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	require.NoError(t, err)
	require.Empty(t, notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateChatShared)))

	acl, err := client.GetChatACL(ctx, chat.ID)
	require.NoError(t, err)
	require.Len(t, acl.Users, 1)
	require.Equal(t, sharedUser.ID.String(), acl.Users[0].ID.String())
	require.Equal(t, map[uuid.UUID]optimus-ide-collabsdk.ChatRole{
		sharedUser.ID: optimus-ide-collabsdk.ChatRoleRead,
	}, chatUserRoles(acl.Users))
	require.Equal(t, map[uuid.UUID]optimus-ide-collabsdk.ChatRole{
		sharedGroup.ID: optimus-ide-collabsdk.ChatRoleRead,
	}, chatGroupRoles(acl.Groups))
	require.Len(t, acl.Groups, 1)
	require.Equal(t, sharedGroup.ID.String(), acl.Groups[0].ID.String())
	require.Empty(t, acl.Groups[0].Members)
	require.Equal(t, 1, acl.Groups[0].TotalMemberCount)

	sharedACL, err := sharedClientExp.GetChatACL(ctx, chat.ID)
	require.NoError(t, err)
	require.Equal(t, chatUserRoles(acl.Users), chatUserRoles(sharedACL.Users))
	require.Equal(t, chatGroupRoles(acl.Groups), chatGroupRoles(sharedACL.Groups))
	require.Len(t, sharedACL.Groups, 1)
	require.Empty(t, sharedACL.Groups[0].Members)
	require.Equal(t, 1, sharedACL.Groups[0].TotalMemberCount)

	sharedChat, err := sharedClientExp.GetChat(ctx, chat.ID)
	require.NoError(t, err)
	require.Equal(t, chat.ID, sharedChat.ID)
	require.Equal(t, optimus-ide-collabdtest.FirstUserParams.Username, sharedChat.OwnerUsername)
	require.Equal(t, optimus-ide-collabdtest.FirstUserParams.Name, sharedChat.OwnerName)
	require.Len(t, sharedChat.Files, 1)
	require.Equal(t, uploaded.ID, sharedChat.Files[0].ID)

	messages, err := sharedClientExp.GetChatMessages(ctx, chat.ID, nil)
	require.NoError(t, err)
	require.NotEmpty(t, messages.Messages)

	got, contentType, err := sharedClientExp.GetChatFile(ctx, uploaded.ID)
	require.NoError(t, err)
	require.Contains(t, contentType, "text/plain")
	require.Equal(t, data, got)
	_, _, err = nonSharedClientExp.GetChatFile(ctx, uploaded.ID)
	requireSDKError(t, err, http.StatusNotFound)

	groupChat, err := groupMemberClientExp.GetChat(ctx, chat.ID)
	require.NoError(t, err)
	require.Equal(t, chat.ID, groupChat.ID)

	_, err = sharedClientExp.CreateChatMessage(ctx, chat.ID, optimus-ide-collabsdk.CreateChatMessageRequest{
		Content: []optimus-ide-collabsdk.ChatInputPart{{
			Type: optimus-ide-collabsdk.ChatInputPartTypeText,
			Text: "should not send",
		}},
	})
	requireSDKError(t, err, http.StatusNotFound)

	err = sharedClientExp.UpdateChat(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatRequest{
		Title: ptr.Ref("should not rename"),
	})
	requireSDKError(t, err, http.StatusNotFound)

	err = sharedClientExp.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			groupMember.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	requireSDKError(t, err, http.StatusForbidden)

	err = sharedClientExp.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			uuid.NewString(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	requireSDKError(t, err, http.StatusForbidden)

	err = client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			strings.ToUpper(firstUser.UserID.String()): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	sdkErr := requireSDKError(t, err, http.StatusBadRequest)
	require.Equal(t, "Cannot change your own chat sharing role.", sdkErr.Message)

	err = client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedUser.ID.String(): optimus-ide-collabsdk.ChatRoleDeleted,
		},
	})
	require.NoError(t, err)
	require.Empty(t, notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateChatShared)))
	_, err = sharedClientExp.GetChat(ctx, chat.ID)
	requireSDKError(t, err, http.StatusNotFound)
	_, err = groupMemberClientExp.GetChat(ctx, chat.ID)
	require.NoError(t, err)

	mAudit.ResetLogs()
	err = client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		GroupRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedGroup.ID.String(): optimus-ide-collabsdk.ChatRoleDeleted,
		},
	})
	require.NoError(t, err)
	require.Empty(t, notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateChatShared)))
	require.True(t, mAudit.Contains(t, database.AuditLog{
		Action:       database.AuditActionWrite,
		ResourceType: database.ResourceTypeChat,
		ResourceID:   chat.ID,
		UserID:       firstUser.UserID,
	}))
	_, err = groupMemberClientExp.GetChat(ctx, chat.ID)
	requireSDKError(t, err, http.StatusNotFound)
}

func TestChatACLSharingNotifiesDirectReadersOnly(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	notifyEnq := &notificationstest.FakeEnqueuer{}
	client, db := newChatClientWithDatabase(t, func(opts *optimus-ide-collabdtest.Options) {
		opts.NotificationsEnqueuer = notifyEnq
	})
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	_ = createChatModelConfig(t, client)

	// A non-owner org admin can share another user's chat via ActionShare.
	adminClient, admin := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID, rbac.ScopedRoleOrgAdmin(firstUser.OrganizationID))
	adminExp := optimus-ide-collabsdk.NewExperimentalClient(adminClient)
	_, groupMember := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)
	_, directUser := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)

	// The group contains both the sharing initiator and another member.
	group := dbgen.Group(t, db, database.Group{OrganizationID: firstUser.OrganizationID})
	dbgen.GroupMember(t, db, database.GroupMemberTable{GroupID: group.ID, UserID: admin.ID})
	dbgen.GroupMember(t, db, database.GroupMemberTable{GroupID: group.ID, UserID: groupMember.ID})

	chat := createChatForSharing(ctx, t, client, firstUser.OrganizationID, "admin shared chat")

	// Share with the group (which grants access to groupMember) and with
	// directUser via the user ACL. Only directUser should be notified.
	err := adminExp.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			directUser.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
		GroupRoles: map[string]optimus-ide-collabsdk.ChatRole{
			group.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	require.NoError(t, err)

	var sent []*notificationstest.FakeNotification
	testutil.Eventually(ctx, t, func(context.Context) bool {
		sent = notifyEnq.Sent(notificationstest.WithTemplateID(notifications.TemplateChatShared))
		return len(sent) == 1
	}, testutil.IntervalFast)
	// Only the direct user-ACL grant is notified. The group member, initiator
	// (admin), and owner (firstUser) are never notified.
	require.Equal(t, directUser.ID, sent[0].UserID)
	for _, notification := range sent {
		require.NotEqual(t, groupMember.ID, notification.UserID)
		require.NotEqual(t, admin.ID, notification.UserID)
		require.NotEqual(t, firstUser.UserID, notification.UserID)
	}
}

func TestChatACLSubChatInheritance(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	client, db := newChatClientWithDatabase(t)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	modelConfig := createChatModelConfig(t, client)
	sharedClient, sharedUser := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)
	sharedClientExp := optimus-ide-collabsdk.NewExperimentalClient(sharedClient)

	root := createChatForSharing(ctx, t, client, firstUser.OrganizationID, "root chat")
	child := dbgen.Chat(t, db, database.Chat{
		OrganizationID:    firstUser.OrganizationID,
		OwnerID:           firstUser.UserID,
		ParentChatID:      uuid.NullUUID{UUID: root.ID, Valid: true},
		LastModelConfigID: modelConfig.ID,
		Title:             "child chat",
	})

	err := client.UpdateChatACL(ctx, root.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedUser.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	require.NoError(t, err)

	sharedChild, err := sharedClientExp.GetChat(ctx, child.ID)
	require.NoError(t, err)
	require.Equal(t, child.ID, sharedChild.ID)
	require.NotNil(t, sharedChild.RootChatID)
	require.Equal(t, root.ID, *sharedChild.RootChatID)

	_, err = sharedClientExp.GetChat(ctx, root.ID)
	require.NoError(t, err)

	err = client.UpdateChatACL(ctx, child.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedUser.ID.String(): optimus-ide-collabsdk.ChatRoleDeleted,
		},
	})
	sdkErr := requireSDKError(t, err, http.StatusBadRequest)
	require.Equal(t, "Chat ACLs can only be set on root chats.", sdkErr.Message)

	_, err = client.GetChatACL(ctx, child.ID)
	sdkErr = requireSDKError(t, err, http.StatusBadRequest)
	require.Equal(t, "Chat ACLs can only be set on root chats.", sdkErr.Message)
}

func TestChatACLValidation(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	client := newChatClient(t)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	_ = createChatModelConfig(t, client)
	chat := createChatForSharing(ctx, t, client, firstUser.OrganizationID, "validation chat")
	missingUserID := uuid.New()
	missingGroupID := uuid.New()

	tests := []struct {
		name           string
		req            optimus-ide-collabsdk.UpdateChatACL
		wantValidation optimus-ide-collabsdk.ValidationError
	}{
		{
			name: "InvalidRole",
			req: optimus-ide-collabsdk.UpdateChatACL{
				UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
					uuid.NewString(): optimus-ide-collabsdk.ChatRole("write"),
				},
			},
			wantValidation: optimus-ide-collabsdk.ValidationError{
				Field:  "user_roles",
				Detail: `role "write" is not a valid chat role`,
			},
		},
		{
			name: "InvalidUserUUID",
			req: optimus-ide-collabsdk.UpdateChatACL{
				UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
					"not-a-uuid": optimus-ide-collabsdk.ChatRoleRead,
				},
			},
			wantValidation: optimus-ide-collabsdk.ValidationError{
				Field:  "user_roles",
				Detail: "not-a-uuid is not a valid UUID.",
			},
		},
		{
			name: "InvalidGroupUUID",
			req: optimus-ide-collabsdk.UpdateChatACL{
				GroupRoles: map[string]optimus-ide-collabsdk.ChatRole{
					"not-a-uuid": optimus-ide-collabsdk.ChatRoleRead,
				},
			},
			wantValidation: optimus-ide-collabsdk.ValidationError{
				Field:  "group_roles",
				Detail: "not-a-uuid is not a valid UUID.",
			},
		},
		{
			name: "MissingUser",
			req: optimus-ide-collabsdk.UpdateChatACL{
				UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
					missingUserID.String(): optimus-ide-collabsdk.ChatRoleRead,
				},
			},
			wantValidation: optimus-ide-collabsdk.ValidationError{
				Field:  "user_roles",
				Detail: "user with ID " + missingUserID.String() + " does not exist",
			},
		},
		{
			name: "MissingGroup",
			req: optimus-ide-collabsdk.UpdateChatACL{
				GroupRoles: map[string]optimus-ide-collabsdk.ChatRole{
					missingGroupID.String(): optimus-ide-collabsdk.ChatRoleRead,
				},
			},
			wantValidation: optimus-ide-collabsdk.ValidationError{
				Field:  "group_roles",
				Detail: "group with ID " + missingGroupID.String() + " does not exist",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := testutil.Context(t, testutil.WaitLong)
			err := client.UpdateChatACL(ctx, chat.ID, tt.req)
			sdkErr := requireSDKError(t, err, http.StatusBadRequest)
			require.Equal(t, "Invalid request to update chat ACL.", sdkErr.Message)
			require.Contains(t, sdkErr.Validations, tt.wantValidation)
		})
	}
}

func TestSharedReaderStreamChat(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	client, db := newChatClientWithDatabase(t)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	modelConfig := createChatModelConfig(t, client)
	sharedClient, sharedUser := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID)
	sharedClientExp := optimus-ide-collabsdk.NewExperimentalClient(sharedClient)
	chat := dbgen.Chat(t, db, database.Chat{
		OrganizationID:    firstUser.OrganizationID,
		OwnerID:           firstUser.UserID,
		LastModelConfigID: modelConfig.ID,
		Title:             "shared stream chat",
	})
	insertAssistantCostMessage(t, db, chat.ID, modelConfig.ID, 0)

	err := client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			sharedUser.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	require.NoError(t, err)

	events, closer, err := sharedClientExp.StreamChat(ctx, chat.ID, nil)
	require.NoError(t, err)
	t.Cleanup(func() { _ = closer.Close() })

	foundAssistantMessage := false
	for !foundAssistantMessage {
		select {
		case <-ctx.Done():
			require.FailNow(t, "timed out waiting for shared stream chat event")
		case event, ok := <-events:
			require.True(t, ok, "stream closed before expected event")
			require.Equal(t, chat.ID, event.ChatID)
			require.NotEqual(t, optimus-ide-collabsdk.ChatStreamEventTypeError, event.Type)
			if event.Type == optimus-ide-collabsdk.ChatStreamEventTypeMessage &&
				event.Message != nil &&
				event.Message.Role == optimus-ide-collabsdk.ChatMessageRoleAssistant {
				foundAssistantMessage = true
			}
		}
	}
	require.NoError(t, closer.Close())

	persisted, err := db.GetChatByID(dbauthz.AsSystemRestricted(ctx), chat.ID)
	require.NoError(t, err)
	require.False(t, persisted.LastReadMessageID.Valid)
}

//nolint:tparallel,paralleltest // Subtests share a single optimus-ide-collabdtest instance.
func TestListChatsSharedScope(t *testing.T) {
	t.Parallel()

	ctx := testutil.Context(t, testutil.WaitLong)
	client, db := newChatClientWithDatabase(t)
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	modelConfig := createChatModelConfig(t, client)
	viewerClient, viewer := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID, rbac.ScopedRoleAgentsAccess(firstUser.OrganizationID))
	viewerClientExp := optimus-ide-collabsdk.NewExperimentalClient(viewerClient)
	sharedChat := dbgen.Chat(t, db, database.Chat{
		OrganizationID:    firstUser.OrganizationID,
		OwnerID:           firstUser.UserID,
		LastModelConfigID: modelConfig.ID,
		Title:             "shared with viewer",
	})
	viewerChat := dbgen.Chat(t, db, database.Chat{
		OrganizationID:    firstUser.OrganizationID,
		OwnerID:           viewer.ID,
		LastModelConfigID: modelConfig.ID,
		Title:             "viewer owned",
	})
	unsharedChat := dbgen.Chat(t, db, database.Chat{
		OrganizationID:    firstUser.OrganizationID,
		OwnerID:           firstUser.UserID,
		LastModelConfigID: modelConfig.ID,
		Title:             "not shared with viewer",
	})

	err := client.UpdateChatACL(ctx, sharedChat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			viewer.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	require.NoError(t, err)

	for _, tc := range []struct {
		name     string
		opts     *optimus-ide-collabsdk.ListChatsOptions
		expected map[uuid.UUID]struct{}
		shared   map[uuid.UUID]bool
	}{
		{
			name:     "default owned only",
			expected: map[uuid.UUID]struct{}{viewerChat.ID: {}},
			shared:   map[uuid.UUID]bool{viewerChat.ID: false},
		},
		{
			name: "created by me only",
			opts: &optimus-ide-collabsdk.ListChatsOptions{
				Source: optimus-ide-collabsdk.ChatListSourceCreatedByMe,
			},
			expected: map[uuid.UUID]struct{}{viewerChat.ID: {}},
			shared:   map[uuid.UUID]bool{viewerChat.ID: false},
		},
		{
			name: "shared with me only",
			opts: &optimus-ide-collabsdk.ListChatsOptions{
				Source: optimus-ide-collabsdk.ChatListSourceSharedWithMe,
			},
			expected: map[uuid.UUID]struct{}{sharedChat.ID: {}},
			shared:   map[uuid.UUID]bool{sharedChat.ID: true},
		},
		{
			name: "created by me and shared with me",
			opts: &optimus-ide-collabsdk.ListChatsOptions{
				Query: "source:created_by_me,shared_with_me",
			},
			expected: map[uuid.UUID]struct{}{viewerChat.ID: {}, sharedChat.ID: {}},
			shared:   map[uuid.UUID]bool{viewerChat.ID: false, sharedChat.ID: true},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chats, err := viewerClientExp.ListChats(ctx, tc.opts)
			require.NoError(t, err)
			require.Equal(t, tc.expected, chatIDSet(chats))
			require.NotContains(t, chatIDSet(chats), unsharedChat.ID)
			for _, chat := range chats {
				expectedShared, ok := tc.shared[chat.ID]
				require.True(t, ok, "missing shared assertion for chat %s", chat.ID)
				require.Equal(t, expectedShared, chat.Shared)
			}
		})
	}
}

//nolint:paralleltest // This test verifies a process-wide RBAC kill switch.
func TestChatSharingDisabled(t *testing.T) {
	previous := rbac.ChatACLDisabled()
	rbac.SetChatACLDisabled(false)
	rbac.ReloadBuiltinRoles(nil)
	t.Cleanup(func() {
		rbac.ReloadBuiltinRoles(nil)
		rbac.SetChatACLDisabled(previous)
	})

	ctx := testutil.Context(t, testutil.WaitLong)
	values := optimus-ide-collabdtest.DeploymentValues(t)
	values.DisableChatSharing = true
	store, pubsub := dbtestutil.NewDB(t)
	client := newChatClient(t, func(opts *optimus-ide-collabdtest.Options) {
		opts.DeploymentValues = values
		opts.Database = store
		opts.Pubsub = pubsub
	})
	firstUser := optimus-ide-collabdtest.CreateFirstUser(t, client.Client)
	modelConfig := createChatModelConfig(t, client)
	viewerClient, viewer := optimus-ide-collabdtest.CreateAnotherUser(t, client.Client, firstUser.OrganizationID, rbac.ScopedRoleAgentsAccess(firstUser.OrganizationID))
	viewerClientExp := optimus-ide-collabsdk.NewExperimentalClient(viewerClient)

	chat := dbgen.Chat(t, store, database.Chat{
		OrganizationID:    firstUser.OrganizationID,
		OwnerID:           firstUser.UserID,
		LastModelConfigID: modelConfig.ID,
		Title:             "disabled sharing",
	})
	err := store.UpdateChatACLByID(ctx, database.UpdateChatACLByIDParams{
		ID: chat.ID,
		UserACL: database.ChatACL{
			viewer.ID.String(): database.ChatACLEntry{Permissions: []policy.Action{policy.ActionRead}},
		},
		GroupACL: database.ChatACL{},
	})
	require.NoError(t, err)

	_, err = viewerClientExp.GetChat(ctx, chat.ID)
	requireSDKError(t, err, http.StatusNotFound)

	_, err = client.GetChatACL(ctx, chat.ID)
	sdkErr := requireSDKError(t, err, http.StatusForbidden)
	require.Equal(t, "Chat sharing is disabled for this deployment.", sdkErr.Message)

	err = client.UpdateChatACL(ctx, chat.ID, optimus-ide-collabsdk.UpdateChatACL{
		UserRoles: map[string]optimus-ide-collabsdk.ChatRole{
			viewer.ID.String(): optimus-ide-collabsdk.ChatRoleRead,
		},
	})
	requireSDKError(t, err, http.StatusForbidden)

	ownerChats, err := client.ListChats(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, map[uuid.UUID]struct{}{chat.ID: {}}, chatIDSet(ownerChats))

	viewerChats, err := viewerClientExp.ListChats(ctx, nil)
	require.NoError(t, err)
	require.Empty(t, viewerChats)
}

func createChatForSharing(
	ctx context.Context,
	t *testing.T,
	client *optimus-ide-collabsdk.ExperimentalClient,
	organizationID uuid.UUID,
	text string,
	fileIDs ...uuid.UUID,
) optimus-ide-collabsdk.Chat {
	t.Helper()

	content := []optimus-ide-collabsdk.ChatInputPart{{
		Type: optimus-ide-collabsdk.ChatInputPartTypeText,
		Text: text,
	}}
	for _, fileID := range fileIDs {
		content = append(content, optimus-ide-collabsdk.ChatInputPart{
			Type:   optimus-ide-collabsdk.ChatInputPartTypeFile,
			FileID: fileID,
		})
	}
	chat, err := client.CreateChat(ctx, optimus-ide-collabsdk.CreateChatRequest{
		OrganizationID: organizationID,
		Content:        content,
	})
	require.NoError(t, err)
	return chat
}

func chatUserRoles(users []optimus-ide-collabsdk.ChatUser) map[uuid.UUID]optimus-ide-collabsdk.ChatRole {
	roles := make(map[uuid.UUID]optimus-ide-collabsdk.ChatRole, len(users))
	for _, user := range users {
		roles[user.ID] = user.Role
	}
	return roles
}

func chatGroupRoles(groups []optimus-ide-collabsdk.ChatGroup) map[uuid.UUID]optimus-ide-collabsdk.ChatRole {
	roles := make(map[uuid.UUID]optimus-ide-collabsdk.ChatRole, len(groups))
	for _, group := range groups {
		roles[group.ID] = group.Role
	}
	return roles
}

func chatIDSet(chats []optimus-ide-collabsdk.Chat) map[uuid.UUID]struct{} {
	ids := make(map[uuid.UUID]struct{}, len(chats))
	for _, chat := range chats {
		ids[chat.ID] = struct{}{}
	}
	return ids
}
