#!/bin/bash

set -euo pipefail

if [[ ! -f "./cert.pem" || ! -f "./key.pem" ]]; then
    echo "CA certificate or key not found. Exiting."
    exit 1
fi

declare -A users
users=( ["admin"]="admin-db" ["gacha"]="gacha-db" ["market"]="market-db" ["user"]="user-db")

cd ../

for user in "${!users[@]}"; do
    CN="${users[$user]}"

    BASE_DIR="$user/postgres/certs"
    mkdir -p "$BASE_DIR"

    openssl genrsa -out "$BASE_DIR/server.key" 4096 || { echo "Failed to generate key for $user"; exit 1; }

    openssl req -new -key "$BASE_DIR/server.key" -out "$BASE_DIR/server.csr" -subj "/CN=$CN" || { echo "Failed to generate CSR for $user"; exit 1; }

    echo `pwd`

    openssl x509 -req -in "$BASE_DIR/server.csr" -CA ./cacerts/cert.pem -CAkey ./cacerts/key.pem -CAcreateserial -out "$BASE_DIR/server.crt" -days 365 || { echo "Failed to generate certificate for $user"; exit 1; }

    rm "$BASE_DIR/server.csr" || { echo "Failed to remove CSR file for $user"; exit 1; }
done

rm ./cacerts/cert.srl

read -p "Do you want to change ownership of the generated certificates to the postgres user (999) [y/n]? " answer

case $answer in
    [Yy]* )
        for user in "${!users[@]}"; do
            BASE_DIR="$user/postgres/certs"
            sudo chown 999:999 "$BASE_DIR/server."* || { echo "Failed to change ownership for $user"; exit 1; }
            echo "Ownership changed for $user."
        done
        ;;
    * )
        echo "Ownership not changed for any user."
        ;;
esac

echo "All certificates generated successfully."
