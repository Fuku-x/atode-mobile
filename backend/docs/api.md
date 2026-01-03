# Backend API Spec

- OpenAPI定義: [`openapi.yaml`](./openapi.yaml)

## Base

- Base URL: (environment-dependent)
- Content-Type: `application/json`

## Examples (curl)

Set these env vars for convenience:

```bash
export API_BASE_URL="http://localhost:8080"
export FIREBASE_ID_TOKEN="<FIREBASE_ID_TOKEN>"
export REQUEST_ID="dev-req-001"
```

### GET /me

```bash
curl -sS "$API_BASE_URL/me" \
  -H "Authorization: Bearer $FIREBASE_ID_TOKEN" \
  -H "X-Request-Id: $REQUEST_ID"
```

### GET /tasks

```bash
curl -sS "$API_BASE_URL/tasks" \
  -H "Authorization: Bearer $FIREBASE_ID_TOKEN" \
  -H "X-Request-Id: $REQUEST_ID"
```

### POST /tasks

```bash
curl -sS -X POST "$API_BASE_URL/tasks" \
  -H "Authorization: Bearer $FIREBASE_ID_TOKEN" \
  -H "X-Request-Id: $REQUEST_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk"}'
```

### PUT /tasks/{id}

```bash
export TASK_ID="<TASK_UUID>"

curl -sS -X PUT "$API_BASE_URL/tasks/$TASK_ID" \
  -H "Authorization: Bearer $FIREBASE_ID_TOKEN" \
  -H "X-Request-Id: $REQUEST_ID" \
  -H "Content-Type: application/json" \
  -d '{"isDone":true, "dueAt":null}'
```

### DELETE /tasks/{id}

```bash
curl -sS -X DELETE "$API_BASE_URL/tasks/$TASK_ID" \
  -H "Authorization: Bearer $FIREBASE_ID_TOKEN" \
  -H "X-Request-Id: $REQUEST_ID" \
  -i
```

## Authentication

- Required endpoints use **Firebase ID Token**.
- Send token via HTTP header:

```http
Authorization: Bearer <FIREBASE_ID_TOKEN>
```

If missing/invalid:

- Status: `401`
- Body: unified error response (see below)

## Request ID

- You can optionally provide a request id:

```http
X-Request-Id: <string>
```

- Server always responds with:

```http
X-Request-Id: <string>
```

- Error responses include `request_id` (same value as `X-Request-Id`).

## Unified Error Response

All error responses are JSON:

```json
{
  "code": "string",
  "message": "string",
  "request_id": "string"
}
```

Notes:

- `request_id` may be omitted when not available, but the middleware sets it for all requests.

### Common error codes

- `unauthorized`
- `invalid_json`
- `method_not_allowed`
- `internal_server_error`
- `not_found`

(Endpoint-specific codes are documented per endpoint.)

## Error code matrix (status x code)

This section is intended for frontend implementers to quickly map failure handling.

### General

| Status | code | When |
| --- | --- | --- |
| 401 | `unauthorized` | Missing/invalid Firebase ID token |
| 405 | `method_not_allowed` | Unsupported HTTP method |
| 500 | `internal_server_error` | Unexpected server failure |
| 404 | `not_found` | Resource not found (or path not matched) |

### Tasks

| Status | code | When |
| --- | --- | --- |
| 400 | `invalid_json` | Invalid JSON OR unknown field in PUT body |
| 400 | `title_required` | `title` is missing/empty after trim |
| 400 | `no_fields` | PUT body has no updatable fields |
| 400 | `invalid_task_id` | `/tasks/{id}` where `{id}` is not a UUID |
| 400 | `invalid_due_at` | `dueAt` is not RFC3339/RFC3339Nano string |
| 400 | `invalid_scheduled_at` | `scheduledAt` is not RFC3339/RFC3339Nano string |

## Data Types

### Timestamp

All timestamps are RFC3339/RFC3339Nano strings produced by Go `time.Time` JSON encoding.

Recommendations:

