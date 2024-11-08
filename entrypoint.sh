#!/bin/bash

set -e

if [ -z "$GITHUB_TOKEN" ]; then
  echo "GITHUB_TOKEN environment variable not set"
  exit 1
fi

echo "GITHUB_TOKEN=$GITHUB_TOKEN" > .env

./github_profile