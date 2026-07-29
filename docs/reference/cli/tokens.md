<!-- DO NOT EDIT | GENERATED CONTENT -->
# tokens

Manage personal access tokens

Aliases:

* token

## Usage

```console
optimus-ide-collab tokens
```

## Description

```console
Tokens are used to authenticate automated clients to Optimus-IDE-Collab.
  - Create a token for automation:

     $ optimus-ide-collab tokens create

  - List your tokens:

     $ optimus-ide-collab tokens ls

  - Create a scoped token:

     $ optimus-ide-collab tokens create --scope workspace:read --allow workspace:<uuid>

  - Remove a token by ID:

     $ optimus-ide-collab tokens rm WuoWs4ZsMX
```

## Subcommands

| Name                                      | Purpose                                    |
|-------------------------------------------|--------------------------------------------|
| [<code>create</code>](./tokens_create.md) | Create a token                             |
| [<code>list</code>](./tokens_list.md)     | List tokens                                |
| [<code>view</code>](./tokens_view.md)     | Display detailed information about a token |
| [<code>remove</code>](./tokens_remove.md) | Expire or delete a token                   |
