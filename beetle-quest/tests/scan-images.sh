#!/usr/bin/env bash

set -euo pipefail

images=(
    "beetle-quest-auth-service:latest"
    "beetle-quest-market-service:latest"
    "beetle-quest-user-service:latest"
    "beetle-quest-admin-service:latest"
    "beetle-quest-gacha-service:latest"
    "beetle-quest-static-service:latest"
    "beetle-quest-user-db:latest"
    "beetle-quest-gacha-db:latest"
    "beetle-quest-market-db:latest"
    "beetle-quest-admin-db:latest"
)

output_dir="./trivy_scan_results"
mkdir -p "$output_dir"



for image in "${images[@]}"; do
    output_file="$output_dir/${image//:latest/}_scan.txt"

    echo "Scanning image: $image"
    docker run --rm -v /var/run/docker.sock:/var/run/docker.sock:ro -v ./../../.trivy_cache:/.cache/ aquasec/trivy image -f table "$image" &> "$output_file"

    sed -n '14,$p' "$output_file" 2>&1 | tee "$output_dir/${image//:latest/}_table.txt"
    rm "$output_file"

    echo "Scan results saved to: $output_file"
done

echo "Scanning completed for all images."
