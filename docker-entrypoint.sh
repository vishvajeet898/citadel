#!/bin/sh

set -e

sh ./fetch-and-validate-secrets.sh

exec "$@"
