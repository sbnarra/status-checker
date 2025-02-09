#!/bin/bash

if [ "$(id -u)" -ne 0 ]; then
    echo "This installer must be run as root, elevating..."
    exec sudo bash "$0" "$@"
fi

echo "Uninstalling Status Checker!"

if [ -e /etc/systemd/system/status-checker.service ]; then
  systemctl stop status-checker
  systemctl disable status-checker
  rm /etc/systemd/system/status-checker.service
fi

rm -rf /opt/status-checker

echo "Status Checker Uninstalled!"