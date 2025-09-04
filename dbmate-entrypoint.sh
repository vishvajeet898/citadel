#!/bin/sh

set -e

s4cmd get s3://${S3_BUCKET}/citadel-microservice/migrate.env /.env

exec "$@"
