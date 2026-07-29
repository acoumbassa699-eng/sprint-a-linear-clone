package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/buildinfo"
)

// GetClientInfo returns the MCP client information to use when initializing MCP connections.
// This provides a consistent way for all proxy implementations to report client information.
func GetClientInfo() mcp.Implementation {
	return mcp.Implementation{
		Name:    "optimus-ide-collab/aibridge",
		Version: buildinfo.Version(),
	}
}