- Prefer UTC with `Z` suffix (e.g. `2025-01-01T00:00:00Z`).
- Frontend should send timestamp strings in the same format.

### Task

```json
{
  "id": "uuid",
  "userId": "uuid",
  "title": "string",
  "isDone": true,
  "dueAt": "2025-01-01T00:00:00Z",
  "scheduledAt": "2025-01-01T00:00:00Z",
  "createdAt": "2025-01-01T00:00:00Z",
  "updatedAt": "2025-01-01T00:00:00Z"
}
```

Notes:

- `dueAt` and `scheduledAt` may be omitted when `null` in DB (`omitempty`).

---

# Endpoints

## GET /healthz

Health check endpoint.

### Response

- `200` (text/plain): `ok`
- `503` (text/plain): `unhealthy`

---

## GET /me

Returns authenticated user info from the request context.

### Auth

Required.

### Response 200

```json
{
  "user_id": "uuid",
  "firebase_uid": "string",
  "email": "string"
}
```

### Errors

- `401 unauthorized`
- `500 internal_server_error`

---

## GET /tasks

List tasks for the authenticated user.

### Auth

Required.

### Response 200

```json
{
  "tasks": [
    {
      "id": "uuid",
      "userId": "uuid",
      "title": "string",
      "isDone": false,
      "createdAt": "2025-01-01T00:00:00Z",
      "updatedAt": "2025-01-01T00:00:00Z"
    }
  ]
}
```

### Sort order

Tasks are returned in the following order:

- `is_done` ASC
- `due_at` ASC (NULLs are treated as infinity / last)
- `created_at` DESC

### Errors

- `401 unauthorized`
- `500 internal_server_error`

---

## POST /tasks

Create a new task.

### Auth

Required.

### Request body

```json
{
  "title": "string"
}
```

Rules:

- `title` is required and trimmed.
- Empty title is rejected.

### Response 201

`Task`

### Errors

- `400 invalid_json`
- `400 title_required` (`message`: `title is required`)
- `401 unauthorized`
- `500 internal_server_error`

---

## PUT /tasks/{id}

Partial update of a task owned by the authenticated user.

### Auth

Required.

### Path params

- `id`: task UUID

If `id` is not a valid UUID:

- `400 invalid_task_id`

If path does not match `/tasks/{id}` exactly:

- `404 not_found`

### Request body

This endpoint accepts a partial JSON object. Unknown fields are rejected.

Allowed fields:

- `title`: string
- `isDone`: boolean
- `dueAt`: RFC3339/RFC3339Nano string, or `null` to clear
- `scheduledAt`: RFC3339/RFC3339Nano string, or `null` to clear

Examples:

```json
{ "isDone": true }
```

```json
{ "dueAt": null }
```

### Semantics of omitted vs null

- If a field is **omitted**, it is **not changed**.
- If `dueAt`/`scheduledAt` is **explicitly `null`**, the value is **cleared**.
- If `title`/`isDone` is **explicitly `null`**, it is treated as invalid JSON for that field type.

Rules:

- At least one field must be present.
- If `title` is present, it is trimmed and must be non-empty.
- If `dueAt`/`scheduledAt` is present and not `null`, it must be a valid timestamp string.

### Response 200

`Task`

### Errors

- `400 invalid_json` (invalid JSON or contains unknown fields)
- `400 title_required`
- `400 no_fields` (`message`: `no fields to update`)
- `400 invalid_due_at`
- `400 invalid_scheduled_at`
- `400 invalid_task_id`
- `401 unauthorized`
- `404 not_found`
- `500 internal_server_error`

---

## DELETE /tasks/{id}

Delete a task owned by the authenticated user.

### Auth

Required.

### Path params

- `id`: task UUID

If `id` is not a valid UUID:

- `400 invalid_task_id`

If path does not match `/tasks/{id}` exactly:

- `404 not_found`

### Response 204

No content.

### Errors

- `400 invalid_task_id`
- `401 unauthorized`
- `404 not_found`
- `500 internal_server_error`
