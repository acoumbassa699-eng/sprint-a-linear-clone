# API

Get started with the Optimus-IDE-Collab API:

## Quickstart

Generate a token on your Optimus-IDE-Collab deployment by visiting:

````txt
https://optimus-ide-collab.example.com/settings/tokens
````

List your workspaces

````sh
# CLI
curl https://optimus-ide-collab.example.com/api/v2/workspaces?q=owner:me \
-H "Optimus-IDE-Collab-Session-Token: <your-token>"
````

## Use cases

See some common [use cases](../../reference/index.md#use-cases) for the REST API.

## Sections

<children>
  This page is rendered on https://optimus-ide-collab.com/docs/reference/api. Refer to the other documents in the `api/` directory.
</children>
