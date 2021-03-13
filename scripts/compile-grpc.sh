#!/usr/bin/env sh

set -e
set -x

cleanup() {
	docker rm extract-user-contract-grpc-builder
}

trap 'cleanup' EXIT

if [ $# -eq 0 ]; then
	current_directory=$(dirname "$0")
else
	current_directory="$1"
fi

cd "$current_directory"/..

docker build -f docker/Dockerfile.buildGrpcContract -t user-contract-grpc-builder .
docker create --name extract-user-contract-grpc-builder user-contract-grpc-builder
docker cp extract-user-contract-grpc-builder:/src/contract/grpc/go/user.pb.go ./contract/grpc/go

