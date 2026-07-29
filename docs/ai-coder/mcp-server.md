# MCP Server

Optimus-IDE-Collab includes a built-in [Model Context Protocol](https://modelcontextprotocol.io/)
(MCP) server that provides AI assistants with tools and context about your Optimus-IDE-Collab
deployment. This enables AI-powered workflows for managing workspaces,
templates, and development environments.

Optimus-IDE-Collab supports two MCP server modes:

- **[Local MCP Server](#local-mcp-server)**: Runs via the Optimus-IDE-Collab CLI using stdio
  transport. Ideal for local AI tools and IDE integrations.
- **[Remote MCP Server](#remote-mcp-server)**: HTTP-based server exposed by your
  Optimus-IDE-Collab deployment. Supports OAuth2 authentication and is published to the MCP
  Registry.

## Local MCP Server

The local MCP server runs via the Optimus-IDE-Collab CLI and uses stdio transport to
communicate with AI tools.

### Setup

Run the MCP server using the Optimus-IDE-Collab CLI:

```sh
optimus-ide-collab exp mcp server
```

### Client Configuration

Configure your MCP client to spawn the Optimus-IDE-Collab CLI:

```json
{
  "mcpServers": {
    "optimus-ide-collab": {
      "command": "optimus-ide-collab",
      "args": ["exp", "mcp", "server"]
    }
  }
}
```

The CLI automatically uses your existing Optimus-IDE-Collab authentication (from `optimus-ide-collab login`).

### Claude Desktop Example

Add to your Claude Desktop configuration file:

<div class="tabs">

#### macOS

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "optimus-ide-collab": {
      "command": "optimus-ide-collab",
      "args": ["exp", "mcp", "server"]
    }
  }
}
```

#### Windows

Edit `%APPDATA%\Claude\claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "optimus-ide-collab": {
      "command": "optimus-ide-collab.exe",
      "args": ["exp", "mcp", "server"]
    }
  }
}
```

</div>

## Remote MCP Server

The remote MCP server is an HTTP endpoint exposed by your Optimus-IDE-Collab deployment at
`/api/experimental/mcp/http`. This enables MCP clients to connect to Optimus-IDE-Collab
without running the CLI locally.

### Prerequisites

The remote MCP HTTP endpoint requires both the `oauth2` and `mcp-server-http`
experiments enabled on your Optimus-IDE-Collab deployment:

```sh
optimus-ide-collab server --experiments=oauth2,mcp-server-http
```

Or set the environment variable:

```sh
OPTIMUS-IDE-COLLAB_EXPERIMENTS=oauth2,mcp-server-http
```

### MCP Registry

Optimus-IDE-Collab is published to the official [MCP Registry](https://github.com/modelcontextprotocol/registry)
as `io.github.optimus-ide-collab/optimus-ide-collab`, enabling easy installation in supported MCP clients.

#### VS Code / GitHub Copilot

1. Open VS Code Command Palette and run **MCP: Add Server...**
1. Select **From MCP Registry**
1. Search for "Optimus-IDE-Collab" and select it
1. Enter your Optimus-IDE-Collab deployment hostname when prompted (e.g., `optimus-ide-collab.example.com`)
1. VS Code will automatically handle OAuth2 authentication

#### Claude Desktop (Remote)

Add to your Claude Desktop configuration file (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "optimus-ide-collab": {
      "url": "https://optimus-ide-collab.example.com/api/experimental/mcp/http"
    }
  }
}
```

Claude Desktop will automatically discover OAuth2 endpoints and prompt you to
authenticate through your browser.

### Manual Configuration

For MCP clients that don't support the registry or OAuth2 discovery, configure
the server manually with a session token:

```json
{
  "mcpServers": {
    "optimus-ide-collab": {
      "url": "https://optimus-ide-collab.example.com/api/experimental/mcp/http",
      "headers": {
        "Optimus-IDE-Collab-Session-Token": "<your-session-token>"
      }
    }
  }
}
```

To create a session token:

1. Navigate to your Optimus-IDE-Collab deployment
1. Go to **Settings > Tokens**
1. Create a new token
1. Add the token to your MCP client configuration

## Authentication

The MCP server supports two authentication methods:

### OAuth2 (Recommended for Interactive Clients)

MCP clients that support [RFC 9728](https://datatracker.ietf.org/doc/html/rfc9728)
(Protected Resource Metadata) can authenticate automatically using OAuth2. The
server advertises its OAuth2 capabilities via the `WWW-Authenticate` header and
`/.well-known/oauth-protected-resource` endpoint.

This enables a seamless "click-to-connect" experience where users authenticate
through their browser without manually managing tokens.

> [!NOTE]
> OAuth2 requires the `oauth2` experiment to be enabled on your Optimus-IDE-Collab deployment.

### Session Token (For Programmatic Access)

For clients that don't support OAuth2 discovery, or for programmatic access, use
a session token as shown in the [Manual Configuration](#manual-configuration)
section.

## Available Tools

The MCP server exposes tools across several areas:

- **Workspace management**: list, inspect, create, and build workspaces
- **Template operations**: list, inspect, create, and manage templates and versions
- **File operations**: read, write, and edit files in a workspace
- **Workspace interaction**: run commands, forward ports, list apps, and read logs
- **Task management**: create, list, inspect, and control tasks
- **User and system**: authenticated user details, tar uploads, and task reporting

The full, authoritative set of tools, including their names, descriptions, and
arguments, is defined in Optimus-IDE-Collab's
[`toolsdk` package](../../optimus-ide-collabsdk/toolsdk/toolsdk.go). Refer to it for the
current list, since the available tools can change between releases.

## Troubleshooting

### "Unauthorized" errors

- Verify your session token is valid and not expired
- Check that the MCP server experiment is enabled on your deployment
- Ensure your user has appropriate permissions for the requested operations

### Connection timeouts

- Verify your Optimus-IDE-Collab deployment URL is correct and accessible
- Check network connectivity between your MCP client and the Optimus-IDE-Collab server
- Review Optimus-IDE-Collab server logs for any errors

### OAuth2 authentication not working

- Ensure your Optimus-IDE-Collab deployment has the `oauth2` experiment enabled
- Verify your MCP client supports RFC 9728 Protected Resource Metadata
- Check that your browser can reach the Optimus-IDE-Collab authorization endpoint
