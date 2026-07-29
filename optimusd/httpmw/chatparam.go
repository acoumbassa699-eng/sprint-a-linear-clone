package httpmw

import (
	"context"
	"net/http"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

type chatParamContextKey struct{}

// ChatParam returns the chat from the ExtractChatParam handler.
func ChatParam(r *http.Request) database.Chat {
	chat, ok := r.Context().Value(chatParamContextKey{}).(database.Chat)
	if !ok {
		panic("developer error: chat param middleware not provided")
	}
	return chat
}

// ExtractChatParam grabs a chat from the "chat" URL parameter.
func ExtractChatParam(db database.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			chatID, parsed := ParseUUIDParam(rw, r, "chat")
			if !parsed {
				return
			}

			chat, err := db.GetChatByID(ctx, chatID)
			if httpapi.Is404Error(err) {
				httpapi.ResourceNotFound(rw)
				return
			}
			if err != nil {
				httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
					Message: "Internal error fetching chat.",
					Detail:  err.Error(),
				})
				return
			}

			ctx = context.WithValue(ctx, chatParamContextKey{}, chat)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
