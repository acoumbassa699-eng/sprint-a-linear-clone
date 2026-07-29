package httpapi

import (
	"net/textproto"
	"strings"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// StripOptimus-IDE-CollabCookies removes the session token from the cookie header provided.
func StripOptimus-IDE-CollabCookies(header string) string {
	header = textproto.TrimString(header)
	cookies := []string{}

	var part string
	for len(header) > 0 { // continue since we have rest
		part, header, _ = strings.Cut(header, ";")
		part = textproto.TrimString(part)
		if part == "" {
			continue
		}
		name, _, _ := strings.Cut(part, "=")
		if name == optimus-ide-collabsdk.SessionTokenCookie ||
			name == optimus-ide-collabsdk.OAuth2StateCookie ||
			name == optimus-ide-collabsdk.OAuth2RedirectCookie ||
			name == optimus-ide-collabsdk.PathAppSessionTokenCookie ||
			// This uses a prefix check because the subdomain cookie is unique
			// per workspace proxy and is based on a hash of the workspace proxy
			// subdomain hostname. See the workspaceapps package for more
			// details.
			strings.HasPrefix(name, optimus-ide-collabsdk.SubdomainAppSessionTokenCookie) ||
			name == optimus-ide-collabsdk.SignedAppTokenCookie {
			continue
		}
		cookies = append(cookies, part)
	}
	return strings.Join(cookies, "; ")
}
