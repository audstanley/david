# Deployment Guide

## Prerequisites

- Go 1.24+ installed
- Linux/Unix system (also works on macOS and Windows)
- User with appropriate permissions
- 100MB+ disk space for database and files

## Installation

### Build from Source

```bash
# Clone repository
git clone https://github.com/audstanley/david
cd david

# Build binaries
mage Build

# Install binaries
sudo cp dist/david /usr/local/bin/
sudo cp dist/dcrypt /usr/local/bin/
sudo chmod +x /usr/local/bin/david
sudo chmod +x /usr/local/bin/dcrypt
```

### Download Pre-built Binary

Download from releases (coming soon) and extract:

```bash
tar xvf david-linux-amd64.tar.gz
sudo cp david /usr/local/bin/
```

## Configuration

### Create Configuration File

```bash
sudo mkdir -p /etc/david
sudo cp config.yaml.example /etc/david/config.yaml
sudo nano /etc/david/config.yaml
```

### Generate JWT Secret

```bash
# Generate a random secret
openssl rand -hex 32

# Add to environment or config
export DAVID_JWT_SECRET="<generated-secret>"
```

### Create Data Directory

```bash
sudo mkdir -p /var/lib/david
sudo chown -R $USER:$USER /var/lib/david
```

### Create First Admin User

```bash
# Generate password hash
dcrypt passwd --password "your-secure-password"

# Add to config.yaml
users:
  admin:
    password: "$2a$10$..."  # Hash from above
    role: "admin"
    email: "admin@example.com"
```

## systemd Setup

### Copy Service File

```bash
sudo cp contrib/systemd/david.service /etc/systemd/system/
sudo systemctl daemon-reload
```

### Enable and Start

```bash
sudo systemctl enable david
sudo systemctl start david
sudo systemctl status david
```

### View Logs

```bash
sudo journalctl -u david -f
```

## TLS Configuration

### Generate Self-Signed Certificate (for testing)

```bash
openssl req -x509 -newkey rsa:2048 \
  -keyout /etc/david/key.pem \
  -out /etc/david/cert.pem \
  -days 365 -nodes
```

### Update Configuration

```yaml
webdav:
  tls:
    enabled: true
    cert_file: "/etc/david/cert.pem"
    key_file: "/etc/david/key.pem"
```

### Use Production Certificate

For production, use a certificate from Let's Encrypt or your CA:

```bash
# Example with Certbot
sudo certbot certonly --standalone -d yourdomain.com
sudo systemctl reload david
```

## Reverse Proxy (Optional)

### nginx Example

```nginx
server {
    listen 80;
    server_name yourdomain.com;
    
    location /webdav {
        proxy_pass http://127.0.0.1:8000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebDAV headers
        proxy_set_header Content-Type application/octet-stream;
    }
}
```

### Caddy Example

```caddyfile
yourdomain.com {
    reverse_proxy localhost:8000
}
```

## Firewall Configuration

### UFW (Ubuntu)

```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8000/tcp  # If not using reverse proxy
sudo ufw enable
```

### firewalld (CentOS/RHEL)

```bash
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --permanent --add-port=8000/tcp
sudo firewall-cmd --reload
```

## Monitoring

### Health Check

```bash
curl http://localhost:8080/api/health
```

### Systemd Status

```bash
sudo systemctl status david
```

### Logs

```bash
# Recent logs
sudo journalctl -u david -n 100

# All logs
sudo journalctl -u david

# Live tail
sudo journalctl -u david -f
```

## Backup

### Database Backup

```bash
# Stop service
sudo systemctl stop david

# Backup databases
tar czf /backup/david-backup-$(date +%Y%m%d).tar.gz /var/lib/david/data

# Restart service
sudo systemctl start david
```

### File Backup

```bash
# If using WebDAV for storage
rsync -av /var/lib/david/ /backup/david-storage/
```

## Upgrading

```bash
# Stop service
sudo systemctl stop david

# Backup
sudo tar czf /backup/david-config.tar.gz /etc/david

# Install new version
sudo cp dist/david /usr/local/bin/

# Start service
sudo systemctl start david

# Verify
sudo systemctl status david
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
sudo lsof -i :8000

# Kill process or change port in config
```

### Permission Errors

```bash
# Fix directory permissions
sudo chown -R $USER:$USER /var/lib/david
sudo chmod -R 755 /var/lib/david
```

### TLS Errors

```bash
# Check certificate
openssl x509 -in /etc/david/cert.pem -text -noout

# Check key
openssl rsa -in /etc/david/key.pem -check
```

### Database Corruption

```bash
# Stop service
sudo systemctl stop david

# Try to recover (LevelDB has recovery built-in)
sudo systemctl start david

# If still failing, restore from backup
```

## Security Checklist

- [ ] JWT secret is unique and secure
- [ ] TLS is enabled in production
- [ ] Strong passwords for all users
- [ ] Rate limiting is enabled
- [ ] Audit logging is enabled
- [ ] Firewall is configured
- [ ] Regular backups are scheduled
- [ ] System is kept up to date
- [ ] Logs are monitored

## Support

- Documentation: `docs/` directory
- Issues: https://github.com/audstanley/david/issues
- Configuration: `docs/CONFIGURATION.md`
