#!/bin/bash
set -e
[ -f "$(dirname "$0")/.env" ] && source "$(dirname "$0")/.env"
echo "username=x-access-token"
echo "password=${KHAN_GH_TOKEN}"
