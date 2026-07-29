<!-- DO NOT EDIT | GENERATED CONTENT -->
# reset-password

Directly connect to the database to reset a user's password

## Usage

```console
optimus-ide-collab reset-password [flags] <username>
```

## Options

### --postgres-url

|             |                                       |
|-------------|---------------------------------------|
| Type        | <code>string</code>                   |
| Environment | <code>$OPTIMUS-IDE-COLLAB_PG_CONNECTION_URL</code> |

URL of a PostgreSQL database to connect to.

### --postgres-connection-auth

|             |                                        |
|-------------|----------------------------------------|
| Type        | <code>password\|awsiamrds</code>       |
| Environment | <code>$OPTIMUS-IDE-COLLAB_PG_CONNECTION_AUTH</code> |
| Default     | <code>password</code>                  |

Type of auth to use when connecting to postgres.
