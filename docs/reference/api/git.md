# Git

## Get user external auths

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/external-auth \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`GET /api/v2/external-auth`

### Example responses

> 200 Response

```json
{
  "authenticated": true,
  "created_at": "2019-08-24T14:15:22Z",
  "expires": "2019-08-24T14:15:22Z",
  "has_refresh_token": true,
  "provider_id": "string",
  "updated_at": "2019-08-24T14:15:22Z",
  "validate_error": "string"
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                           |
|--------|---------------------------------------------------------|-------------|------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.ExternalAuthLink](schemas.md#optimus-ide-collabsdkexternalauthlink) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Get external auth by ID

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/external-auth/{externalauth} \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`GET /api/v2/external-auth/{externalauth}`

### Parameters

| Name           | In   | Type           | Required | Description     |
|----------------|------|----------------|----------|-----------------|
| `externalauth` | path | string(string) | true     | Git Provider ID |

### Example responses

> 200 Response

```json
{
  "app_install_url": "string",
  "app_installable": true,
  "authenticated": true,
  "device": true,
  "display_name": "string",
  "installations": [
    {
      "account": {
        "avatar_url": "string",
        "id": 0,
        "login": "string",
        "name": "string",
        "profile_url": "string"
      },
      "configure_url": "string",
      "id": 0
    }
  ],
  "supports_revocation": true,
  "user": {
    "avatar_url": "string",
    "id": 0,
    "login": "string",
    "name": "string",
    "profile_url": "string"
  }
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                   |
|--------|---------------------------------------------------------|-------------|----------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.ExternalAuth](schemas.md#optimus-ide-collabsdkexternalauth) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Delete external auth user link by ID

### Code samples

```sh
# Example request using curl
curl -X DELETE http://optimus-ide-collab-server:8080/api/v2/external-auth/{externalauth} \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`DELETE /api/v2/external-auth/{externalauth}`

### Parameters

| Name           | In   | Type           | Required | Description     |
|----------------|------|----------------|----------|-----------------|
| `externalauth` | path | string(string) | true     | Git Provider ID |

### Example responses

> 200 Response

```json
{
  "token_revocation_error": "string",
  "token_revoked": true
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                                                       |
|--------|---------------------------------------------------------|-------------|----------------------------------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.DeleteExternalAuthByIDResponse](schemas.md#optimus-ide-collabsdkdeleteexternalauthbyidresponse) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Get external auth device by ID

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/external-auth/{externalauth}/device \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`GET /api/v2/external-auth/{externalauth}/device`

### Parameters

| Name           | In   | Type           | Required | Description     |
|----------------|------|----------------|----------|-----------------|
| `externalauth` | path | string(string) | true     | Git Provider ID |

### Example responses

> 200 Response

```json
{
  "device_code": "string",
  "expires_in": 0,
  "interval": 0,
  "user_code": "string",
  "verification_uri": "string"
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                               |
|--------|---------------------------------------------------------|-------------|----------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.ExternalAuthDevice](schemas.md#optimus-ide-collabsdkexternalauthdevice) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Post external auth device by ID

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/external-auth/{externalauth}/device \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`POST /api/v2/external-auth/{externalauth}/device`

### Parameters

| Name           | In   | Type           | Required | Description          |
|----------------|------|----------------|----------|----------------------|
| `externalauth` | path | string(string) | true     | External Provider ID |

### Responses

| Status | Meaning                                                         | Description | Schema |
|--------|-----------------------------------------------------------------|-------------|--------|
| 204    | [No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5) | No Content  |        |

To perform this operation, you must be authenticated. [Learn more](authentication.md).
