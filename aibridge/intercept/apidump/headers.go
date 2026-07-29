package apidump

// sensitiveRequestHeaders are headers that should be redacted from request dumps.
var sensitiveRequestHeaders = map[string]struct{}{
	"Api-Key":                             {},
	"Authorization":                       {},
	"Cookie":                              {},
	"Proxy-Authorization":                 {},
	"X-Amz-Security-Token":                {},
	"X-Api-Key":                           {},
	"X-Auth-Token":                        {},
	"X-Optimus-IDE-Collab-AI-Governance-Session-Token": {},
	"X-Optimus-IDE-Collab-AI-Governance-Token":         {},
}

// sensitiveResponseHeaders are headers that should be redacted from response dumps.
// Note: header names use Go's canonical form (http.CanonicalHeaderKey).
var sensitiveResponseHeaders = map[string]struct{}{
	"Set-Cookie":         {},
	"Www-Authenticate":   {},
	"Proxy-Authenticate": {},
}
