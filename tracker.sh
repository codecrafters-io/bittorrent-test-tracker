#!/bin/sh
set -e

sed -i "s|REDIS_URL|${REDIS_URL}|" /etc/chihaya.yaml
echo ""
cat /etc/chihaya.yaml
echo ""

exec /go/bin/chihaya
