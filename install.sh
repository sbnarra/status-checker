#!/bin/bash
set -e

arch=$(uname -m)
if [ "$arch" == "i386" ]; then
  arch="386"
elif [ "$arch" == "x86_64" ]; then
  arch="amd64"
elif [ "$arch" == "armv6l" ]; then
  arch="arm_v6"
elif [ "$arch" == "armv7l" ]; then
  arch="arm_v7"
elif [ "$arch" == "aarch64" ]; then
  arch="arm64"
else
  arch=$arch
fi

download_url=$(curl -s https://api.github.com/repos/sbnarra/status-checks/releases/latest | jq -r '.assets[] | select(.name == "status-checker_'$arch'") | .url')

