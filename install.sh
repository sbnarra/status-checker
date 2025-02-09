#!/bin/bash

INSTALL_DIR=${INSTALL_DIR:-/opt/status-checker}

if [ "$(id -u)" -ne 0 ]; then
    echo "This installer must be run as root, elevating..."
    exec sudo bash "$0" "$@"
fi

echo "Installing Status Checker!"

if [ -e /etc/systemd/system/status-checker.service ]; then
  echo "stopping systemd status checker..."
  systemctl stop status-checker
fi

set -e

arch=$(uname -m)
if [ "$arch" == "i386" ]; then arch="386"
elif [ "$arch" == "x86_64" ]; then arch="amd64"
elif [ "$arch" == "aarch64" ]; then arch="arm64"
elif [ "$arch" == "armv6l" ]; then arch="arm_v6"
elif [ "$arch" == "armv7l" ]; then arch="arm_v7"
# elif [ "$arch" == "mips64" ]; then arch="mips64"
# elif [ "$arch" == "mips64le" ]; then arch="mips64le"
# elif [ "$arch" == "ppc64le" ]; then arch="ppc64le"
# elif [ "$arch" == "riscv64" ]; then arch="riscv64"
# elif [ "$arch" == "s390x" ]; then arch="s390x"
else arch=$arch
fi

echo "...installing into $INSTALL_DIR"
mkdir -p $INSTALL_DIR

download_bin_url=$(curl -s https://api.github.com/repos/sbnarra/status-checker/releases/latest | jq -r '.assets[] | select(.name == "status-checker_'$arch'") | .browser_download_url')
curl -sL -o $INSTALL_DIR/status-checker $download_bin_url
chmod +x $INSTALL_DIR/status-checker
[ -e $INSTALL_DIR/checks.yaml ] || \
  curl -s -o $INSTALL_DIR/checks.yaml https://raw.githubusercontent.com/sbnarra/status-checker/refs/heads/main/config/checks.yaml
[ -e $INSTALL_DIR/config.env ] || \
  cat <<EOF >$INSTALL_DIR/config.env
CHECKS_PATH=$INSTALL_DIR/checks.yaml
BIND_ADDR=:9944
# SERVER_ENABLED=true
# DEBUG=true
# HISTORY_DIR=$INSTALL_DIR/history
# MIN_HISTORY=100
# HISTORY_CHECK_SIZE_LIMIT=10MB
# PROMETHEUS_ENABLED=true
# SLACK_HOOK_URL=
EOF

if [ -e /etc/systemd/system ]; then
  echo "...configuring systemd startup"
  cat <<EOF >/etc/systemd/system/status-checker.service
[Unit]
Description=Status Checker
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=$INSTALL_DIR/status-checker
EnvironmentFile=$INSTALL_DIR/config.env
StandardOutput=inherit
StandardError=inherit

[Install]
WantedBy=multi-user.target
EOF
  systemctl daemon-reload
  systemctl enable status-checker
  systemctl start status-checker
  RESTART_CMD="systemctl restart status-checker"
else
  echo "...didn't find service manager, starting now but will need manually starting on reboot"
  $INSTALL_DIR/status-checker $INSTALL_DIR/checks.yaml
  RESTART_CMD="kill $?; $INSTALL_DIR/status-checker $INSTALL_DIR/checks.yaml"
fi

cat <<EOF
Status Checker Installed!

...to view status checks open:
  http://localhost:9944
...to edit status checks open:
  $INSTALL_DIR/checks.yaml
...to edit runtime config open:
  $INSTALL_DIR/config.env
...restart service once done editing checks/config: 
  $RESTART_CMD
...to uninstall:
  curl -s https://raw.githubusercontent.com/sbnarra/status-checker/refs/heads/main/uninstall.sh | bash
EOF