#!/bin/sh
set -e

sed -i "s|REDIS_URL|${REDIS_URL}|" /etc/chihaya.yaml
sed -i "s|redis://:|redis://x:|" /etc/chihaya.yaml # Add dummy username
echo ""
cat /etc/chihaya.yaml
echo ""

exec /go/bin/chihaya
