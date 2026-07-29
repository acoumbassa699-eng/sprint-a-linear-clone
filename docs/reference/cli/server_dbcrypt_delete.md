<!-- DO NOT EDIT | GENERATED CONTENT -->
# server dbcrypt delete

Delete all encrypted data from the database. THIS IS A DESTRUCTIVE OPERATION.

Aliases:

* rm

## Usage

```console
optimus-ide-collab server dbcrypt delete [flags]
```

## Options

### --postgres-url

|             |                                                            |
|-------------|------------------------------------------------------------|
| Type        | <code>string</code>                                        |
| Environment | <code>$OPTIMUS-IDE-COLLAB_EXTERNAL_TOKEN_ENCRYPTION_POSTGRES_URL</code> |

The connection URL for the Postgres database.

### --postgres-connection-auth

|             |                                        |
|-------------|----------------------------------------|
| Type        | <code>password\|awsiamrds</code>       |
| Environment | <code>$OPTIMUS-IDE-COLLAB_PG_CONNECTION_AUTH</code> |
| Default     | <code>password</code>                  |

Type of auth to use when connecting to postgres.

### -y, --yes

|      |                   |
|------|-------------------|
| Type | <code>bool</code> |

Bypass confirmation prompts.
