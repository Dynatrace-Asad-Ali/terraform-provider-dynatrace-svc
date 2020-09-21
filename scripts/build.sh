#!/usr/bin/env bash

VERSION=0.1.0-pre

# Delete the old dir
echo "==> Removing old directory..."
rm -rf bin/*
mkdir -p bin/windows/
mkdir -p bin/linux/
mkdir -p bin/darwin/

echo "==> Creating new binaries for windows, darwin and linux amd64..."
env GOOS=linux GOARCH=amd64 go build -o bin/linux/terraform-provider-dynatrace_v${VERSION}
zip bin/linux/terraform-provider-dynatrace_${VERSION}_linux_amd64.zip bin/linux/terraform-provider-dynatrace_v${VERSION}

env GOOS=darwin GOARCH=amd64 go build -o bin/darwin/terraform-provider-dynatrace_v${VERSION}
zip bin/darwin/terraform-provider-dynatrace_${VERSION}_darwin_amd64.zip bin/darwin/terraform-provider-dynatrace_v${VERSION}

env GOOS=windows GOARCH=amd64 go build -o bin/windows/terraform-provider-dynatrace_v${VERSION}
zip bin/windows/terraform-provider-dynatrace_${VERSION}_windows_amd64.zip bin/windows/terraform-provider-dynatrace_v${VERSION}