#!/usr/bin/env bash
# Rebuild and restart the production compose stack on the EC2 host.
# Run from the repo root after `git pull`.
#
# Previous deploys pruned images *after* `docker compose up --build`.
# If the disk is already full, unpacking the frontend image fails with
# "no space left on device". Reclaim unused images and BuildKit cache
# first. Never prune volumes — Postgres data lives there.
set -euo pipefail

echo "Disk before Docker cleanup:"
df -h / || true
docker system df || true

docker container prune -f || true
docker image prune -af || true
docker builder prune -af || true

echo "Disk after Docker cleanup:"
df -h / || true
docker system df || true

docker compose pull
docker compose up -d --build --remove-orphans
docker image prune -f || true
