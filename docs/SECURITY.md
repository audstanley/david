# Security Guide - david CalDAV Server

## Authentication

### JWT Tokens

- **Access tokens**: 1 hour expiry (configurable)
- **Refresh tokens**: 7 days expiry (configurable)
- **Secret**: Use strong random 32+ character secret
- **Blacklist**: Supported for token revocation

**Generate secure secret**:
```bash
head -c 32 /dev/urandom | base64
```

### API Keys

- Generated with bcrypt hash (cost 12)
- Prefix: `david_`
- Can be revoked per user
- Store securely - shown only once

### Basic Auth

- Username: User email
- Password: User password (bcrypt hashed)
- Recommended for CalDAV clients

## Password Security

### Requirements

- Minimum 8 characters
- Bcrypt hashing (default cost 12)
- Salting automatic

### Recommendations

- Use 16+ character passwords
- Enable password complexity validation
- Consider password manager integration

## Rate Limiting

### Default Settings

- 60 requests per minute
- 10 burst requests
- Per user/IP

### Configuration

```yaml
auth:
  rate_limit:
    requests_per_minute: 60
    burst: 10
    enabled: true
```

### Protection

- Prevents brute force attacks
- Mitigates DoS
- Per-user granularity

## TLS/SSL

### Required in Production

- Always enable TLS for production
- Use strong ciphers only
- Certificate valid 365+ days

### Cipher Configuration

```nginx
ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256';
ssl_prefer_server_ciphers on;
```

### Certificate Best Practices

1. Use Let's Encrypt for automatic renewal
2. Keep certificates updated
3. Use strong key lengths (2048+ RSA, 256+ ECC)
4. Enable HSTS

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
```

## Audit Logging

### Enabled by Default

- User login attempts
- Calendar operations
- Event modifications
- Admin actions

### Log Contents

- Timestamp
- User ID
- Action
- IP address
- Result

### Retention

- Default: Keep all logs
- Configure based on compliance needs

## Input Validation

### All Inputs Validated

- Calendar UIDs
- Event data
- User credentials
- API parameters

### Sanitization

- HTML escaped in descriptions
- XSS prevention
- SQL injection prevention (LevelDB)

## File Upload Security

### ICS Validation

- Parse before storing
- Validate required fields
- Reject malformed files
- Size limits configurable

### Content-Type Checking

```
text/calendar; charset=utf-8
```

## Database Security

### LevelDB Configuration

- File system permissions: 700
- Data directory: Private
- Backup encryption recommended

### Access Control

- File-level permissions
- Database user isolation
- Network isolation (localhost only)

## Network Security

### Firewall Rules

```bash
# Only allow necessary ports
ufw allow 443/tcp  # HTTPS
ufw deny from any to any port 22  # SSH (use bastion)
```

### Bind Address

- Production: `0.0.0.0` (with TLS)
- Development: `127.0.0.1` only

## Admin Security

### Role-Based Access

- **Admin**: Full access
- **Manager**: Calendar management
- **User**: Own calendars only

### Admin Account

- Use strong password
- Enable 2FA (if available)
- Limit admin accounts
- Monitor admin actions

## Headers & CSP

### Recommended Headers

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
```

## Backup Security

### Encrypted Backups

```bash
# Encrypt backup
gpg --encrypt --recipient admin@david.com backup.tar.gz
```

### Secure Storage

- Encrypted disk
- Offsite replication
- Access control

## Vulnerability Management

### Regular Updates

- Update dependencies monthly
- Monitor CVE notifications
- Test updates in staging

### Security Scanning

```bash
# Dependency scanning
go mod audit

# Static analysis
golangci-lint run
```

## Incident Response

### Breach Detection

- Monitor audit logs
- Alert on anomalies
- Track unusual patterns

### Response Steps

1. **Identify**: Determine scope
2. **Contain**: Isolate affected systems
3. **Eradicate**: Remove attack vector
4. **Recover**: Restore from backup
5. **Review**: Update security measures

### Notification

- Document all incidents
- Notify affected users
- Report to authorities if required

## Compliance

### GDPR Considerations

- Data encryption at rest
- User data export capability
- Right to deletion
- Data retention policies

### Data Minimization

- Collect only necessary data
- Regular data cleanup
- Anonymize old data

## Checklist

### Before Deployment

- [ ] TLS certificates configured
- [ ] Strong JWT secret set
- [ ] Rate limiting enabled
- [ ] Firewall rules configured
- [ ] Admin accounts secured
- [ ] Backup strategy in place
- [ ] Log monitoring configured
- [ ] Audit logging enabled

### Regular Maintenance

- [ ] Update dependencies monthly
- [ ] Rotate secrets annually
- [ ] Review audit logs
- [ ] Check certificate expiry
- [ ] Test backup restoration
- [ ] Review access logs
- [ ] Update firewall rules

## Resources

- OWASP Top 10: https://owasp.org/www-project-top-ten/
- TLS Best Practices: https://ssl-config.mozilla.org/
- Go Security: https://github.com/golang/go/wiki/Security
