package oauth2provider

import (
	"net/http"

	"github.com/google/uuid"

	"cdr.dev/slog/v3"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/audit"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/database/dbtime"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpapi"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabd/httpmw"
	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// GetAppSecrets returns an http.HandlerFunc that handles GET /oauth2-provider/apps/{app}/secrets
func GetAppSecrets(db database.Store) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		app := httpmw.OAuth2ProviderApp(r)
		dbSecrets, err := db.GetOAuth2ProviderAppSecretsByAppID(ctx, app.ID)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error getting OAuth2 client secrets.",
				Detail:  err.Error(),
			})
			return
		}
		secrets := []optimus-ide-collabsdk.OAuth2ProviderAppSecret{}
		for _, secret := range dbSecrets {
			secrets = append(secrets, optimus-ide-collabsdk.OAuth2ProviderAppSecret{
				ID:                    secret.ID,
				LastUsedAt:            optimus-ide-collabsdk.NullTime{NullTime: secret.LastUsedAt},
				ClientSecretTruncated: secret.DisplaySecret,
			})
		}
		httpapi.Write(ctx, rw, http.StatusOK, secrets)
	}
}

// CreateAppSecret returns an http.HandlerFunc that handles POST /oauth2-provider/apps/{app}/secrets
func CreateAppSecret(db database.Store, auditor *audit.Auditor, logger slog.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			ctx               = r.Context()
			app               = httpmw.OAuth2ProviderApp(r)
			aReq, commitAudit = audit.InitRequest[database.OAuth2ProviderAppSecret](rw, &audit.RequestParams{
				Audit:   *auditor,
				Log:     logger,
				Request: r,
				Action:  database.AuditActionCreate,
			})
		)
		defer commitAudit()
		secret, err := GenerateSecret()
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Failed to generate OAuth2 client secret.",
				Detail:  err.Error(),
			})
			return
		}
		dbSecret, err := db.InsertOAuth2ProviderAppSecret(ctx, database.InsertOAuth2ProviderAppSecretParams{
			ID:           uuid.New(),
			CreatedAt:    dbtime.Now(),
			SecretPrefix: []byte(secret.Prefix),
			HashedSecret: secret.Hashed,
			// DisplaySecret is the last six characters of the original unhashed secret.
			// This is done so they can be differentiated and it matches how GitHub
			// displays their client secrets.
			DisplaySecret: secret.Formatted[len(secret.Formatted)-6:],
			AppID:         app.ID,
		})
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error creating OAuth2 client secret.",
				Detail:  err.Error(),
			})
			return
		}
		aReq.New = dbSecret
		httpapi.Write(ctx, rw, http.StatusCreated, optimus-ide-collabsdk.OAuth2ProviderAppSecretFull{
			ID:               dbSecret.ID,
			ClientSecretFull: secret.Formatted,
		})
	}
}

// DeleteAppSecret returns an http.HandlerFunc that handles DELETE /oauth2-provider/apps/{app}/secrets/{secretID}
func DeleteAppSecret(db database.Store, auditor *audit.Auditor, logger slog.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var (
			ctx               = r.Context()
			secret            = httpmw.OAuth2ProviderAppSecret(r)
			aReq, commitAudit = audit.InitRequest[database.OAuth2ProviderAppSecret](rw, &audit.RequestParams{
				Audit:   *auditor,
				Log:     logger,
				Request: r,
				Action:  database.AuditActionDelete,
			})
		)
		aReq.Old = secret
		defer commitAudit()
		err := db.DeleteOAuth2ProviderAppSecretByID(ctx, secret.ID)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, optimus-ide-collabsdk.Response{
				Message: "Internal error deleting OAuth2 client secret.",
				Detail:  err.Error(),
			})
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
