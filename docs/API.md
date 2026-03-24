# API Reference - david CalDAV Server

## Base URLs

- Development: `http://localhost:8080`
- Production: `https://your-domain.com`

## Authentication

All protected endpoints require authentication via one of these methods:

### Bearer Token (JWT)

```bash
Authorization: Bearer <access_token>
```

### API Key

```bash
Authorization: API-Key <api_key>
```

### Basic Auth

```bash
Authorization: Basic <base64(username:password)>
```

## Endpoints

### Health Check

**GET** `/api/health`

Return server health status.

**Response**: `200 OK`
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "uptime": "24h30m15s"
}
```

### Auth Endpoints

#### Login

**POST** `/api/auth/login`

Authenticate user and return tokens.

**Request**:
```json
{
  "username": "user@example.com",
  "password": "secret123"
}
```

**Response**: `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### Refresh Token

**POST** `/api/auth/refresh`

Get new access token using refresh token.

**Request**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response**: `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

#### Logout

**POST** `/api/auth/logout`

Logout and invalidate token.

**Headers**: `Authorization: Bearer <token>`

**Response**: `200 OK`

### Calendar Endpoints

#### List Calendars

**GET** `/api/calendars`

List all accessible calendars.

**Query Parameters**:
- `page` (int) - Page number (default: 1)
- `limit` (int) - Items per page (default: 50, max: 100)

**Response**: `200 OK`
```json
{
  "calendars": [
    {
      "uid": "cal-123",
      "display_name": "Work Calendar",
      "description": "Work events",
      "color": "#FF0000",
      "is_public": false,
      "owner_id": "user-456",
      "created": "2026-03-01T00:00:00Z",
      "updated": "2026-03-24T00:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 50
}
```

#### Get Calendar

**GET** `/api/calendars/{uid}`

Get calendar details.

**Response**: `200 OK`
```json
{
  "uid": "cal-123",
  "display_name": "Work Calendar",
  "description": "Work events",
  "color": "#FF0000",
  "is_public": false,
  "timezone": "America/New_York",
  "owner_id": "user-456"
}
```

#### Create Calendar

**POST** `/api/calendars`

Create new calendar.

**Request**:
```json
{
  "display_name": "New Calendar",
  "description": "My new calendar",
  "color": "#00FF00",
  "timezone": "UTC"
}
```

**Response**: `201 Created`
```json
{
  "uid": "cal-789",
  "display_name": "New Calendar",
  "created": "2026-03-24T00:00:00Z"
}
```

#### Update Calendar

**PUT** `/api/calendars/{uid}`

Update calendar.

**Request**:
```json
{
  "display_name": "Updated Name",
  "color": "#0000FF"
}
```

**Response**: `200 OK`

#### Delete Calendar

**DELETE** `/api/calendars/{uid}`

Delete calendar.

**Response**: `204 No Content`

### Event Endpoints

#### List Events

**GET** `/api/calendars/{calendar_uid}/events`

List events in calendar.

**Query Parameters**:
- `start` (string) - Start date (ISO 8601)
- `end` (string) - End date (ISO 8601)
- `page` (int) - Page number
- `limit` (int) - Items per page

**Response**: `200 OK`
```json
{
  "events": [
    {
      "uid": "event-123",
      "calendar_uid": "cal-456",
      "summary": "Team Meeting",
      "description": "Weekly sync",
      "location": "Conference Room A",
      "dtstart": "2026-03-25T10:00:00Z",
      "dtend": "2026-03-25T11:00:00Z",
      "all_day": false,
      "status": "CONFIRMED"
    }
  ],
  "total": 25,
  "page": 1
}
```

#### Get Event

**GET** `/api/events/{uid}`

Get event details.

**Response**: `200 OK`
```json
{
  "uid": "event-123",
  "calendar_uid": "cal-456",
  "summary": "Team Meeting",
  "description": "Weekly sync",
  "location": "Conference Room A",
  "dtstart": "2026-03-25T10:00:00Z",
  "dtend": "2026-03-25T11:00:00Z",
  "organizer": "organizer@example.com",
  "attendees": [
    {
      "email": "user@example.com",
      "status": "ACCEPTED"
    }
  ]
}
```

#### Create Event

**POST** `/api/calendars/{calendar_uid}/events`

Create new event.

**Request**:
```json
{
  "summary": "New Event",
  "description": "Event description",
  "location": "Location",
  "dtstart": "2026-03-25T10:00:00Z",
  "dtend": "2026-03-25T11:00:00Z",
  "all_day": false
}
```

**Response**: `201 Created`

#### Update Event

**PUT** `/api/events/{uid}`

Update event.

**Response**: `200 OK`

#### Delete Event

**DELETE** `/api/events/{uid}`

Delete event.

**Response**: `204 No Content`

### User Endpoints

#### List Users

**GET** `/api/users`

List all users (admin only).

**Response**: `200 OK`
```json
{
  "users": [
    {
      "id": "user-123",
      "username": "john@example.com",
      "email": "john@example.com",
      "display_name": "John Doe",
      "role": "user",
      "created": "2026-03-01T00:00:00Z"
    }
  ]
}
```

#### Get User

**GET** `/api/users/{id}`

Get user details.

**Response**: `200 OK`

#### Create User

**POST** `/api/users`

Create new user (admin only).

**Request**:
```json
{
  "username": "new@example.com",
  "password": "securepassword",
  "email": "new@example.com",
  "display_name": "New User",
  "role": "user"
}
```

**Response**: `201 Created`

#### Update User

**PUT** `/api/users/{id}`

Update user.

**Request**:
```json
{
  "display_name": "Updated Name",
  "role": "manager"
}
```

**Response**: `200 OK`

#### Delete User

**DELETE** `/api/users/{id}`

Delete user (admin only).

**Response**: `204 No Content`

#### Create API Key

**POST** `/api/users/{id}/api-keys`

Create API key for user.

**Response**: `201 Created`
```json
{
  "key_id": "key-123",
  "api_key": "david_abc123def456...",
  "name": "Desktop App",
  "created": "2026-03-24T00:00:00Z"
}
```

#### List API Keys

**GET** `/api/users/{id}/api-keys`

List API keys for user.

**Response**: `200 OK`

#### Revoke API Key

**DELETE** `/api/api-keys/{key_id}`

Revoke API key.

**Response**: `204 No Content`

### Admin Endpoints

#### System Statistics

**GET** `/api/admin/statistics`

Get system-wide statistics.

**Response**: `200 OK`
```json
{
  "system": {
    "total_users": 50,
    "total_calendars": 120,
    "total_events": 5000,
    "total_todos": 200,
    "total_journals": 50,
    "total_storage_bytes": 1073741824,
    "startup_time": "2026-03-24T00:00:00Z"
  },
  "users": [...],
  "calendars": [...]
}
```

#### Audit Logs

**GET** `/api/admin/audit`

Get audit log entries.

**Query Parameters**:
- `page` (int)
- `limit` (int)
- `user_id` (string) - Filter by user
- `action` (string) - Filter by action

**Response**: `200 OK`

#### Rate Limits

**GET** `/api/admin/rate-limits`

List all rate limits.

**POST** `/api/admin/rate-limits`

Create rate limit.

**Request**:
```json
{
  "identifier": "192.168.1.100",
  "requests_per_minute": 30,
  "burst": 5
}
```

**DELETE** `/api/admin/rate-limits/{id}`

Remove rate limit.

**Response**: `204 No Content`

### WebDAV/CalDAV Endpoints

All WebDAV endpoints are available at `/webdav/`.

#### Calendar Collection

**PROPFIND** `/webdav/{user_id}/calendars/{calendar_uid}/`

Get calendar properties.

**MKCALENDAR** `/webdav/{user_id}/calendars/{calendar_uid}/`

Create calendar.

**REPORT** `/webdav/{user_id}/calendars/{calendar_uid}/`

Query calendar events.

```xml
<C:calendar-query>
  <D:prop>
    <C:calendar-data/>
  </D:prop>
  <C:comp-filter name="VEVENT">
    <C:time-range start="20260301T000000Z" end="20260331T235959Z"/>
  </C:comp-filter>
</C:calendar-query>
```

#### Event File

**GET** `/webdav/{user_id}/calendars/{calendar_uid}/{event_uid}.ics`

Download event as ICS.

**PUT** `/webdav/{user_id}/calendars/{calendar_uid}/{event_uid}.ics`

Upload/update event.

**DELETE** `/webdav/{user_id}/calendars/{calendar_uid}/{event_uid}.ics`

Delete event.

### Error Responses

All errors follow this format:

**400 Bad Request**
```json
{
  "error": "validation_error",
  "message": "Invalid request format"
}
```

**401 Unauthorized**
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token"
}
```

**403 Forbidden**
```json
{
  "error": "forbidden",
  "message": "Insufficient permissions"
}
```

**404 Not Found**
```json
{
  "error": "not_found",
  "message": "Resource not found"
}
```

**409 Conflict**
```json
{
  "error": "conflict",
  "message": "Resource already exists"
}
```

**429 Too Many Requests**
```json
{
  "error": "rate_limit_exceeded",
  "message": "Too many requests, please try again later",
  "retry_after": 60
}
```

**500 Internal Server Error**
```json
{
  "error": "internal_error",
  "message": "An unexpected error occurred"
}
```

### Rate Limiting

API requests are rate limited by default:

- 60 requests per minute per user
- 10 burst requests
- Adjusted via configuration

Rate limit headers included in responses:

```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1711234567
```

### Versioning

API version is included in URL path:

- Current: `/api/v1/...`
- Future: `/api/v2/...`

### Content Types

- Request: `application/json` (except file uploads)
- Response: `application/json`

File uploads/downloads use: `text/calendar; charset=utf-8`
