# Site Pinger

# Install as service

Create `/usr/lib/systemd/user/pinger.service` file with content.

```
[Unit]
Description=Site Pinger
[Service]
Type=simple
ExecStart=/data/pinger/pingdaemon --c=/data/pinger/config.yml
[Install]
WantedBy=multi-user.target
```

Run to enable service
`systemctl enable /usr/lib/systemd/user/pinger.service`

Run to start service
`systemctl start /usr/lib/systemd/user/pinger.service`
