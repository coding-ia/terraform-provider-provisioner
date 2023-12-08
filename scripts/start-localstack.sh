#!/bin/bash
set -e

export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-2
export AWS_SESSION_TOKEN=test

mkdir ~/.aws && cp ./scripts/fake.aws ~/.aws/credentials
docker-compose -f ./scripts/docker-compose.yaml up -d
sleep 30s
aws --endpoint-url=http://localhost:4566 sns create-topic --name go-test-topic