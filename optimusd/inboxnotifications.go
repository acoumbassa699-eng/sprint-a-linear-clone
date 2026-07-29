package optimus-ide-collabd

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw/loggermw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/notifications"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/pubsub"
	markdown "github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/render"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
	"github.com/optimus-ide-collab/websocket"
)

const (
	notificationFormatMarkdown  = "markdown"
	notificationFormatPlaintext = "plaintext"
)

var fallbackIcons = map[uuid.UUID]string{
	// workspace related notifications
	notifications.TemplateWorkspaceCreated:           optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceManuallyUpdated:   optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceDeleted:           optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceAutobuildFailed:   optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceDormant:           optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceAutoUpdated:       optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceMarkedForDeletion: optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceAutostopReminder:  optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceManualBuildFailed: optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceOutOfMemory:       optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,
	notifications.TemplateWorkspaceOutOfDisk:         optimus-ide-collabsdk.InboxNotificationFallbackIconWorkspace,

	// account related notifications
	notifications.TemplateUserAccountCreated:           optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,
	notifications.TemplateUserAccountDeleted:           optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,
	notifications.TemplateUserAccountSuspended:         optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,
	notifications.TemplateUserAccountActivated:         optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,
	notifications.TemplateYourAccountSuspended:         optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,
	notifications.TemplateYourAccountActivated:         optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,
	notifications.TemplateUserRequestedOneTimePasscode: optimus-ide-collabsdk.InboxNotificationFallbackIconAccount,

	// template related notifications
	notifications.TemplateTemplateDeleted:             optimus-ide-collabsdk.InboxNotificationFallbackIconTemplate,
	notifications.TemplateTemplateDeprecated:          optimus-ide-collabsdk.InboxNotificationFallbackIconTemplate,
	notifications.TemplateWorkspaceBuildsFailedReport: optimus-ide-collabsdk.InboxNotificationFallbackIconTemplate,

	// chat related notifications
	notifications.TemplateChatAutoArchiveDigest: optimus-ide-collabsdk.InboxNotificationFallbackIconOther,
	notifications.TemplateChatShared:            optimus-ide-collabsdk.InboxNotificationFallbackIconOther,
}

func ensureNotificationIcon(notif optimus-ide-collabsdk.InboxNotification) optimus-ide-collabsdk.InboxNotification {
	if notif.Icon != "" {
		return notif
	}

	fallbackIcon, ok := fallbackIcons[notif.TemplateID]
	if !ok {
		fallbackIcon = optimus-ide-collabsdk.InboxNotificationFallbackIconOther
	}

	notif.Icon = fallbackIcon
	return notif
}

// convertInboxNotificationResponse works as a util function to transform a database.InboxNotification to optimus-ide-collabsdk.InboxNotification
func convertInboxNotificationResponse(ctx context.Context, logger slog.Logger, notif database.InboxNotification) optimus-ide-collabsdk.InboxNotification {
	convertedNotif := optimus-ide-collabsdk.InboxNotification{
		ID:         notif.ID,
		UserID:     notif.UserID,
		TemplateID: notif.TemplateID,
		Targets:    notif.Targets,
		Title:      notif.Title,
		Content:    notif.Content,
		Icon:       notif.Icon,
		Actions: func() []optimus-ide-collabsdk.InboxNotificationAction {
			var actionsList []optimus-ide-collabsdk.InboxNotificationAction
			err := json.Unmarshal([]byte(notif.Actions), &actionsList)
			if err != nil {
				logger.Error(ctx, "unmarshal inbox notification actions", slog.Error(err))
			}
			return actionsList
		}(),
		ReadAt: func() *time.Time {
			if !notif.ReadAt.Valid {
				return nil
			}
			return &notif.ReadAt.Time
		}(),
		CreatedAt: notif.CreatedAt,
	}

	return ensureNotificationIcon(convertedNotif)
}

