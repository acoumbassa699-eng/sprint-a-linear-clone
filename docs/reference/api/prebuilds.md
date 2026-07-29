# Prebuilds

## Get prebuilds settings

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/prebuilds/settings \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`GET /api/v2/prebuilds/settings`

### Example responses

> 200 Response

```json
{
  "reconciliation_paused": true
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                             |
|--------|---------------------------------------------------------|-------------|--------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.PrebuildsSettings](schemas.md#optimus-ide-collabsdkprebuildssettings) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Update prebuilds settings

### Code samples

```sh
# Example request using curl
curl -X PUT http://optimus-ide-collab-server:8080/api/v2/prebuilds/settings \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`PUT /api/v2/prebuilds/settings`

> Body parameter

```json
{
  "reconciliation_paused": true
}
```

### Parameters

| Name   | In   | Type                                                               | Required | Description                |
|--------|------|--------------------------------------------------------------------|----------|----------------------------|
| `body` | body | [optimus-ide-collabsdk.PrebuildsSettings](schemas.md#optimus-ide-collabsdkprebuildssettings) | true     | Prebuilds settings request |

### Example responses

> 200 Response

```json
{
  "reconciliation_paused": true
}
```

### Responses

| Status | Meaning                                                         | Description  | Schema                                                             |
|--------|-----------------------------------------------------------------|--------------|--------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)         | OK           | [optimus-ide-collabsdk.PrebuildsSettings](schemas.md#optimus-ide-collabsdkprebuildssettings) |
| 304    | [Not Modified](https://tools.ietf.org/html/rfc7232#section-4.1) | Not Modified |                                                                    |

To perform this operation, you must be authenticated. [Learn more](authentication.md).
