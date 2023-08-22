#!/bin/sh
set -e

sed -i "s|REDIS_URL|${REDIS_URL}|" /etc/chihaya.yaml
sed -i "s|redis://:|redis://|" /etc/chihaya.yaml # Add dummy username
sed -i "s|0.0.0.0:8080|0.0.0.0:${PORT}|" /etc/chihaya.yaml
echo ""
cat /etc/chihaya.yaml
echo ""

exec /go/bin/chihaya