// watchInboxNotifications watches for new inbox notifications and sends them to the client.
// The client can specify a list of target IDs to filter the notifications.
// @Summary Watch for new inbox notifications
// @ID watch-for-new-inbox-notifications
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Notifications
// @Param targets query string false "Comma-separated list of target IDs to filter notifications"
// @Param templates query string false "Comma-separated list of template IDs to filter notifications"
// @Param read_status query string false "Filter notifications by read status. Possible values: read, unread, all"
// @Param format query string false "Define the output format for notifications title and body." enums(plaintext,markdown)
// @Success 200 {object} optimus-ide-collabsdk.GetInboxNotificationResponse
// @Router /api/v2/notifications/inbox/watch [get]
func (api *API) watchInboxNotifications(rw http.ResponseWriter, r *http.Request) {
	p := httpapi.NewQueryParamParser()
	vals := r.URL.Query()

	var (
		ctx    = r.Context()
		apikey = httpmw.APIKey(r)

		targets    = p.UUIDs(vals, []uuid.UUID{}, "targets")
		templates  = p.UUIDs(vals, []uuid.UUID{}, "templates")
		readStatus = p.String(vals, "all", "read_status")
		format     = p.String(vals, notificationFormatMarkdown, "format")
		logger     = api.Logger.Named("inbox_notifications_watcher")
	)
	p.ErrorExcessParams(vals)
	if len(p.Errors) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Query parameters have invalid values.",
			Validations: p.Errors,
		})
		return
	}

	if !slices.Contains([]string{
		string(database.InboxNotificationReadStatusAll),
		string(database.InboxNotificationReadStatusRead),
		string(database.InboxNotificationReadStatusUnread),
	}, readStatus) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "starting_before query parameter should be any of 'all', 'read', 'unread'.",
		})
		return
	}

	notificationCh := make(chan optimus-ide-collabsdk.InboxNotification, 10)

	closeInboxNotificationsSubscriber, err := api.Pubsub.SubscribeWithErr(pubsub.InboxNotificationForOwnerEventChannel(apikey.UserID),
		pubsub.HandleInboxNotificationEvent(
			func(ctx context.Context, payload pubsub.InboxNotificationEvent, err error) {
				if err != nil {
					api.Logger.Error(ctx, "inbox notification event", slog.Error(err))
					return
				}

				// HandleInboxNotificationEvent cb receives all the inbox notifications - without any filters excepted the user_id.
				// Based on query parameters defined above and filters defined by the client - we then filter out the
				// notifications we do not want to forward and discard it.

				// filter out notifications that don't match the targets
				if len(targets) > 0 {
					for _, target := range targets {
						if isFound := slices.Contains(payload.InboxNotification.Targets, target); !isFound {
							return
						}
					}
				}

				// filter out notifications that don't match the templates
				if len(templates) > 0 {
					if isFound := slices.Contains(templates, payload.InboxNotification.TemplateID); !isFound {
						return
					}
				}

				// filter out notifications that don't match the read status
				if readStatus != "" {
					if readStatus == string(database.InboxNotificationReadStatusRead) {
						if payload.InboxNotification.ReadAt == nil {
							return
						}
					} else if readStatus == string(database.InboxNotificationReadStatusUnread) {
						if payload.InboxNotification.ReadAt != nil {
							return
						}
					}
				}

				// keep a safe guard in case of latency to push notifications through websocket
				select {
				case notificationCh <- ensureNotificationIcon(payload.InboxNotification):
				default:
					api.Logger.Error(ctx, "failed to push consumed notification into websocket handler, check latency")
				}
			},
		))
	if err != nil {
		api.Logger.Error(ctx, "subscribe to inbox notification event", slog.Error(err))
		return
	}
	defer closeInboxNotificationsSubscriber()

	conn, err := websocket.Accept(rw, r, nil)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to upgrade connection to websocket.",
			Detail:  err.Error(),
		})
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_ = conn.CloseRead(context.Background())

	ctx, wsNetConn := optimus-ide-collabsdk.WebsocketNetConn(ctx, conn, websocket.MessageText)
	defer wsNetConn.Close()

	ctx = api.wsWatcher.Watch(ctx, logger, conn)

	enoptimus-ide-collab := json.NewEnoptimus-ide-collab(wsNetConn)

	// Log the request immediately instead of after it completes.
	if rl := loggermw.RequestLoggerFromContext(ctx); rl != nil {
		rl.WriteLog(ctx, http.StatusAccepted)
	}

	for {
		select {
		case <-api.ctx.Done():
			return

		case <-ctx.Done():
			return

		case notif := <-notificationCh:
			unreadCount, err := api.Database.CountUnreadInboxNotificationsByUserID(ctx, apikey.UserID)
			if err != nil {
				api.Logger.Error(ctx, "failed to count unread inbox notifications", slog.Error(err))
				return
			}

			// By default, notifications are stored as markdown
			// We can change the format based on parameter if required
			if format == notificationFormatPlaintext {
				notif.Title, err = markdown.PlaintextFromMarkdown(notif.Title)
				if err != nil {
					api.Logger.Error(ctx, "failed to convert notification title to plain text", slog.Error(err))
					return
				}

				notif.Content, err = markdown.PlaintextFromMarkdown(notif.Content)
				if err != nil {
					api.Logger.Error(ctx, "failed to convert notification content to plain text", slog.Error(err))
					return
				}
			}

			if err := enoptimus-ide-collab.Encode(optimus-ide-collabsdk.GetInboxNotificationResponse{
				Notification: notif,
				UnreadCount:  int(unreadCount),
			}); err != nil {
				api.Logger.Error(ctx, "encode notification", slog.Error(err))
				return
			}
		}
	}
}

