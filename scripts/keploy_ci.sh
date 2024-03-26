#!/usr/bin/env bash

# Download and install keploy
curl --silent --location "https://github.com/keploy/keploy/releases/latest/download/keploy_linux_amd64.tar.gz" | tar xz -C /tmp
sudo mkdir -p /usr/local/bin && sudo mv /tmp/keploy /usr/local/bin && keploy

# Build the Go server
go build -cover -o main ./cmd/server/main.go

# Run tests with keploy
sudo -E keploy test -c "./main" --withCoverage --delay 10

# Generate coverage profile
go tool covdata textfmt -i="./keploy/coverage-reports" -o coverage-profile

# Check coverage against threshold
coverage_percentage=$(go tool cover -func=coverage-profile | grep 'total' | tail -n 1 | awk '{print $3}')
coverage_percentage=${coverage_percentage%\%}  # Remove the percentage sign
threshold=80
echo "Required threshold: ${threshold}%"
if (( $(awk -v num="$coverage_percentage" -v thresh="$threshold" 'BEGIN { if (num < thresh) print 1; else print 0 }') )); then
    echo "Coverage below threshold: $coverage_percentage%"
    exit 1
else
    echo "Coverage meets threshold: $coverage_percentage"
fi

