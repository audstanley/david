# Configuration Reference

## Overview

david uses Viper for configuration management, supporting YAML config files and environment variables.

## Configuration File

The primary configuration file is `config.yaml`. david searches for this file in several locations:

1. `./config/config.yaml`
2. `$HOME/.swd/config.yaml`
3. `$HOME/.david/config.yaml`
4. `./config.yaml`

Alternatively, specify with `--config` flag:
```bash
david server --config /path/to/config.yaml
```

## Environment Variables

All configuration options can be overridden via environment variables with the `DAVID_` prefix and dots replaced with underscores.

| Config Path | Environment Variable |
|-------------|---------------------|
| `database.base_dir` | `DAVID_DATABASE_BASE_DIR` |
| `webdav.address` | `DAVID_WEBDAV_ADDRESS` |
| `webdav.port` | `DAVID_WEBDAV_PORT` |
| `api.address` | `DAVID_API_ADDRESS` |
| `api.port` | `DAVID_API_PORT` |
| `auth.jwt.secret` | `DAVID_JWT_SECRET` |
| `log.level` | `DAVID_LOG_LEVEL` |

**Important:** `DAVID_JWT_SECRET` must be set via environment variable for security.

## Configuration Sections

### Database

```yaml
database:
  base_dir: "/var/lib/david"
  calendars_db: "calendars.db"
  events_db: "events.db"
  todos_db: "todos.db"
  journals_db: "journals.db"
  freebusy_db: "freebusy.db"
  timezones_db: "timezones.db"
  users_db: "users.db"
  jwt_blacklist_db: "jwt_blacklist.db"
  audit_db: "audit.db"
```

- `base_dir` - Base directory for all database files
- All `*_db` fields - Database filenames (stored in `base_dir/data/`)

### WebDAV

```yaml
webdav:
  address: "127.0.0.1"
  port: "8000"
  prefix: "/webdav"
  tls:
    enabled: false
    cert_file: "/etc/david/cert.pem"
    key_file: "/etc/david/key.pem"
```

- `address` - Bind address
- `port` - Listening port
- `prefix` - URL prefix for WebDAV paths
- `tls.*` - TLS configuration (if enabled)

### API

```yaml
api:
  address: "127.0.0.1"
  port: "8080"
  cors_enabled: true
  cors_origins:
    - "http://localhost:3000"
  auth_methods:
    - "jwt"
    - "api_key"
    - "basic"
  rate_limiting:
    enabled: true
    default_limit: 100
    window: "1m"
    per_ip: true
```

- `address` - API bind address
- `port` - API listening port
- `cors_enabled` - Enable CORS headers
- `cors_origins` - Allowed origins
- `auth_methods` - Enabled authentication methods
- `rate_limiting.*` - Rate limiting configuration

### Authentication

```yaml
auth:
  jwt:
    secret: ""  # Set via DAVID_JWT_SECRET env var
    expiry: "24h"
    refresh_expiry: "7d"
  api_key:
    prefix: "david_"
    expiry: "365d"
```

- `jwt.secret` - **REQUIRED** - JWT signing secret (use env var)
- `jwt.expiry` - Access token lifetime
- `jwt.refresh_expiry` - Refresh token lifetime
- `api_key.prefix` - API key prefix
- `api_key.expiry` - API key lifetime

### Recurrence

```yaml
recurrence:
  cache_expiry_days: 30
  max_cache_instances: 1000
```

- `cache_expiry_days` - How long to cache recurrence instances
- `max_cache_instances` - Maximum instances to cache per event

### Logging

```yaml
log:
  level: "info"
  production: false
  rotation:
    max_size: 100
    max_backups: 5
    max_age: 30
```

- `level` - Log level (debug, info, warn, error, fatal)
- `production` - Use JSON format (true) or text (false)
- `rotation.max_size` - Max log file size in MB
- `rotation.max_backups` - Number of backup files
- `rotation.max_age` - Max age of log files in days

### Admin

```yaml
admin:
  statistics:
    enabled: true
    cache_duration: "5m"
  audit:
    enabled: true
    retention_days: 90
```

- `statistics.enabled` - Enable statistics collection
- `statistics.cache_duration` - How long to cache stats
- `audit.enabled` - Enable audit logging
- `audit.retention_days` - How long to keep audit logs

### Hash

```yaml
hash:
  algorithm: "argon2"  # bcrypt, argon2, or scrypt
  params:
    bcrypt_cost: 10
```

- `algorithm` - Password hashing algorithm
- `params.*` - Algorithm-specific parameters

### CORS

```yaml
cors:
  origin: "*"
  credentials: true
```

- `origin` - Allowed origin or "*" for all
- `credentials` - Allow credentials (cookies, auth headers)

## CLI Overrides

Configuration can also be overridden via CLI flags:

```bash
david server \
  --config config.yaml \
  --host 0.0.0.0 \
  --port 9000 \
  --debug \
  --hash-algorithm bcrypt
```

**Priority:** CLI flags > Environment variables > Config file > Defaults

## Security Recommendations

1. **JWT Secret:** Always use `DAVID_JWT_SECRET` environment variable
2. **TLS:** Enable TLS in production
3. **Rate Limiting:** Keep enabled to prevent abuse
4. **Audit Logging:** Enable for security monitoring
5. **CORS:** Restrict to specific origins in production

## Development vs Production

### Development
```yaml
log:
  level: "debug"
  production: false

api:
  rate_limiting:
    enabled: false  # Disable for development
```

### Production
```yaml
log:
  level: "info"
  production: true

admin:
  audit:
    enabled: true
    retention_days: 90
```
