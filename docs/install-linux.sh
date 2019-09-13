#!/bin/bash

PMZERO_REPO=thenets/pmzero

if ! [ -x "$(command -v curl)" ]; then
    echo "[ERROR] Command 'curl' not found!"
    exit 1
fi

if [ -x "$(command -v sudo)" ]; then
    SUDO_PREFIX=sudo
fi

set -e
set -x

LATEST_TAG_NAME=$(curl --silent "https://api.github.com/repos/${PMZERO_REPO}/tags" | grep -Po '"name": "\K.*?(?=")' | head -1)

$SUDO_PREFIX rm -f /usr/bin/pmzero

$SUDO_PREFIX curl -o /usr/bin/pmzero -L https://github.com/${PMZERO_REPO}/releases/download/${LATEST_TAG_NAME}/pmzero-linux-x86_64

$SUDO_PREFIX chmod +x /usr/bin/pmzero