// listInboxNotifications lists the notifications for the user.
// @Summary List inbox notifications
// @ID list-inbox-notifications
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Notifications
// @Param targets query string false "Comma-separated list of target IDs to filter notifications"
// @Param templates query string false "Comma-separated list of template IDs to filter notifications"
// @Param read_status query string false "Filter notifications by read status. Possible values: read, unread, all"
// @Param starting_before query string false "ID of the last notification from the current page. Notifications returned will be older than the associated one" format(uuid)
// @Success 200 {object} optimus-ide-collabsdk.ListInboxNotificationsResponse
// @Router /api/v2/notifications/inbox [get]
func (api *API) listInboxNotifications(rw http.ResponseWriter, r *http.Request) {
	p := httpapi.NewQueryParamParser()
	vals := r.URL.Query()

	var (
		ctx    = r.Context()
		apikey = httpmw.APIKey(r)

		targets        = p.UUIDs(vals, nil, "targets")
		templates      = p.UUIDs(vals, nil, "templates")
		readStatus     = p.String(vals, "all", "read_status")
		startingBefore = p.UUID(vals, uuid.Nil, "starting_before")
	)
	p.ErrorExcessParams(vals)
	if len(p.Errors) > 0 {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message:     "Query parameters have invalid values.",
			Validations: p.Errors,
		})
		return
	}

	if !slices.Contains([]string{
		string(database.InboxNotificationReadStatusAll),
		string(database.InboxNotificationReadStatusRead),
		string(database.InboxNotificationReadStatusUnread),
	}, readStatus) {
		httpapi.Write(ctx, rw, http.StatusBadRequest, optimus-ide-collabsdk.Response{
			Message: "starting_before query parameter should be any of 'all', 'read', 'unread'.",
		})
		return
	}

	createdBefore := dbtime.Now()
	if startingBefore != uuid.Nil {
		lastNotif, err := api.Database.GetInboxNotificationByID(ctx, startingBefore)
		if err == nil {
			createdBefore = lastNotif.CreatedAt
		}
	}

	notifs, err := api.Database.GetFilteredInboxNotificationsByUserID(ctx, database.GetFilteredInboxNotificationsByUserIDParams{
		UserID:       apikey.UserID,
		Templates:    templates,
		Targets:      targets,
		ReadStatus:   database.InboxNotificationReadStatus(readStatus),
		CreatedAtOpt: createdBefore,
	})
	if err != nil {
		api.Logger.Error(ctx, "failed to get filtered inbox notifications", slog.Error(err))
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to get filtered inbox notifications.",
		})
		return
	}

	unreadCount, err := api.Database.CountUnreadInboxNotificationsByUserID(ctx, apikey.UserID)
	if err != nil {
		api.Logger.Error(ctx, "failed to count unread inbox notifications", slog.Error(err))
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to count unread inbox notifications.",
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.ListInboxNotificationsResponse{
		Notifications: func() []optimus-ide-collabsdk.InboxNotification {
			notificationsList := make([]optimus-ide-collabsdk.InboxNotification, 0, len(notifs))
			for _, notification := range notifs {
				notificationsList = append(notificationsList, convertInboxNotificationResponse(ctx, api.Logger, notification))
			}
			return notificationsList
		}(),
		UnreadCount: int(unreadCount),
	})
}

