# Running the server binary with systemd (production)

The steps below assume a built binary located at `/opt/modulr-anchors/bin/modulr-anchors`.
Adjust the paths and user names to match your deployment target.

## 1) Create a dedicated user and directories
```bash
sudo useradd --system --home /opt/modulr-anchors --shell /usr/sbin/nologin modulr
sudo mkdir -p /opt/modulr-anchors/bin /opt/modulr-anchors/config /opt/modulr-anchors/logs
sudo chown -R modulr:modulr /opt/modulr-anchors
```

Place the compiled binary at `/opt/modulr-anchors/bin/modulr-anchors` and make it executable.
Keep configuration files (if any) under `/opt/modulr-anchors/config`.

## 2) Example systemd unit file
Create `/etc/systemd/system/modulr-anchors.service` with:

```ini
[Unit]
Description=Modulr Anchors server
After=network.target
Wants=network-online.target

[Service]
User=modulr
Group=modulr
WorkingDirectory=/opt/modulr-anchors
ExecStart=/opt/modulr-anchors/bin/modulr-anchors --config /opt/modulr-anchors/config/config.yaml
Restart=always
RestartSec=3
StartLimitIntervalSec=0
Environment=MODULR_ENV=production

# Hardening (adjust if the binary needs extra capabilities)
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=true
ReadWritePaths=/opt/modulr-anchors

[Install]
WantedBy=multi-user.target
```

### Notes
- Set `--config` and any other flags or environment variables your binary needs.
- If the binary binds to privileged ports or requires capabilities, remove/adjust the hardening options accordingly.

## 3) Reload systemd and start the service
```bash
sudo systemctl daemon-reload
sudo systemctl enable modulr-anchors.service
sudo systemctl start modulr-anchors.service
```

## 4) Monitoring and logs
- View status: `sudo systemctl status modulr-anchors.service`
- Follow logs (journald): `sudo journalctl -u modulr-anchors.service -f`
- If you still want file-based logs, write them under `/opt/modulr-anchors/logs` so they remain writable with `ReadWritePaths`.

## 5) Deploying updates
```bash
# Stop the service
sudo systemctl stop modulr-anchors.service

# Replace the binary
sudo install -m 0755 modulr-anchors /opt/modulr-anchors/bin/modulr-anchors
sudo chown modulr:modulr /opt/modulr-anchors/bin/modulr-anchors

# Start the service again
sudo systemctl start modulr-anchors.service
```

## 6) Optional: health checks and watchdog
- To enforce a periodic heartbeat, set `WatchdogSec=30s` and have the binary notify systemd via `sd_notify`.
- For graceful reloads, add an `ExecReload` command that your binary understands (e.g., sending `SIGHUP`).
