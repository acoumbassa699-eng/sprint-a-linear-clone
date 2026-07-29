# AI Providers

## List AI providers

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/ai/providers \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`GET /api/v2/ai/providers`

### Example responses

> 200 Response

```json
[
  {
    "api_keys": [
      {
        "created_at": "2019-08-24T14:15:22Z",
        "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
        "masked": "string"
      }
    ],
    "base_url": "string",
    "created_at": "2019-08-24T14:15:22Z",
    "display_name": "string",
    "enabled": true,
    "icon": "string",
    "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
    "name": "string",
    "settings": {},
    "type": "openai",
    "updated_at": "2019-08-24T14:15:22Z"
  }
]
```

### Responses

| Status | Meaning                                                 | Description | Schema                                                        |
|--------|---------------------------------------------------------|-------------|---------------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | array of [optimus-ide-collabsdk.AIProvider](schemas.md#optimus-ide-collabsdkaiprovider) |

<h3 id="list-ai-providers-responseschema">Response Schema</h3>

Status Code **200**

| Name             | Type                                                                 | Required | Restrictions | Description |
|------------------|----------------------------------------------------------------------|----------|--------------|-------------|
| `[array item]`   | array                                                                | false    |              |             |
| `» api_keys`     | array                                                                | false    |              |             |
| `»» created_at`  | string(date-time)                                                    | false    |              |             |
| `»» id`          | string(uuid)                                                         | false    |              |             |
| `»» masked`      | string                                                               | false    |              |             |
| `» base_url`     | string                                                               | false    |              |             |
| `» created_at`   | string(date-time)                                                    | false    |              |             |
| `» display_name` | string                                                               | false    |              |             |
| `» enabled`      | boolean                                                              | false    |              |             |
| `» icon`         | string                                                               | false    |              |             |
| `» id`           | string(uuid)                                                         | false    |              |             |
| `» name`         | string                                                               | false    |              |             |
| `» settings`     | [optimus-ide-collabsdk.AIProviderSettings](schemas.md#optimus-ide-collabsdkaiprovidersettings) | false    |              |             |
| `» type`         | [optimus-ide-collabsdk.AIProviderType](schemas.md#optimus-ide-collabsdkaiprovidertype)         | false    |              |             |
| `» updated_at`   | string(date-time)                                                    | false    |              |             |

#### Enumerated Values

| Property | Value(s)                                                                                                |
|----------|---------------------------------------------------------------------------------------------------------|
| `type`   | `anthropic`, `azure`, `bedrock`, `copilot`, `google`, `openai`, `openai-compat`, `openrouter`, `vercel` |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Create an AI provider

### Code samples

```sh
# Example request using curl
curl -X POST http://optimus-ide-collab-server:8080/api/v2/ai/providers \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`POST /api/v2/ai/providers`

> Body parameter

```json
{
  "api_keys": [
    "string"
  ],
  "base_url": "string",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "name": "string",
  "settings": {},
  "type": "openai"
}
```

### Parameters

| Name   | In   | Type                                                                           | Required | Description                |
|--------|------|--------------------------------------------------------------------------------|----------|----------------------------|
| `body` | body | [optimus-ide-collabsdk.CreateAIProviderRequest](schemas.md#optimus-ide-collabsdkcreateaiproviderrequest) | true     | Create AI provider request |

### Example responses

> 201 Response

```json
{
  "api_keys": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "masked": "string"
    }
  ],
  "base_url": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "settings": {},
  "type": "openai",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Responses

| Status | Meaning                                                      | Description | Schema                                               |
|--------|--------------------------------------------------------------|-------------|------------------------------------------------------|
| 201    | [Created](https://tools.ietf.org/html/rfc7231#section-6.3.2) | Created     | [optimus-ide-collabsdk.AIProvider](schemas.md#optimus-ide-collabsdkaiprovider) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Get an AI provider

### Code samples

```sh
# Example request using curl
curl -X GET http://optimus-ide-collab-server:8080/api/v2/ai/providers/{idOrName} \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`GET /api/v2/ai/providers/{idOrName}`

### Parameters

| Name       | In   | Type   | Required | Description         |
|------------|------|--------|----------|---------------------|
| `idOrName` | path | string | true     | Provider ID or name |

### Example responses

> 200 Response

```json
{
  "api_keys": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "masked": "string"
    }
  ],
  "base_url": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "settings": {},
  "type": "openai",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                               |
|--------|---------------------------------------------------------|-------------|------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.AIProvider](schemas.md#optimus-ide-collabsdkaiprovider) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Delete an AI provider

### Code samples

```sh
# Example request using curl
curl -X DELETE http://optimus-ide-collab-server:8080/api/v2/ai/providers/{idOrName} \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`DELETE /api/v2/ai/providers/{idOrName}`

### Parameters

| Name       | In   | Type   | Required | Description         |
|------------|------|--------|----------|---------------------|
| `idOrName` | path | string | true     | Provider ID or name |

### Responses

| Status | Meaning                                                         | Description | Schema |
|--------|-----------------------------------------------------------------|-------------|--------|
| 204    | [No Content](https://tools.ietf.org/html/rfc7231#section-6.3.5) | No Content  |        |

To perform this operation, you must be authenticated. [Learn more](authentication.md).

## Update an AI provider

### Code samples

```sh
# Example request using curl
curl -X PATCH http://optimus-ide-collab-server:8080/api/v2/ai/providers/{idOrName} \
  -H 'Content-Type: application/json' \
  -H 'Accept: application/json' \
  -H 'Optimus-IDE-Collab-Session-Token: API_KEY'
```

`PATCH /api/v2/ai/providers/{idOrName}`

> Body parameter

```json
{
  "api_keys": [
    {
      "api_key": "string",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08"
    }
  ],
  "base_url": "string",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "settings": {}
}
```

### Parameters

| Name       | In   | Type                                                                           | Required | Description                |
|------------|------|--------------------------------------------------------------------------------|----------|----------------------------|
| `idOrName` | path | string                                                                         | true     | Provider ID or name        |
| `body`     | body | [optimus-ide-collabsdk.UpdateAIProviderRequest](schemas.md#optimus-ide-collabsdkupdateaiproviderrequest) | true     | Update AI provider request |

### Example responses

> 200 Response

```json
{
  "api_keys": [
    {
      "created_at": "2019-08-24T14:15:22Z",
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "masked": "string"
    }
  ],
  "base_url": "string",
  "created_at": "2019-08-24T14:15:22Z",
  "display_name": "string",
  "enabled": true,
  "icon": "string",
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "name": "string",
  "settings": {},
  "type": "openai",
  "updated_at": "2019-08-24T14:15:22Z"
}
```

### Responses

| Status | Meaning                                                 | Description | Schema                                               |
|--------|---------------------------------------------------------|-------------|------------------------------------------------------|
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | OK          | [optimus-ide-collabsdk.AIProvider](schemas.md#optimus-ide-collabsdkaiprovider) |

To perform this operation, you must be authenticated. [Learn more](authentication.md).
