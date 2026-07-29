# Authorization

## List API key scopes

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/auth/scopes \
  -H 'Accept: application/json'
```

`GET /api/v2/auth/scopes`

### Example responses

> 200 Response

```json
{
  "external": [
    "all"
  ]
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                                   |
|--------|---------------------------------------------------------|-------------|--------------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.ExternalAPIKeyScopes](schemas.md#optimus-ide-collabsdkexternalapikeyscopes) |

## Check authorization

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/authcheck \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`POST /api/v2/authcheck`

> Body parameter

```json
{
  "checks": {
    "property1": {
      "action": "create",
      "object": {
        "any_org": true,
        "organization_id": "string",
        "owner_id": "string",
        "resource_id": "string",
        "resource_type": "*"
      }
    },
    "property2": {
      "action": "create",
      "object": {
        "any_org": true,
        "organization_id": "string",
        "owner_id": "string",
        "resource_id": "string",
        "resource_type": "*"
      }
    }
  }
}
```

### Parameters

| Name   | In   | Type                                                                     | Required | Description           |
|--------|------|--------------------------------------------------------------------------|----------|-----------------------|
| `body` | body | [optimus-ide-collabsdk.AuthorizationRequest](schemas.md#optimus-ide-collabsdkauthorizationrequest) | true     | Authorization request |

### Example responses

> 200 Response

```json
{
  "property1": true,
  "property2": true
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                                     |
|--------|---------------------------------------------------------|-------------|----------------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.AuthorizationResponse](schemas.md#optimus-ide-collabsdkauthorizationresponse) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Log in user

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/users/login \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json'
```

`POST /api/v2/users/login`

> Body parameter

```json
{
  "email": "user@example.com",
  "password": "string"
}
```

### Parameters

| Name   | In   | Type                                                                             | Required | Description   |
|--------|------|----------------------------------------------------------------------------------|----------|---------------|
| `body` | body | [optimus-ide-collabsdk.LoginWithPasswordRequest](schemas.md#optimus-ide-collabsdkloginwithpasswordrequest) | true     | Login request |

### Example responses

> 201 Response

```json
{
  "session_token": "string"
}
```

### Responses

| Status | Meaning                                                      | Description | Schema                                                                             |
|--------|--------------------------------------------------------------|-------------|------------------------------------------------------------------------------------|
| 201    | [Created](https://tools.ietf.org/html/rfc7231#section-6.3.2) | Created     | [optimus-ide-collabsdk.LoginWithPasswordResponse](schemas.md#optimus-ide-collabsdkloginwithpasswordresponse) |

## Change password with a one-time passcode

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/users/otp/change-password \
  -H 'Content-Type: application/json'
```

`POST /api/v2/users/otp/change-password`

> Body parameter

```json
{
  "email": "user@example.com",
  "one_time_passcode": "string",
  "password": "string"
}
```

### Parameters

| Name   | In   | Type                                                                                                             | Required | Description             |
|--------|------|------------------------------------------------------------------------------------------------------------------|----------|-------------------------|
| `body` | body | [optimus-ide-collabsdk.ChangePasswordWithOneTimePasscodeRequest](schemas.md#optimus-ide-collabsdkchangepasswordwithonetimepassoptimus-ide-collabequest) | true     | Change password request |

### Responses

| Status | Meaning                                                         | Description | Schema |
|--------|-----------------------------------------------------------------|-------------|--------|
| 204    | [No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5) | No Content  |        |

## Request one-time passcode

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/users/otp/request \
  -H 'Content-Type: application/json'
```

`POST /api/v2/users/otp/request`

> Body parameter

```json
{
  "email": "user@example.com"
}
```

### Parameters

| Name   | In   | Type                                                                                       | Required | Description               |
|--------|------|--------------------------------------------------------------------------------------------|----------|---------------------------|
| `body` | body | [optimus-ide-collabsdk.RequestOneTimePasscodeRequest](schemas.md#optimus-ide-collabsdkrequestonetimepassoptimus-ide-collabequest) | true     | One-time passcode request |

### Responses

| Status | Meaning                                                         | Description | Schema |
|--------|-----------------------------------------------------------------|-------------|--------|
| 204    | [No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5) | No Content  |        |

## Validate user password

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/users/validate-password \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`POST /api/v2/users/validate-password`

> Body parameter

```json
{
  "password": "string"
}
```

### Parameters

| Name   | In   | Type                                                                                   | Required | Description                    |
|--------|------|----------------------------------------------------------------------------------------|----------|--------------------------------|
| `body` | body | [optimus-ide-collabsdk.ValidateUserPasswordRequest](schemas.md#optimus-ide-collabsdkvalidateuserpasswordrequest) | true     | Validate user password request |

### Example responses

> 200 Response

```json
{
  "details": "string",
  "valid": true
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                                                   |
|--------|---------------------------------------------------------|-------------|------------------------------------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.ValidateUserPasswordResponse](schemas.md#optimus-ide-collabsdkvalidateuserpasswordresponse) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Convert user from password to oauth authentication

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/users/{user}/convert-login \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`POST /api/v2/users/{user}/convert-login`

> Body parameter

```json
{
  "password": "string",
  "to_type": ""
}
```

### Parameters

| Name   | In   | Type                                                                   | Required | Description          |
|--------|------|------------------------------------------------------------------------|----------|----------------------|
| `user` | path | string                                                                 | true     | User ID, name, or me |
| `body` | body | [optimus-ide-collabsdk.ConvertLoginRequest](schemas.md#optimus-ide-collabsdkconvertloginrequest) | true     | Convert request      |

### Example responses

> 201 Response

```json
{
  "expires_at": "2019-08-24T14:15:22Z",
  "state_string": "string",
  "to_type": "",
  "user_id": "a169451c-8525-4352-b8ca-070dd449a1a5"
}
```

### Responses

| Status | Meaning                                                      | Description | Schema                                                                         |
|--------|--------------------------------------------------------------|-------------|--------------------------------------------------------------------------------|
| 201    | [Created](https://tools.ietf.org/html/rfc7231#section-6.3.2) | Created     | [optimus-ide-collabsdk.OAuthConversionResponse](schemas.md#optimus-ide-collabsdkoauthconversionresponse) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).
