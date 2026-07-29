<!-- DO NOT EDIT | GENERATED CONTENT -->
# oauth2-provider

Manage Optimus-IDE-Collab OAuth2 provider settings

## Usage

```console
optimus-ide-collab oauth2-provider
```

## Description

```console
Administrators can use these commands to change OAuth2 provider settings.
  - Enable dynamic client registration (RFC 7591), allowing OAuth2/MCP clients to
self-register without an admin creating an app first:

     $ optimus-ide-collab oauth2-provider dcr enable

  - Disable dynamic client registration. Clients that already registered are
unaffected; only new self-registration attempts are rejected:

     $ optimus-ide-collab oauth2-provider dcr disable
```

## Subcommands

| Name                                         | Purpose                                              |
|----------------------------------------------|------------------------------------------------------|
| [<code>dcr</code>](./oauth2-provider_dcr.md) | Manage OAuth2 dynamic client registration (RFC 7591) |
