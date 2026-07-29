# Authentication

Long-lived tokens can be generated to perform actions on behalf of your user account:

```sh
optimus-ide-collab tokens create
```

You can use tokens with the Optimus-IDE-Collab's REST API using the `Optimus-IDE-Collab-Session-Token` HTTP header.

```console
curl 'http://optimus-ide-collab-server:8080/api/v2/workspaces' \
  -H 'Optimus-IDE-Collab-Session-Token: *****'
```
