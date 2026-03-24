# Deployment Guide - david CalDAV Server

## Prerequisites

- Linux system (systemd-based)
- Go 1.21+ (only for building from source)
- 1GB+ RAM
- 10GB+ disk space
- OpenSSL for TLS certificates

## Quick Start

### Option 1: Binary Release (Recommended)

```bash
# Download latest release
curl -LO https://github.com/audstanley/david/releases/latest/download/david-linux-amd64

# Make executable
chmod +x david-linux-amd64

# Move to system path
sudo mv david-linux-amd64 /usr/local/bin/david

# Create user
sudo useradd -r -s /bin/false david
sudo mkdir -p /var/lib/david /var/log/david /etc/david
sudo chown -R david:david /var/lib/david /var/log/david

# Copy configuration
sudo cp config.yaml.production /etc/david/config.yaml

# Edit configuration
sudo nano /etc/david/config.yaml
# Update: JWT secret, API key, TLS certificates, data_dir

# Install systemd service
sudo cp contrib/systemd/david.service /etc/systemd/system/
sudo systemctl daemon-reload

# Start service
sudo systemctl enable david
sudo systemctl start david

# Check status
sudo systemctl status david
```

### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/audstanley/david.git
cd david

# Build
go build -o david -ldflags="-X main.version=1.0.0" ./cmd/david

# Run setup
./david setup

# Create admin user
./david user create admin

# Run server
./david server --config config.yaml
```

## Configuration

### Essential Settings

```yaml
# /etc/david/config.yaml

server:
  host: "0.0.0.0"
  port: 443
  tls:
    enabled: true
    cert_file: "/etc/ssl/certs/david.crt"
    key_file: "/etc/ssl/private/david.key"

data_dir: "/var/lib/david"

auth:
  jwt:
    secret: "<STRONG-RANDOM-SECRET>"
```

### Generate Secret

```bash
# Generate secure secret
head -c 32 /dev/urandom | base64
```

## TLS Setup

```bash
# Generate self-signed certificate (for testing)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/ssl/private/david.key \
  -out /etc/ssl/certs/david.crt \
  -subj "/CN=your-domain.com"

# Or use Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot certonly --standalone -d your-domain.com
```

## Firewall Configuration

```bash
# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Or for custom port
sudo ufw allow 8443/tcp

# Enable firewall
sudo ufw enable
```

## Reverse Proxy (Optional)

### Nginx Example

```nginx
server {
    listen 443 ssl http2;
    server_name cal.example.com;

    ssl_certificate /etc/ssl/certs/david.crt;
    ssl_certificate_key /etc/ssl/private/david.key;

    location /webdav/ {
        proxy_pass http://127.0.0.1:8443/webdav/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebDAV headers
        proxy_set_header Content-Length $content_length;
        proxy_buffering off;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8443/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Apache Example

```apache
<VirtualHost *:443>
    ServerName cal.example.com
    
    SSLEngine on
    SSLCertificateFile /etc/ssl/certs/david.crt
    SSLCertificateKeyFile /etc/ssl/private/david.key
    
    ProxyPass /webdav/ http://127.0.0.1:8443/webdav/
    ProxyPassReverse /webdav/ http://127.0.0.1:8443/webdav/
    
    ProxyPass /api/ http://127.0.0.1:8443/api/
    ProxyPassReverse /api/ http://127.0.0.1:8443/api/
</VirtualHost>
```

## Backup Strategy

```bash
# Create backup script
cat > /usr/local/bin/david-backup << 'EOF'
#!/bin/bash
BACKUP_DIR="/var/backups/david"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR
tar -czf $BACKUP_DIR/david_$TIMESTAMP.tar.gz /var/lib/david
find $BACKUP_DIR -name "david_*.tar.gz" -mtime +7 -delete
EOF

chmod +x /usr/local/bin/david-backup

# Add to cron
echo "0 2 * * * /usr/local/bin/david-backup" | sudo tee /etc/cron.d/david-backup
```

## Monitoring

### Prometheus Metrics

david exposes metrics at `/metrics`:

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'david'
    static_configs:
      - targets: ['localhost:8443']
    metrics_path: '/metrics'
```

### Health Check

```bash
# Check service status
curl -k https://localhost:8443/health

# Check metrics
curl -k https://localhost:8443/metrics
```

### Log Rotation

```bash
# Configure logrotate
cat > /etc/logrotate.d/david << 'EOF'
/var/log/david*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0640 david david
    postrotate
        systemctl reload david
    endscript
}
EOF
```

## Troubleshooting

### Service Won't Start

```bash
# Check logs
sudo journalctl -u david -f

# Check config
david server --config /etc/david/config.yaml --debug

# Check port
sudo ss -tlnp | grep 8443

# Check permissions
ls -la /var/lib/david
```

### Connection Issues

```bash
# Test WebDAV
curl -k https://localhost:8443/webdav/

# Test API
curl -k https://localhost:8443/api/health

# Check firewall
sudo ufw status
```

### Database Corruption

```bash
# Stop service
sudo systemctl stop david

# Backup
sudo cp -r /var/lib/david /var/lib/david.backup

# Try to recover
david backup --restore /var/lib/david.backup
```

## Security Best Practices

1. **Use strong secrets** - Generate with `head -c 32 /dev/urandom | base64`
2. **Enable TLS** - Always use HTTPS in production
3. **Regular updates** - Keep tzdata updated
4. **Monitor logs** - Set up log monitoring
5. **Backup regularly** - Automate backups
6. **Rate limiting** - Enable and configure
7. **Audit logging** - Keep audit logs enabled
8. **Firewall** - Restrict access to necessary ports
9. **Least privilege** - Run as non-root user
10. **Harden OS** - Apply security updates regularly

## Performance Tuning

### Increase File Descriptors

```bash
# Edit systemd service
cat > /etc/systemd/system/david.service.d/limits.conf << 'EOF'
[Service]
LimitNOFILE=65536
EOF

sudo systemctl daemon-reload
sudo systemctl restart david
```

### Cache Configuration

```yaml
# config.yaml
cache:
  enabled: true
  duration: "10m"
  max_size: 1073741824  # 1GB
```

## Uninstall

```bash
# Stop service
sudo systemctl stop david
sudo systemctl disable david

# Remove service
sudo rm /etc/systemd/system/david.service
sudo systemctl daemon-reload

# Remove files
sudo rm -rf /var/lib/david /var/log/david /etc/david
sudo rm /usr/local/bin/david
sudo userdel david
```

## Support

- Issues: https://github.com/audstanley/david/issues
- Documentation: https://github.com/audstanley/david/docs
- Community: Join our Discord/Slack (if available)
