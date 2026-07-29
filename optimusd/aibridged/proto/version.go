package proto

import "github.com/optimus-ide-collab/optimus-ide-collab/v2/apiversion"

// Version history:
//
// API v1.0:
//   - Initial version. Serves the Recorder, MCPConfigurator, and Authorizer
//     services to embedded and standalone AI Gateway daemons.
//
// API v1.1:
//   - Adds the ProviderConfigurator service with the GetAIProviders unary RPC,
//     letting embedded and standalone gateways fetch provider configuration
//     over DRPC instead of reading the database directly.
//
// API v1.2:
//   - Adds the ProviderConfigurator.WatchAIProviders streaming RPC, pushing a
//     change signal to gateways so a running standalone gateway refetches its
//     provider set when the provider configuration changes.
const (
	CurrentMajor = 1
	CurrentMinor = 2
)

// VersionQueryParam is the URL query parameter the standalone AI Gateway
// uses to advertise its aibridged API version when dialing optimus-ide-collabd's serve
// endpoint, and that optimus-ide-collabd reads to negotiate compatibility.
const VersionQueryParam = "version"

// CurrentVersion is the current aibridged API version.
// Breaking changes to the aibridged API **MUST** increment CurrentMajor above.
// Non-breaking changes to the aibridged API **MUST** increment CurrentMinor
// above.
var CurrentVersion = apiversion.New(CurrentMajor, CurrentMinor)
