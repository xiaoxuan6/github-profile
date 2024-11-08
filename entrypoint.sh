#!/bin/bash

set -e

if [ -f ".env" ];then
  ENV_GITHUB_TOKEN=$(cut -f2 -d"\"" .env)
  if [ -n "$ENV_GITHUB_TOKEN" ]; then
    GITHUB_TOKEN="$ENV_GITHUB_TOKEN"
  fi
fi

if [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_TOKEN environment variable not set"
  exit 1
fi

echo "GITHUB_TOKEN=$GITHUB_TOKEN" > .env

./github_profile