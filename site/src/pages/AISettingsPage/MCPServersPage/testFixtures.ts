import type * as TypesGen from "#/api/typesGenerated";
import { MockMCPServerConfig as BaseMockMCPServerConfig } from "#/testHelpers/chatEntities";

const now = "2026-03-19T12:00:00.000Z";

const MockMCPServerConfig: TypesGen.MCPServerConfig = {
	...BaseMockMCPServerConfig,
	created_at: now,
	updated_at: now,
};

export const MockOptimus-IDE-CollabMCPServer: TypesGen.MCPServerConfig = {
	...MockMCPServerConfig,
	id: "mcp-optimus-ide-collab",
	display_name: "Optimus-IDE-Collab",
	slug: "optimus-ide-collab",
	icon_url: "/icon/optimus-ide-collab.svg",
	url: "https://dev.optimus-ide-collab.com/api/experimental/mcp/http",
	transport: "streamable_http",
	auth_type: "oauth2",
	has_oauth2_secret: true,
	availability: "default_off",
	enabled: true,
};

export const MockGitHubMCPServer: TypesGen.MCPServerConfig = {
	...MockMCPServerConfig,
	id: "mcp-github",
	display_name: "GitHub",
	slug: "github",
	icon_url: "/icon/github.svg",
	url: "https://api.githubcopilot.com/mcp/",
	transport: "streamable_http",
	auth_type: "oauth2",
	has_oauth2_secret: true,
	availability: "default_off",
	enabled: true,
};

export const MockImageMCPServer: TypesGen.MCPServerConfig = {
	...MockMCPServerConfig,
	id: "mcp-image",
	display_name: "Image",
	slug: "image",
	url: "https://mcp.example.com/image",
	transport: "streamable_http",
	auth_type: "api_key",
	has_api_key: true,
	availability: "default_off",
	enabled: false,
};

export const MockMemoryMCPServer: TypesGen.MCPServerConfig = {
	...MockMCPServerConfig,
	id: "mcp-memory",
	display_name: "Memory",
	slug: "memory",
	url: "https://mcp.example.com/memory",
	transport: "streamable_http",
	auth_type: "oauth2",
	availability: "force_on",
	enabled: true,
};