// updateInboxNotificationReadStatus changes the read status of a notification.
// @Summary Update read status of a notification
// @ID update-read-status-of-a-notification
// @Security Optimus-IDE-CollabSessionToken
// @Produce json
// @Tags Notifications
// @Param id path string true "id of the notification"
// @Success 200 {object} optimus-ide-collabsdk.Response
// @Router /api/v2/notifications/inbox/{id}/read-status [put]
func (api *API) updateInboxNotificationReadStatus(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		apikey = httpmw.APIKey(r)
	)

	notificationID, ok := httpmw.ParseUUIDParam(rw, r, "id")
	if !ok {
		return
	}

	var body optimus-ide-collabsdk.UpdateInboxNotificationReadStatusRequest
	if !httpapi.Read(ctx, rw, r, &body) {
		return
	}

	err := api.Database.UpdateInboxNotificationReadStatus(ctx, database.UpdateInboxNotificationReadStatusParams{
		ID: notificationID,
		ReadAt: func() sql.NullTime {
			if body.IsRead {
				return sql.NullTime{
					Time:  dbtime.Now(),
					Valid: true,
				}
			}

			return sql.NullTime{}
		}(),
	})
	if err != nil {
		api.Logger.Error(ctx, "failed to update inbox notification read status", slog.Error(err))
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to update inbox notification read status.",
		})
		return
	}

	unreadCount, err := api.Database.CountUnreadInboxNotificationsByUserID(ctx, apikey.UserID)
	if err != nil {
		api.Logger.Error(ctx, "failed to call count unread inbox notifications", slog.Error(err))
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to call count unread inbox notifications.",
		})
		return
	}

	updatedNotification, err := api.Database.GetInboxNotificationByID(ctx, notificationID)
	if err != nil {
		api.Logger.Error(ctx, "failed to get notification by id", slog.Error(err))
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to get notification by id.",
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, optimus-ide-collabsdk.UpdateInboxNotificationReadStatusResponse{
		Notification: convertInboxNotificationResponse(ctx, api.Logger, updatedNotification),
		UnreadCount:  int(unreadCount),
	})
}

// markAllInboxNotificationsAsRead marks as read all unread notifications for authenticated user.
// @Summary Mark all unread notifications as read
// @ID mark-all-unread-notifications-as-read
// @Security Optimus-IDE-CollabSessionToken
// @Tags Notifications
// @Success 204
// @Router /api/v2/notifications/inbox/mark-all-as-read [put]
func (api *API) markAllInboxNotificationsAsRead(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		apikey = httpmw.APIKey(r)
	)

	err := api.Database.MarkAllInboxNotificationsAsRead(ctx, database.MarkAllInboxNotificationsAsReadParams{
		UserID: apikey.UserID,
		ReadAt: sql.NullTime{Time: dbtime.Now(), Valid: true},
	})
	if err != nil {
		api.Logger.Error(ctx, "failed to mark all unread notifications as read", slog.Error(err))
		httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
			Message: "Failed to mark all unread notifications as read.",
		})
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
